package password

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestPassword(t *testing.T) {
	passwordHashed, err := base64.StdEncoding.DecodeString("ZTk1YzEzYzQ3YmUyOTIxZWVmMmRkNjQwY2Q3NGFhZjNhZjk1MDI2MjRjZjE4Y2NlZjc0NDY2YjhjMmJlMWYyYw==")
	if err != nil {
		t.Errorf("decode base64 failed:%v", err)
	}
	t.Logf("----password hash=%s", string(passwordHashed))
	salt := "WkvWoO7OjUwOaL+M2yq2fQ=="
	saltByte, _ := base64.StdEncoding.DecodeString(salt)
	hash := HashPassword("%RrMki8q", saltByte)
	t.Logf("----password hash=%s", hex.EncodeToString(hash))
}
