package regx

import "testing"

func TestEamilRegex(t *testing.T) {
	email := "uzg108@outlook.com"
	matched, err := EamilRegex(email)
	if err != nil {
		t.Errorf("EamilRegex(%s) failed: %v", email, err)
	}
	if !matched {
		t.Errorf("EamilRegex(%s) failed: %v", email, err)
	}
}

func TestUsernameRegex(t *testing.T) {
	username := "野兽先辈"
	matched, err := UsernameRegex(username)
	if err != nil {
		t.Errorf("UsernameRegex(%s) failed: %v", username, err)
	}
	if !matched {
		t.Errorf("UsernameRegex(%s) failed: %v", username, err)
	}
}
