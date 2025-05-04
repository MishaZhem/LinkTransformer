package httpgin

import (
	"github.com/gin-gonic/gin"
)

type linkResponse struct {
	Link string `json:"link"`
}

func LinkSuccessResponse(link string) gin.H {
	return gin.H{
		"data": linkResponse{
			Link: link,
		},
		"error": nil,
	}
}

func ErrorResponse(err error) gin.H {
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
