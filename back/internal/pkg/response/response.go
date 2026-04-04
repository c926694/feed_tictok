package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "ok",
		Data:    data,
	})
}

func Message(c *gin.Context, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
	})
}

func Fail(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
	})
}
