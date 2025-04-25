package httpgin

import (
	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup) {
	r.GET("", mainPage()) // Метод для получения главной страницы
}
