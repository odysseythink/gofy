package middleware

import (
	"github.com/gin-gonic/gin"
)

// 处理跨域请求,支持options访问
func NeedInit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// 处理请求
	}
}
