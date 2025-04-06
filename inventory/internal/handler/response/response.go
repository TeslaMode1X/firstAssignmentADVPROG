package handler

import "github.com/gin-gonic/gin"

type BaseResponse struct {
	Message string `json:"message"`
}

func Response(c *gin.Context, status int, message string) {
	c.JSON(status, BaseResponse{Message: message})
}
