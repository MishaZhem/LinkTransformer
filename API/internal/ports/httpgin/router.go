package httpgin

import (
	"LinkTransformer/internal/app"

	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.GET("/shorter/:link", generateLink(a))
	r.GET("/:link", redirectLink(a))
	r.GET("/stats/:link", getStatistic(a))
	r.GET("/clicks/:link", getTotalClicks(a))
}
