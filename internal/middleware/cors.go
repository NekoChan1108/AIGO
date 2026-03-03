package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CorsMiddleWare 跨域中间件
func CorsMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin")

		// 处理请求头信息
		var headerKeys []string
		for k := range ctx.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}

		// 设置CORS头部
		if origin != "" {
			ctx.Header("Access-Control-Allow-Origin", "*") // 允许所有域，生产环境应具体指定
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			ctx.Header("Access-Control-Allow-Headers", "Authorization, X-New-Access-Token, Content-Length, X-CSRF-Token, Token, Session, X-Requested-With, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language, Cache-Control, Content-Type, Pragma")
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			ctx.Header("Access-Control-Max-Age", "172800") // 预检请求缓存时间（秒）
			ctx.Header("Access-Control-Allow-Credentials", "false")
		}

		// 处理OPTIONS预检请求
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
