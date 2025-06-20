package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorType 错误类型
type ErrorType string

const (
	ErrInvalidRequest ErrorType = "INVALID_REQUEST"
	ErrUnauthorized   ErrorType = "UNAUTHORIZED"
	ErrForbidden      ErrorType = "FORBIDDEN"
	ErrNotFound       ErrorType = "NOT_FOUND"
	ErrDatabase       ErrorType = "DATABASE_ERROR"
	ErrInternal       ErrorType = "INTERNAL_ERROR"
	ErrFileOperation  ErrorType = "FILE_OPERATION_ERROR"
)

// Error 自定义错误类型
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error 实现error接口
func (e *Error) Error() string {
	return e.Message
}

// NewError 创建新的错误
func NewError(message string, code int) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Write 将错误响应写入HTTP响应
func (e *ErrorResponse) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	return json.NewEncoder(w).Encode(e)
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(err error) (int, *ErrorResponse) {
	if err == nil {
		return http.StatusInternalServerError, &ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	if e, ok := err.(*Error); ok {
		return e.Code, &ErrorResponse{
			Code:    e.Code,
			Message: e.Message,
		}
	}

	// 默认返回500错误
	return http.StatusInternalServerError, &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}

// WrapError 包装错误
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return &Error{
			Message: fmt.Sprintf("%s: %s", message, e.Message),
			Code:    e.Code,
		}
	}

	return &Error{
		Message: fmt.Sprintf("%s: %s", message, err.Error()),
		Code:    http.StatusInternalServerError,
	}
}

// getStatusCode 获取HTTP状态码
func getStatusCode(errType ErrorType) int {
	switch errType {
	case ErrInvalidRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrDatabase:
		return http.StatusInternalServerError
	case ErrInternal:
		return http.StatusInternalServerError
	case ErrFileOperation:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// IsInvalidRequest 检查是否是无效请求错误
func IsInvalidRequest(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == http.StatusBadRequest
	}
	return false
}

// IsUnauthorized 检查是否是未授权错误
func IsUnauthorized(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == http.StatusUnauthorized
	}
	return false
}

// IsForbidden 检查是否是禁止访问错误
func IsForbidden(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == http.StatusForbidden
	}
	return false
}

// IsNotFound 检查是否是未找到错误
func IsNotFound(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == http.StatusNotFound
	}
	return false
}

// IsDatabase 检查是否是数据库错误
func IsDatabase(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == http.StatusInternalServerError && e.Message == "database error"
	}
	return false
}

// IsInternal 检查是否是内部服务器错误
func IsInternal(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == http.StatusInternalServerError
	}
	return false
}

// IsFileOperation 检查是否是文件操作错误
func IsFileOperation(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == http.StatusInternalServerError && e.Message == "file operation error"
	}
	return false
} 