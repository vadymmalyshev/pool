package apierrors

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime"
)

const defaultApiErrorMessage = "Bad request"

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func NewApiErr(code int, message string) *ApiError {
	_, fn, line, _ := runtime.Caller(1)
	logrus.Errorf("%s:%d %v", fn, line, message)
	return &ApiError{
		Code:    code,
		Message: message,
	}
}

func NewApiErrorBadId(id string) *ApiError {
	message := fmt.Sprintf("Invalid ID:%s", id)
	_, fn, line, _ := runtime.Caller(1)
	logrus.Errorf("%s:%d %v", fn, line, message)
	return &ApiError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewDefApiError() *ApiError {
	return &ApiError{
		Code:    http.StatusBadRequest,
		Message: defaultApiErrorMessage,
	}
}

// Handle error and log line num and func name
func HandleError(err error) (b bool) {
	if err != nil {
		// 1 - log where the error happened
		_, fn, line, _ := runtime.Caller(1)
		logrus.Errorf("%s:%d %v", fn, line, err)
		b = true
	}
	return
}

// Handles error in api-handler, makes JSON abort if error exist, return true if aborted
func AbortWithApiError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	_, fn, line, _ := runtime.Caller(1)
	switch e := err.(type) {
	case *ApiError:
		logrus.Warnf("%s:%d aborted: %v", fn, line, e)
		c.AbortWithStatusJSON(e.Code, e)
		return true
	default:
		def := NewDefApiError()
		logrus.Warnf("%s:%d aborted: %v", fn, line, def)
		c.AbortWithStatusJSON(def.Code, def)
		return true
	}
}
