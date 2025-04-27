package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func generateLink() gin.HandlerFunc {
	return func(c *gin.Context) {
		link := c.Param("link")
		c.JSON(http.StatusOK, link)
	}
}

func redirectLink() gin.HandlerFunc {
	return func(c *gin.Context) {
		link := c.Param("link")
		c.JSON(http.StatusOK, link)
	}
}
