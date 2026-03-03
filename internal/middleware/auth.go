package middleware

import (
	"AIGO/config"
	"AIGO/internal/model"
	"AIGO/pkg/utils/cookie"
	"AIGO/pkg/utils/jwt"
	"errors"
	"net/http"
	"strings"
	"time"

	v5 "github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

const (
	AuthUserKey          = "Username"
	NewAccessTokenHeader = "X-New-Access-Token"
)

// extractAccessToken 从请求头中提取token
func extractAccessToken(ctx *gin.Context) string {
	accessToken := ctx.GetHeader("Authorization")
	if accessToken == "" {
		return ""
	}
	// 去掉Bearer前缀
	return strings.TrimPrefix(strings.TrimSpace(accessToken), "Bearer ")
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取token
		accessToken := extractAccessToken(ctx)
		if accessToken == "" {
			ctx.JSON(http.StatusUnauthorized, &model.Response{
				Code: http.StatusUnauthorized,
				Msg:  "unauthorized",
				Data: nil,
			})
			// 中断请求
			ctx.Abort()
			return
		}
		// 校验token
		claims, err := jwt.ValidateTokenByType(accessToken, jwt.AccessTokenType)
		if err == nil {
			// 往上下文里添加用户信息
			ctx.Set(AuthUserKey, claims.Username)
			ctx.Next()
			return
		}
		// 非过期错误
		if !errors.Is(err, v5.ErrTokenExpired) {
			ctx.JSON(http.StatusUnauthorized, &model.Response{
				Code: http.StatusUnauthorized, // 通知前端重新登录获取新的token
				Msg:  err.Error(),
				Data: nil,
			})
			// 中断请求
			ctx.Abort()
			return
		}
		// 过期错误
		// 尝试续期 先从http only cookie中获取refreshToken
		refreshToken, err := ctx.Cookie(cookie.CookieRefreshTokenKey)
		if err != nil || refreshToken == "" {
			ctx.JSON(http.StatusUnauthorized, &model.Response{
				Code: http.StatusUnauthorized, // 通知前端重新登录获取新的token
				Msg:  "missing refresh cookie",
				Data: nil,
			})
			ctx.Abort()
			return
		}
		// 验证refreshToken
		start := time.Now()
		var newAccessToken, newRefreshToken string
		refreshClaims, err := jwt.ValidateTokenByType(refreshToken, jwt.RefreshTokenType)
		// 1、验证失败 只能退而求其次重新登录
		if err != nil || refreshClaims == nil {
			ctx.JSON(http.StatusUnauthorized, &model.Response{
				Code: http.StatusUnauthorized, // 通知前端重新登录获取新的token
				Msg:  "refresh token invalid, please login again",
				Data: nil,
			})
			// 中断请求
			ctx.Abort()
			return
		}
		// 2、验证成功 但是refreshtoken快过期
		if refreshClaims.ExpiresAt.Time.Sub(start) < time.Duration(config.Cfg.JwtCfg.RefreshExpiration/2)*time.Second {
			// 直接全体刷新
			newAccessToken, newRefreshToken, err = jwt.GenerateTokens(refreshClaims.Username, time.Now())
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, &model.Response{
					Code: http.StatusInternalServerError,
					Msg:  err.Error(),
					Data: nil,
				})
				// 中断请求
				ctx.Abort()
				return
			}
		} else {
			// 3、验证成功 refreshtoken有效期足够
			newAccessToken, err = jwt.RefreshAccessToken(refreshToken)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, &model.Response{
					Code: http.StatusInternalServerError,
					Msg:  err.Error(),
					Data: nil,
				})
				// 中断请求
				ctx.Abort()
				return
			}
		}
		if newRefreshToken != "" {
			// 刷新cookie
			cookie.SetRefreshTokenCookie(ctx, newRefreshToken)
		}
		// 设置新的accessToken 每次响应前端检查请求头是否含有 X-New-Access-Token 有就替换原有的accessToken
		ctx.Header(NewAccessTokenHeader, "Bearer "+newAccessToken)
		// 拿到新的accessToken添加到当前的请求中确保后续使用该中间件的handler能够继续进行请求而不是被拦截
		ctx.Request.Header.Set("Authorization", newAccessToken)
		ctx.Set(AuthUserKey, refreshClaims.Username)
		ctx.Next()
	}
}
