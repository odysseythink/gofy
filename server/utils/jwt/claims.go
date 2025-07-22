package jwtutils

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"mlib.com/mlog"
)

// func ClearToken(c *gin.Context) {
// 	// 增加cookie x-token 向来源的web添加
// 	host, _, err := net.SplitHostPort(c.Request.Host)
// 	if err != nil {
// 		host = c.Request.Host
// 	}

// 	if net.ParseIP(host) != nil {
// 		c.SetCookie("x-token", "", -1, "/", "", false, false)
// 	} else {
// 		c.SetCookie("x-token", "", -1, "/", host, false, false)
// 	}
// }

// func SetToken(c *gin.Context, token string, maxAge int) {
// 	// 增加cookie x-token 向来源的web添加
// 	host, _, err := net.SplitHostPort(c.Request.Host)
// 	if err != nil {
// 		host = c.Request.Host
// 	}

// 	if net.ParseIP(host) != nil {
// 		c.SetCookie("x-token", token, maxAge, "/", "", false, false)
// 	} else {
// 		c.SetCookie("x-token", token, maxAge, "/", host, false, false)
// 	}
// }

func GetToken(c *gin.Context) string {
	// token, _ := c.Cookie("x-token")
	// if token == "" {
	// 	token = c.Request.Header.Get("x-token")
	// }
	// return token

	authHeader := c.Request.Header.Get("Authorization")
	if authHeader != "" {
		// 假设token紧跟在Bearer后面，用空格分隔
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 {
			if strings.ToLower(tokenParts[0]) != "bearer" {
				return ""
			}
			token := tokenParts[1]

			return token
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func GetClaims(c *gin.Context, secretKey string) (*CustomClaims, error) {
	token := GetToken(c)
	claims, err := ParseToken(token, secretKey)
	if err != nil {
		mlog.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context, secretKey string) (string, error) {
	cclaims, err := GetClaims(c, secretKey)
	if err != nil {
		return "", err
	}
	return cclaims.UserID, nil
}

// Custom claims structure
type CustomClaims struct {
	UserID string `json:"user_id,omitempty"`
	jwt.RegisteredClaims
}

func GenToken(user_id, secretKey, iss string, ep time.Duration) (token string, claims CustomClaims, err error) {
	claims = CreateClaims(user_id, iss, ep)
	token, err = CreateToken(claims, secretKey)
	if err != nil {
		return
	}
	return
}

func GetAccountJWTToken(user_id, secretKey, iss string, ep time.Duration) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"dify"},           // 受众
			NotBefore: jwt.NewNumericDate(now.Add(-1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(now.Add(ep)),    // 过期时间 7天  配置文件
			Issuer:    iss,                                // 签名的发行者
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
