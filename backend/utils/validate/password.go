package validate

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/odysseythink/mlog"
	"golang.org/x/crypto/pbkdf2"
)

// // 定义密码规则的正则表达式
// var passwordPattern = regexp.MustCompile(`^(?=.*[a-zA-Z])(?=.*\d).{8,}$`)

// // ValidPassword 检查密码是否符合规则
// func ValidPassword(password string) (string, error) {
// 	if passwordPattern.MatchString(password) {
// 		return password, nil
// 	}
// 	return "", errors.New("not a valid password")
// }

// HashPassword 使用PBKDF2 HMAC SHA256算法对密码进行哈希
func HashPassword(passwordStr string, salt []byte) []byte {
	// dk = hashlib.pbkdf2_hmac("sha256", password_str.encode("utf-8"), salt_byte, 10000)
	// return binascii.hexlify(dk)
	return pbkdf2.Key([]byte(passwordStr), salt, 10000, 32, sha256.New)
	// return hex.EncodeToString(dk)
	// dk := pbkdf2Hmac([]byte(passwordStr), salt, 10000, 32)
	// return dk
}

// comparePassword 比较密码是否匹配
func ComparePassword(passwordStr, passwordHashedBase64, saltBase64 string) bool {
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return false
	}
	mlog.Debugf("--------passwordHashedBase64=%s", passwordHashedBase64)
	hashed := HashPassword(passwordStr, salt)
	passwordHashed, err := base64.StdEncoding.DecodeString(passwordHashedBase64)
	if err != nil {
		return false
	}
	mlog.Debugf("--------hashed=%s", hex.EncodeToString(hashed))
	mlog.Debugf("--------passwordHashed=%s", string(passwordHashed))
	return hex.EncodeToString(hashed) == string(passwordHashed)
}
