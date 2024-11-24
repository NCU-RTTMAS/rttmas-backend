package api

import (
	"net/http"
	"strconv"

	rttmas_service "rttmas-backend/pkg/services"
	"rttmas-backend/pkg/utils/logger"

	"github.com/gin-gonic/gin"
)

func String2Int64(input string) int64 {
	i, err := strconv.Atoi(input)
	if err != nil {
		return 0
	}
	return int64(i)
}

func QueryObjectPath(c *gin.Context) {
	objectType := String2Int64(c.Query("object_type"))
	targetIdentifier := c.Query("target_identifier")
	startTime := String2Int64(c.Query("start_time"))
	endTime := String2Int64(c.Query("end_time"))

	results := rttmas_service.QueryObjectPath(objectType, targetIdentifier, startTime, endTime)
	logger.Info(results)

	c.JSON(http.StatusOK, results)
}
