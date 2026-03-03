package encrypt

import "testing"

func TestMD5(t *testing.T) {
	t.Log(MD5Encrypt("123456"))
}
