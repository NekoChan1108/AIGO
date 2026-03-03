package verification

import (
	"math/rand"
)

const (  
	// 验证码长度
	codeLength = 6
	// 随机种子
	randomSeed = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// GenerateVerificationCode 生成验证码
func GenerateVerificationCode() string {
	code := make([]byte, codeLength)
	// 生成随机数
	for i := range codeLength {
		// 填充随机数
		code[i] = randomSeed[rand.Intn(len(randomSeed))]
	}
	return string(code)
}
