package password

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"regexp"

	"golang.org/x/crypto/pbkdf2"

	"mlib.com/gofy/server/core/exceptions"
)

var passwordPattern = `^(?=.*[a-zA-Z])(?=.*\d).{8,}$`

// ValidPassword checks if the password matches the required pattern.
func ValidPassword(password string) (string, error) {
	pattern := regexp.MustCompile(passwordPattern)
	if pattern.MatchString(password) {
		return password, nil
	}
	return "", exceptions.NewValueError("Password must contain letters and numbers, and the length must be greater than 8.")
}

// HashPassword hashes the password with the given salt.
func HashPassword(passwordStr string, saltByte []byte) []byte {
	dk := pbkdf2.Key([]byte(passwordStr), saltByte, 10000, 32, sha256.New)
	return dk
	// return hex.EncodeToString(dk)
}

// ComparePassword compares the provided password with the hashed password.
func ComparePassword(passwordStr, passwordHashedBase64, saltBase64 string) bool {
	saltByte, _ := base64.StdEncoding.DecodeString(saltBase64)
	hashed := HashPassword(passwordStr, saltByte)
	passwordHashed, _ := base64.StdEncoding.DecodeString(passwordHashedBase64)
	return hex.EncodeToString(hashed) == string(passwordHashed)
}

// GenerateRefreshToken generates a secure random token of the specified length.
func GenerateRefreshToken(length int) string {
	// Generate random bytes
	if length <= 0 {
		length = 64
	}
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)

	// Convert random bytes to hexadecimal string
	token := hex.EncodeToString(randomBytes)
	return token
}
