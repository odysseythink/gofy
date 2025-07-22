package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	"mlib.com/gofy/server/utils"
	"mlib.com/mlog"
)

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					mlog.Error(c.Request.URL.Path,
						"error", err,
						"request", string(httpRequest))
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					mlog.Error("[Recovery from panic]",
						"error", err,
						"request", string(httpRequest),
						"stack", string(debug.Stack()),
					)
				} else {
					mlog.Error("[Recovery from panic]",
						"error", err,
						"request", string(httpRequest),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// ErrorMiddleware is a Gin middleware to handle errors.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				status_code := 0
				var default_data map[string]any
				if real_exp, ok := r.(httpexceptions.HTTPException); ok {
					real_exp.Response(c)
					return
				} else if real_exp, ok := r.(*exceptions.ValueError); ok {
					status_code = http.StatusBadRequest
					default_data = map[string]any{
						"code":    "invalid_param",
						"message": real_exp.Error(),
						"status":  status_code,
					}
				} else if real_exp, ok := r.(*exceptions.AppInvokeQuotaExceededError); ok {
					status_code = http.StatusTooManyRequests
					default_data = map[string]any{
						"code":    "too_many_requests",
						"message": real_exp.Error(),
						"status":  status_code,
					}
				} else {
					mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
					status_code = http.StatusInternalServerError
					default_data = map[string]any{
						"code":    "invalid_param",
						"message": http.StatusText(status_code),
					}
				}

				if status_code == http.StatusBadRequest {
					c.JSON(http.StatusBadRequest, default_data)
				} else {
					c.JSON(status_code, default_data)
				}
				c.Abort()
			}
		}()
		c.Next()

	}
}
