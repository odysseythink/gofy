package jwtutils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/odysseythink/mlog"
)

// type JWT struct {
// 	SigningKey []byte
// }

var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("that's not even a token")
	ErrTokenInvalid     = errors.New("couldn't handle this token")
)

// func NewJWT() *JWT {
// 	return &JWT{
// 		[]byte(global.GVA_CONFIG.JWT.SigningKey),
// 	}
// }

func CreateClaims(userID, iss string, ep time.Duration) CustomClaims {
	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"gofy"},           // 受众
			NotBefore: jwt.NewNumericDate(now.Add(-1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(now.Add(ep)),    // 过期时间 7天  配置文件
			Issuer:    iss,                                // 签名的发行者
		},
	}
	return claims
}

// 创建一个token
func CreateToken(claims CustomClaims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// // CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
// func (j *JWT) CreateTokenByOldToken(oldToken string, claims request.CustomClaims) (string, error) {
// 	v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (any, error) {
// 		return j.CreateToken(claims)
// 	})
// 	return v.(string), err
// }

// 解析 token
func ParseToken(tokenString, secretKey string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i any, e error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				mlog.Error("pase jwt with claims failed:", err)
				return nil, ErrTokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, ErrTokenInvalid

	} else {
		return nil, ErrTokenInvalid
	}
}
