package cookie

import (
	"AIGO/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

const CookieRefreshTokenKey = "refresh_token"

// SetRefreshTokenCookie 设置refreshToken的cookie
func SetRefreshTokenCookie(ctx *gin.Context, refreshToken string) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     CookieRefreshTokenKey,
		Value:    refreshToken,
		Path:     config.Cfg.AppCfg.Path,
		HttpOnly: true,
		MaxAge:   int(config.Cfg.JwtCfg.RefreshExpiration),
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
