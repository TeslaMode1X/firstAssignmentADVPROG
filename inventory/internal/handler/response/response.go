package handler

import "github.com/gin-gonic/gin"

type BaseResponse struct {
	Message interface{} `json:"message"`
}

func Response(c *gin.Context, status int, message interface{}) {
	c.JSON(status, BaseResponse{Message: message})
}
