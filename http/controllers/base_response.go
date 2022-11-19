package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func NewErrorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, BaseResponse{
		Status:  false,
		Message: err,
	})
}
