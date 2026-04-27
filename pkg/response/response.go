package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Code      Code   `json:"code"`                 // 0 表示成功
	Message   string `json:"message"`              // 提示信息
	Data      T      `json:"data"`                 // 返回数据
	RequestID string `json:"request_id,omitempty"` // 可选，链路追踪
}

// JSON 成功返回（无 msg，默认 "success"）
func SuccessNoMsg[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Response[T]{
		Code:      SuccessCode,
		Message:   "success", // 可以改成你系统统一默认成功提示
		Data:      data,
		RequestID: getRequestID(c),
	})
}

// JSON 成功返回
func Success[T any](c *gin.Context, data T, msg string) {
	c.JSON(http.StatusOK, Response[T]{
		Code:      SuccessCode,
		Message:   msg,
		Data:      data,
		RequestID: getRequestID(c),
	})
}

// JSON 失败返回
func Fail(c *gin.Context, code Code, msg string) {
	c.JSON(http.StatusOK, Response[any]{
		Code:      code,
		Message:   msg,
		Data:      nil,
		RequestID: getRequestID(c),
	})
}

// FailAbort 失败返回并中止请求
func FailAbort(c *gin.Context, code Code, msg string) {
	Fail(c, code, msg)
	c.Abort()
}

// getRequestID 从 Context 获取 request_id
func getRequestID(c *gin.Context) string {
	if rid, exists := c.Get("request_id"); exists {
		return rid.(string)
	}
	return ""
}
