package response

import "github.com/gin-gonic/gin"

type BaseResponse struct {
	Op      string      `json:"op,omitempty"`
	Message interface{} `json:"message"`
}

func Response(c *gin.Context, status int, op string, message interface{}) {
	c.JSON(status, BaseResponse{
		Op:      op,
		Message: message,
	})
}
