package router

import (
	"AIGO/internal/handler"

	"github.com/gin-gonic/gin"
)

// R 全局路由
var Router *gin.Engine

func init() {
	r := gin.Default()
	g := r.Group("/api")
	// 路由组
	{
		gHello := g.Group("/hello")
		gHello.GET("", handler.HelloHandler)
	}
	{
		gUser := g.Group("/user")
		gUser.POST("/register", handler.UserRegisterHandler)
		gUser.POST("/login", handler.UserLoginHandler)
	}
	Router = r
}
