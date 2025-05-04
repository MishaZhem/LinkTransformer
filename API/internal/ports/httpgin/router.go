package httpgin

import (
	"LinkTransformer/internal/app"

	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.GET("/shorter/:link", generateLink(a)) // Метод для получения главной страницы
	r.GET("/:link", redirectLink(a))         // Метод для получения главной страницы
}
