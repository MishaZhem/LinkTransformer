package httpgin

import (
	"LinkTransformer/internal/app"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func generateLink(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		link := c.Param("link")

		url, err := a.GenerateLink(context.Background(), link)
		if err != nil {
			getStatusByError(c, err)
			return
		}

		c.JSON(http.StatusOK, LinkSuccessResponse(url))
	}
}

func redirectLink(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		link := c.Param("link")

		url, err := a.RedirectLink(context.Background(), link)
		if err != nil {
			getStatusByError(c, err)
			return
		}

		c.JSON(http.StatusOK, LinkSuccessResponse(url))
	}
}

func getStatistic(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		link := c.Param("link")

		events, err := a.GetStatistics(context.Background(), link)
		if err != nil {
			getStatusByError(c, err)
			return
		}

		c.JSON(http.StatusOK, ClickEventSuccessResponse(events))
	}
}

func getTotalClicks(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		link := c.Param("link")

		clicks, err := a.GetTotalClicks(context.Background(), link)
		if err != nil {
			getStatusByError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  clicks,
			"error": nil,
		})
	}
}

func getStatusByError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, app.ErrForbidden):
		c.JSON(http.StatusForbidden, ErrorResponse(err))
	case errors.Is(err, app.ErrBadRequest):
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}
