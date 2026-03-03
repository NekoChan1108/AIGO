package llm

import (
	"AIGO/internal/llm"
	"AIGO/internal/logic"
	"AIGO/internal/middleware"
	"AIGO/internal/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	CreateSessionAndSendMessageReq struct {
		Question string `json:"question"`
		LLMType  string `json:"llm_type"`
	}
	CreateStreamingSessionAndSendMessageReq struct {
		Question string `json:"question"`
		LLMType  string `json:"llm_type"`
	}
	CreateSessionAndSendMessageResp struct {
		Content   string `json:"content"`
		SessionID string `json:"session_id"`
	}
	DeleteSessionReq struct {
		SessionID string `json:"session_id"`
	}
	ChatSendReq struct {
		Question  string `json:"question"`
		SessionID string `json:"session_id"`
		LLMType   string `json:"llm_type"`
	}
	ChatSendResp struct {
		Content string `json:"content"`
	}
	ChatStreamSendReq struct {
		Question  string `json:"question"`
		SessionID string `json:"session_id"`
		LLMType   string `json:"llm_type"`
	}
	ChatHistoryReq struct {
		SessionID string `json:"session_id"`
	}
)

// setSSEHeaders 设置SSE响应头
func setSSEHeaders(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream") // 设置Content-Type为SSE
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked") // 启用分块传输编码
	ctx.Header("X-Accel-Buffering", "no")      // 禁用Nginx缓冲
}

// GetSessionByUsernameHandler 获取用户会话
func GetSessionByUsernameHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString(middleware.AuthUserKey)
		sessions, err := logic.GetSessionByUsername(ctx, username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, &model.Response{
			Code: http.StatusOK,
			Msg:  "success",
			Data: sessions,
		})
	}
}

// CreateSessionAndSendMessageHandler 创建普通会话并发送消息
func CreateSessionAndSendMessageHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString(middleware.AuthUserKey)
		var req CreateSessionAndSendMessageReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, &model.Response{
				Code: http.StatusBadRequest,
				Msg:  "invalid params",
				Data: nil,
			})
			return
		}
		sessionID, content, err := logic.CreateSessionAndSendMessage(ctx, req.Question, username, llm.LLMType(req.LLMType))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, &model.Response{
			Code: http.StatusOK,
			Msg:  "success",
			Data: &CreateSessionAndSendMessageResp{
				Content:   content,
				SessionID: sessionID,
			},
		})
	}
}

// CreateStreamingSessionAndSendMessageHandler 创建流式会话并发送消息
func CreateStreamingSessionAndSendMessageHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString(middleware.AuthUserKey)
		var req CreateStreamingSessionAndSendMessageReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, &model.Response{
				Code: http.StatusBadRequest,
				Msg:  "invalid params",
				Data: nil,
			})
			return
		}
		// 设置SSE头
		setSSEHeaders(ctx)
		// 创建会话
		sessionID, err := logic.CreateStreamSession(ctx, req.Question, username)
		if err != nil {
			ctx.SSEvent("error", &model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		// 先将会话ID返回给前端 用于侧边栏展示
		fmt.Fprintf(ctx.Writer, "data: {\"session_id\": \"%s\"}\n\n", sessionID)
		ctx.Writer.Flush()
		// 开始发送流式响应
		err = logic.StreamMessage(ctx, req.Question, username, sessionID, http.ResponseWriter(ctx.Writer), llm.LLMType(req.LLMType))
		if err != nil {
			ctx.SSEvent("error", &model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
	}
}

// DeleteSessionHandler 删除会话
func DeleteSessionHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req DeleteSessionReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, &model.Response{
				Code: http.StatusBadRequest,
				Msg:  "invalid params",
				Data: nil,
			})
			return
		}
		err := logic.DeleteSession(ctx, req.SessionID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, &model.Response{
			Code: http.StatusOK,
			Msg:  "success",
			Data: nil,
		})
	}
}

// ChatSendHandler 聊天 会话中返回对问题的解答
func ChatSendHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString(middleware.AuthUserKey)
		var req ChatSendReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, &model.Response{
				Code: http.StatusBadRequest,
				Msg:  "invalid params",
				Data: nil,
			})
			return
		}
		content, err := logic.ChatSend(ctx, req.Question, username, req.SessionID, llm.LLMType(req.LLMType))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, &model.Response{
			Code: http.StatusOK,
			Msg:  "success",
			Data: &ChatSendResp{
				Content: content,
			},
		})
	}
}

// ChatHistoryHandler 聊天 会话中返回所有问题和回答
func ChatHistoryHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ChatHistoryReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, &model.Response{
				Code: http.StatusBadRequest,
				Msg:  "invalid params",
				Data: nil,
			})
			return
		}
		history, err := logic.ChatHistory(ctx, req.SessionID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, &model.Response{
			Code: http.StatusOK,
			Msg:  "success",
			Data: history,
		})
	}
}

// ChatStreamSendHandler 聊天 会话中返回对问题的流式解答
func ChatStreamSendHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString(middleware.AuthUserKey)
		var req ChatStreamSendReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, &model.Response{
				Code: http.StatusBadRequest,
				Msg:  "invalid params",
				Data: nil,
			})
			return
		}
		// 设置SSE头
		setSSEHeaders(ctx)
		// 开始发送流式响应
		err := logic.StreamMessage(ctx, req.Question, username, req.SessionID, http.ResponseWriter(ctx.Writer), llm.LLMType(req.LLMType))
		if err != nil {
			ctx.SSEvent("error", &model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
	}
}
