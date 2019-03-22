package common

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type errorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ErrorResponse aborts request with specific error
func ErrorResponse(c *gin.Context, err error) {
	status := http.StatusBadRequest

	if c.Writer.Status() != http.StatusOK {
		return
	}

	c.AbortWithStatusJSON(status, errorResponse{
		Code:    status,
		Message: err.Error(),
		Error:   errors.Cause(err).Error(),
	})
}
