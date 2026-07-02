package utils

import (
	"github.com/gin-gonic/gin"
)

// Response represents a standard API response structure for success
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents a standard API response structure for failures
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SendSuccess sends a standardized success response
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError sends a standardized error response
func SendError(c *gin.Context, statusCode int, message string, errs interface{}) {
	var detail interface{}
	if err, ok := errs.(error); ok {
		detail = err.Error()
	} else {
		detail = errs
	}

	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  detail,
	})
}
