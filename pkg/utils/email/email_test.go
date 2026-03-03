package email

import "testing"

func TestWelcomeEmail(t *testing.T) {
	SendWelcomeEmail("XXXXXXX@XXXXXX.com", "YajuSenpai")
}

func TestSendVerificationEmail(t *testing.T) {
	SendVerificationEmail("XXXXXXX@XXXXXX.com", "114514")
}
