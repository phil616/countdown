package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Errors  []FieldError `json:"errors,omitempty"`
}

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

const (
	CodeSuccess             = 0
	CodeValidationFailed    = 40001
	CodeTokenMissing        = 40101
	CodeTokenInvalid        = 40102
	CodeCredentialError     = 40103
	CodeLoginLocked         = 40104
	CodeNotFound            = 40401
	CodeBusinessError       = 42201
	CodeInternalServerError = 50001
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
		Meta:    meta,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func BadRequest(c *gin.Context, msg string, errors []FieldError) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeValidationFailed,
		Message: msg,
		Errors:  errors,
	})
}

func Unauthorized(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    code,
		Message: msg,
	})
}

func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: msg,
	})
}

func TooManyRequests(c *gin.Context, msg string) {
	c.JSON(http.StatusTooManyRequests, Response{
		Code:    CodeLoginLocked,
		Message: msg,
	})
}

func BusinessError(c *gin.Context, msg string) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Code:    CodeBusinessError,
		Message: msg,
	})
}

func InternalError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeInternalServerError,
		Message: msg,
	})
}

func CalculateTotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	return pages
}
