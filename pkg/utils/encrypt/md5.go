package encrypt

import (
	"crypto/md5"
	"fmt"
)

const salt = "AIGO"


// MD5Encrypt md5加密
func MD5Encrypt(str string) (string, error) {
	hash := md5.New()
	b := []byte(str + salt)
	_, err := hash.Write(b)
	if err != nil {
		return "", fmt.Errorf("md5 encrypt failed: %v", err)
	}
	return fmt.Sprintf("%x", hash.Sum(b)), nil
}
