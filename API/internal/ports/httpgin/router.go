package httpgin

import (
	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup) {
	r.GET("/shorter/:link", generateLink()) // Метод для получения главной страницы
	r.GET("/:link", redirectLink())         // Метод для получения главной страницы
}
