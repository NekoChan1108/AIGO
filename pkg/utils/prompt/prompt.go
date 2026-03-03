package prompt

import (
	"fmt"
	"strings"

	"github.com/cloudwego/eino/schema"
)

const (
	defaultPrompt = "请你根据你自己的理解以及内部知识库中所有相关内容,再没有就进行网络搜索相关结果,然后直接回答用户问题"
)

// BuildPrompt 构建提示词 取不到内容就直接回答用户问题
func BuildPrompt(query string, docs []*schema.Document) string {
	// 构建提示词
	if len(docs) <= 0 {
		return defaultPrompt
	}
	contextContent := &strings.Builder{}
	// 从topk文档中提取内容
	for idx, doc := range docs {
		n, err := fmt.Fprintf(contextContent, "文档%d: %s\n", idx+1, doc.Content)
		if err != nil || n <= 0 {
			return defaultPrompt
		}
	}
	// 构建提示词
	prompt := fmt.Sprintf("根据以下文档回答用户提出的问题,参考文档:%s\n 用户问题:%s \n 请你提供准确的、完整的回答并且在有必要时引用参考文档中的内容.\n如果文档中没有相关信息就%s.\n", contextContent.String(), query, defaultPrompt)
	return prompt
}
