package utils

import "testing"

func TestGetIP(t *testing.T) {
	t.Logf("---local ip=%s", GetIP())
}
