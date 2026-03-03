package regx

import (
	"fmt"
	"regexp"
)

// EamilRegex 邮箱正则校验
func EamilRegex(email string) (bool, error) {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if err != nil {
		return false, fmt.Errorf("match email failed: %v", err)
	}
	return matched, nil
}

// UsernameRegex 用户名正则校验
func UsernameRegex(username string) (bool, error) {
	// 支持中文、英文等任意unicode字符、数字、下划线，长度2-20位
	matched, err := regexp.MatchString(`^[\p{L}\p{N}_]{2,20}$`, username)
	if err != nil {
		return false, fmt.Errorf("match username failed: %v", err)
	}
	return matched, nil
}
