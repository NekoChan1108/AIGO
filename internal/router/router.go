package router

import (
	"AIGO/config"
	"AIGO/internal/handler"
	"AIGO/internal/handler/auth"
	"AIGO/internal/handler/file"
	"AIGO/internal/handler/llm"
	"AIGO/internal/middleware"
	"AIGO/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// R 全局路由
var Router *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 暴露 prometheus 指标 WrapH 将原生的 http handler 转为 Gin 的 handler
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	g := r.Group(config.Cfg.AppCfg.Path)
	g.Use(middleware.CorsMiddleWare())
	// 路由组
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, &model.Response{
			Code: http.StatusNotFound,
			Msg:  "not found",
			Data: nil,
		})
	})
	{
		gHello := g.Group("/hello")
		gHello.GET("", handler.HelloHandler)
	}
	// 用户路由组
	{
		gUser := g.Group("/auth")
		gUser.POST("/register", auth.UserRegisterHandler)
		gUser.POST("/login", auth.UserLoginHandler)
		gUser.POST("/verification", auth.VerificationCodeHandler)
	}
	// LLM路由组
	{
		gLLM := g.Group("/llm")
		gLLM.Use(middleware.AuthMiddleware())
		gLLM.GET("/session/list", llm.GetSessionByUsernameHandler())
		gLLM.POST("/session/create", llm.CreateSessionAndSendMessageHandler())
		gLLM.DELETE("/session/delete", llm.DeleteSessionHandler())
		gLLM.POST("/chat/send", llm.ChatSendHandler())
		gLLM.POST("/chat/history", llm.ChatHistoryHandler())
		gLLM.POST("/chat/create/stream", llm.CreateStreamingSessionAndSendMessageHandler())
		gLLM.POST("/chat/send/stream", llm.ChatStreamSendHandler())
	}
	// 文件路由组
	{
		gFile := g.Group("/file")
		gFile.Use(middleware.AuthMiddleware())
		gFile.POST("/upload", file.FileHandler())
	}
	Router = r
}
