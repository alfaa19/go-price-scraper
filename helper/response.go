package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	successJson struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func ResponseSuccessJson(c *gin.Context, message string, data interface{}) {

	if message == "" {
		message = "success"
	}

	res := successJson{
		Message: message,
		Success: true,
		Data:    data,
	}

	c.JSON(http.StatusOK, res)
}
