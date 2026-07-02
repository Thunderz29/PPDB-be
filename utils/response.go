package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

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
