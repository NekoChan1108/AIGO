package dao

import (
	"AIGO/config"
	"AIGO/pkg/db"
	"AIGO/pkg/log"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	indexer "github.com/cloudwego/eino-ext/components/indexer/milvus2"
	retriever "github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// TODO 搜索模式的完善 考虑混合搜索 目前只支持余弦相似度 自动解析不同格式文档来选择分割器 http://www.cloudwego.cn/zh/docs/eino/ecosystem_integration/document/

const maxBufferSize = 2 << 20 // 2MB 流式读取单次最大读取大小
const chunkSize = 2 << 9      // 每个文档片段的最大长度
const overlapSize = 2 << 7    // 重叠大小 确保上下文不丢失
// 语义分割符 用于将文档分割为片段 确保上下文不丢失
var (
	semanticRuneSymbol   = []rune{'\n', '.', '。', '!', '！', '?', '？'}
	semanticStringSymbol = []string{"\n\n", "\n", ".", "。", "!", "！", "?", "？"}
)

var (
	embedder                     *ark.Embedder
	ragManager                   *RagManager
	ragManagerOnce, embedderOnce sync.Once // 保证全局一个RagManager实例以及全局一个嵌入器实例
)

// RagManager RAG管理器
type RagManager struct {
	lock         sync.RWMutex                        // 锁 确保对索引器映射的并发访问安全
	embedderMap  map[string]map[string]*ark.Embedder // 嵌入器映射 模型名称  值为嵌入器
	indexerMap   map[string]map[string]*RagIndexer   // 索引器映射 用户名+会话ID  值为索引器
	retrieverMap map[string]map[string]*RagRetriever // 检索器映射 用户名+会话ID  值为检索器
}

// RagIndexer RAG索引器
type RagIndexer struct {
	indexer   *indexer.Indexer
	spliter   document.Transformer // 文档分割器
	username  string               // 用户名 每个用户一个自己的文件夹
	sessionID string               // 会话ID 每个会话一个自己的文件夹
}

// RagRetriever RAG检索器
type RagRetriever struct {
	retriever *retriever.Retriever
	username  string // 用户名 每个用户一个自己的文件夹
	sessionID string // 会话ID 每个会话一个自己的文件夹
}

// newRagManager 创建RAG管理器
func newRagManager() *RagManager {
	return &RagManager{
		lock:         sync.RWMutex{},
		embedderMap:  make(map[string]map[string]*ark.Embedder),
		indexerMap:   make(map[string]map[string]*RagIndexer),
		retrieverMap: make(map[string]map[string]*RagRetriever),
	}
}

// newRagIndexer 创建RAG索引器
func newRagIndexer(ctx context.Context, username, sessionID string) (*RagIndexer, error) {
	// 创建文档分割器
	spliter, err := newSpliter(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create spliter: %v", err)
	}
	// 创建索引器
	indexer, err := indexer.NewIndexer(ctx, &indexer.IndexerConfig{
		Client:     db.MilvusDB,
		Collection: config.Cfg.MilvusCfg.Collection,
		Embedding:  embedder,
		// 向量配置
		Vector: &indexer.VectorConfig{
			Dimension:    config.Cfg.MilvusCfg.Dimension, // 与配置的向量模型的维度一致
			MetricType:   indexer.COSINE,                 // 余弦相似度
			IndexBuilder: indexer.NewAutoIndexBuilder(),  // 自动选择索引类型
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create milvus indexer: %v", err)
	}
	return &RagIndexer{
		indexer:   indexer,
		spliter:   spliter,
		username:  username,
		sessionID: sessionID,
	}, nil
}

// newSpliter 文档分割器
func newSpliter(ctx context.Context) (document.Transformer, error) {
	recursiveSplitter, err := recursive.NewSplitter(ctx, &recursive.Config{
		ChunkSize:   chunkSize,                                    // 分割的块大小
		OverlapSize: overlapSize,                                  // 重叠大小 确保上下文不丢失
		Separators:  semanticStringSymbol,                         // 分隔符
		LenFunc:     func(s string) int { return len([]rune(s)) }, // 字符串长度函数 采用rune长度计算 避免中文等多字节字符计算错误
		KeepType:    recursive.KeepTypeNone,                       // 结尾处不保留分隔符
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create recursive splitter: %v", err)
	}
	return recursiveSplitter, nil
}

// IndexFile 索引文件
func (i *RagIndexer) IndexFile(ctx context.Context, fileId string, file io.Reader) error {
	// 索引文件
	// 先读取用户上传的文件
	bufReader := bufio.NewReader(file)
	// 当前切片的内容
	var curSegment strings.Builder
	// 当前切片的索引
	var curSegmentIdx int
	for {
		// 单字符读取 避免oom
		r, _, err := bufReader.ReadRune()
		if err != nil {
			// 读取到文件末尾 跳出循环
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("failed to read file: %v", err)
		}
		// 写入当前切片
		curSegment.WriteRune(r)
		// 切到最大程度的同时保持语义
		if curSegment.Len() >= maxBufferSize && slices.Contains(semanticRuneSymbol, r) {
			// 写入当前切片 进行索引
			if err := i.processSegment(ctx, fileId, curSegment.String(), curSegmentIdx); err != nil {
				return fmt.Errorf("failed to process segment: %v", err)
			}
			// 重置当前切片
			curSegment.Reset()
			curSegmentIdx++
		}
	}
	// 处理剩余内容
	if curSegment.Len() > 0 {
		// 写入当前切片 进行索引
		if err := i.processSegment(ctx, fileId, curSegment.String(), curSegmentIdx); err != nil {
			return fmt.Errorf("failed to process rest segment: %v", err)
		}
	}
	return nil
}

// processSegment 处理切片
func (i *RagIndexer) processSegment(ctx context.Context, fileId, content string, idx int) error {
	doc := &schema.Document{
		ID:      fmt.Sprintf("%s_seg_%d", fileId, idx), // id确保唯一
		Content: content,
		MetaData: map[string]any{
			// 元数据存储额外信息
			"username":   i.username,
			"session_id": i.sessionID,
			"file_id":    fileId,
			"created_at": time.Now().Format("2006-01-02 15:04:05"),
		},
	}
	// 切分文档
	chunks, err := i.spliter.Transform(ctx, []*schema.Document{doc})
	if err != nil {
		return fmt.Errorf("failed to split document: %v", err)
	}
	for chunkIdx, chunk := range chunks {
		chunk.ID = fmt.Sprintf("%s_seg_%d_chunk_%d", fileId, idx, chunkIdx)
		if chunk.MetaData == nil {
			chunk.MetaData = make(map[string]any)
		}
		// 赋值
		maps.Copy(chunk.MetaData, doc.MetaData)
		// 添加元数据
		chunk.MetaData["chunk_idx"] = chunkIdx
		chunk.MetaData["segment_idx"] = idx
	}
	// 索引文档
	if _, err = i.indexer.Store(ctx, chunks); err != nil {
		return fmt.Errorf("failed to index documents: %v", err)
	}
	return nil
}

// DeleteIndex 删除会话索引
func (i *RagIndexer) DeleteIndex(ctx context.Context) error {
	res, err := db.MilvusDB.Delete(ctx, milvusclient.NewDeleteOption(config.Cfg.MilvusCfg.Collection).WithExpr(
		fmt.Sprintf("metadata[\"username\"] == '%s' AND metadata[\"session_id\"] == '%s'", i.username, i.sessionID),
	))
	if err != nil {
		return fmt.Errorf("failed to delete document from milvus: %v", err)
	}
	if res.DeleteCount <= 0 {
		log.Infof("no documents deleted")
	}
	return nil
}

// newRagRetriever 创建RAG检索器
func newRagRetriever(ctx context.Context, username, sessionID string) (*RagRetriever, error) {
	// 判断当前会话是否有上传文件 查一条并且只返回id就行不要全查浪费时间提高效率
	res, err := db.MilvusDB.Query(ctx, milvusclient.NewQueryOption(config.Cfg.MilvusCfg.Collection).WithFilter(
		fmt.Sprintf("metadata[\"username\"] == '%s' AND metadata[\"session_id\"] == '%s'", username, sessionID),
	).WithLimit(1).WithOutputFields([]string{"id"}...))
	if err != nil {
		return nil, fmt.Errorf("failed to query milvus: %v", err)
	}
	if res.ResultCount <= 0 {
		return nil, nil
	}
	// 创建检索器
	retriever, err := retriever.NewRetriever(ctx, &retriever.RetrieverConfig{
		Client:     db.MilvusDB,
		Collection: config.Cfg.MilvusCfg.Collection,
		Embedding:  embedder,
		TopK:       config.Cfg.MilvusCfg.TOPK,                    // 召回时选择最接近的答案数
		SearchMode: search_mode.NewApproximate(retriever.COSINE), // 搜索模式
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create milvus retriever: %v", err)
	}
	return &RagRetriever{
		retriever: retriever,
		username:  username,
		sessionID: sessionID,
	}, nil
}

// RetrieveFile 从Milvus检索文档
func (r *RagRetriever) RetrieveFile(ctx context.Context, query string) ([]*schema.Document, error) {
	// 从Milvus检索文档取出topk个最接近query的答案
	docs, err := r.retriever.Retrieve(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documents from milvus: %v", err)
	}
	return docs, nil
}

// GetGlobalRagManager 获取全局RAG管理器
func GetGlobalRagManager() *RagManager {
	ragManagerOnce.Do(func() {
		ragManager = newRagManager()
	})
	// double check
	if ragManager == nil {
		ragManager = newRagManager()
	}
	return ragManager
}

// GetOrCreateRagIndexer 获取或创建RAG索引器
func GetOrCreateRagIndexer(ctx context.Context, username, sessionID string) (indexer *RagIndexer, err error) {
	// 获取全局RAG管理器
	ragManager := GetGlobalRagManager()
	if ragManager == nil {
		return nil, fmt.Errorf("rag manager is nil")
	}
	ragManager.lock.RLock()
	// 检查索引器是否已存在
	ragIndexer, ok := ragManager.indexerMap[username][sessionID]
	if ok {
		ragManager.lock.RUnlock()
		return ragIndexer, nil
	}
	ragManager.lock.RUnlock()
	ragManager.lock.Lock()
	defer ragManager.lock.Unlock()
	// 检查索引器是否已存在 double check
	ragIndexer, ok = ragManager.indexerMap[username][sessionID]
	if ok {
		return ragIndexer, nil
	}
	// 初始化索引器
	indexer, err = newRagIndexer(ctx, username, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create rag indexer: %w", err)
	}
	// 存储索引器
	if ragManager.indexerMap[username] == nil {
		ragManager.indexerMap[username] = make(map[string]*RagIndexer)
	}
	ragManager.indexerMap[username][sessionID] = indexer
	return indexer, nil
}

// GetRagRetriever 获取RAG检索器
func GetRagRetriever(ctx context.Context, username, sessionID string) (retriever *RagRetriever, err error) {
	// 获取全局RAG管理器
	ragManager := GetGlobalRagManager()
	if ragManager == nil {
		return nil, fmt.Errorf("rag manager is nil")
	}
	ragManager.lock.RLock()
	// 检查检索器是否已存在
	retriever, ok := ragManager.retrieverMap[username][sessionID]
	if ok {
		ragManager.lock.RUnlock()
		return retriever, nil
	}
	ragManager.lock.RUnlock()
	ragManager.lock.Lock()
	defer ragManager.lock.Unlock()
	// 检查检索器是否已存在 double check
	retriever, ok = ragManager.retrieverMap[username][sessionID]
	if ok {
		return retriever, nil
	}
	// 初始化检索器
	retriever, err = newRagRetriever(ctx, username, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create rag retriever: %w", err)
	}
	if retriever == nil {
		return nil, nil
	}
	// 存储检索器
	if ragManager.retrieverMap[username] == nil {
		ragManager.retrieverMap[username] = make(map[string]*RagRetriever)
	}
	ragManager.retrieverMap[username][sessionID] = retriever
	return retriever, nil
}

// 初始化Embedder
func init() {
	embedderOnce.Do(func() {
		var err error
		embedder, err = ark.NewEmbedder(context.Background(), &ark.EmbeddingConfig{
			Model:  config.Cfg.ModelCfg.EmbeddingModel,
			APIKey: config.Cfg.ModelCfg.ApiKey,
		})
		if err != nil {
			log.Panic(fmt.Errorf("failed to create volcengine embedding model: %w", err))
		}
	})
}
