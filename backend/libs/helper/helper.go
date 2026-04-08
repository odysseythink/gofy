package helper

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateTextHash(text string) string {
	hashText := text + "None"
	hash := sha256.Sum256([]byte(hashText))
	return hex.EncodeToString(hash[:])
}
