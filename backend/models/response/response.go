package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    string `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

const (
	ERROR = 7
	// UNAUTHORIED = 1
	SUCCESS = 0
)

func Result(status int, data any, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		Status:  status,
		Data:    data,
		Message: msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]any{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]any{}, message, c)
}

func OkWithData(data any, c *gin.Context) {
	Result(SUCCESS, data, "成功", c)
}

func OkWithDetailed(data any, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]any{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]any{}, message, c)
}

func NoAuth(message string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Status:  7,
		Data:    nil,
		Message: message,
	})
}

func FailWithDetailed(data any, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}
