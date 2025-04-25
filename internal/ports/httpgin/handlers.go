package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func mainPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "test")
	}
}
