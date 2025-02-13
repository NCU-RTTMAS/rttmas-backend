package api

import (
	"net/http"
	"rttmas-backend/pkg/statistics"

	"github.com/gin-gonic/gin"
)

func GetAverageSpeed(c *gin.Context) {
	geohash := c.Request.URL.Query().Get("geohash")
	result := statistics.GetAverageSpeed(geohash)
	c.JSON(http.StatusOK, result)
}
