package jwt

import (
	"AIGO/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO 添加token黑名单机制

// tokenType token类型
type tokenType string

const (
	AccessTokenType  tokenType = "access_token"
	RefreshTokenType tokenType = "refresh_token"
)

// TokenClaims 自定义的claims
type TokenClaims struct {
	Username  string    // 用户名
	TokenType tokenType // token类型
	jwt.RegisteredClaims
}

// GenAccessToken 生成访问token
/*
@description: 生成访问token
@param username 用户名
@param now 生成token的时间 签发时间
@return string token
@return error 错误
*/
func GenAccessToken(username string, now time.Time) (string, error) {
	claims := TokenClaims{
		Username:  username,
		TokenType: AccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(config.Cfg.JwtCfg.AccessExpiration) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    config.Cfg.JwtCfg.Issuer,
			Subject:   config.Cfg.JwtCfg.Subject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.Cfg.JwtCfg.Secret))
	if err != nil {
		return "", fmt.Errorf("gen access token failed: %v", err)
	}
	return tokenStr, nil
}

// GenRefreshToken 生成刷新token
/*
@description: 生成刷新token
@param username 用户名
@param now 签发时间
@return string token
@return error 错误
*/
func GenRefreshToken(username string, now time.Time) (string, error) {
	claims := TokenClaims{
		Username:  username,
		TokenType: RefreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(config.Cfg.JwtCfg.RefreshExpiration) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    config.Cfg.JwtCfg.Issuer,
			Subject:   config.Cfg.JwtCfg.Subject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.Cfg.JwtCfg.Secret))
	if err != nil {
		return "", fmt.Errorf("gen refresh token failed: %v", err)
	}
	return tokenStr, nil
}

// GenerateTokens 生成双token
/*
@description:
@param username 用户名
@param now 签发时间
@return string token
@return error 错误
*/
func GenerateTokens(username string, now time.Time) (string, string, error) {
	iat := time.Now()
	accessToken, err := GenAccessToken(username, iat)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := GenRefreshToken(username, iat)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// parseToken 解析token
func parseToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(*jwt.Token) (any, error) {
		return []byte(config.Cfg.JwtCfg.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token error: %v", err)
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token is invalid")
}

// RefreshAccessToken 根据刷新token刷新访问token
func RefreshAccessToken(refreshTokenStr string) (string, error) {
	claims, err := parseToken(refreshTokenStr)
	if err != nil {
		return "", err
	}
	// 判断是否是刷新token
	if claims.TokenType != RefreshTokenType {
		return "", fmt.Errorf("token is not refresh token")
	}
	// 重新生成访问token
	return GenAccessToken(claims.Username, time.Now())
}

// ValidateTokenByType 验证token
func ValidateTokenByType(tokenStr string, tokenType tokenType) (*TokenClaims, error) {
	claims, err := parseToken(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != tokenType {
		return nil, fmt.Errorf("expected %v but got %v", tokenType, claims.TokenType)
	}
	if claims.Issuer != config.Cfg.JwtCfg.Issuer {
		return nil, fmt.Errorf("invalid issuer")
	}
	if claims.Subject != config.Cfg.JwtCfg.Subject {
		return nil, fmt.Errorf("invalid subject")
	}
	// 校验签发时间（防止伪造过去的Token）
	if claims.IssuedAt.Time.After(time.Now()) {
		return nil, fmt.Errorf("token issued in the future")
	}
	return claims, nil
}
