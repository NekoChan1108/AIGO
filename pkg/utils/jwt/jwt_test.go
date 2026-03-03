package jwt

import (
	"testing"
	"time"
)

var s1, s2 string

func TestGenToken(t *testing.T) {
	s1, s2, _ = GenerateTokens("test", time.Now())
	t.Log(s1, s2)
}

func TestValidateToken(t *testing.T) {
	t.Log(ValidateTokenByType(s1, AccessTokenType))
	t.Log(ValidateTokenByType(s2, RefreshTokenType))
}
