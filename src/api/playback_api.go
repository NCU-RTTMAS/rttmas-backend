package api

import (
	"net/http"
	"strconv"

	rttmas_models "rttmas-backend/models"
	rttmas_mongo "rttmas-backend/mongo"
	rttmas_service "rttmas-backend/services"

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

	var results interface{}
	if objectType == 1 {
		results = rttmas_service.QueryUserPath(targetIdentifier, startTime, endTime)
	} else if objectType == 2 {
		results = rttmas_service.QueryPlatePath(targetIdentifier, startTime, endTime)
	}

	c.JSON(http.StatusOK, results)
}

func QueryAvailableObjects(c *gin.Context) {
	searchQuery := c.Query("query")

	allUIDs, _ := rttmas_mongo.MongoGetUniqueFieldValues[rttmas_models.UserData](rttmas_mongo.UserDataCollection, "uid", searchQuery)
	allPlateNumbers, _ := rttmas_mongo.MongoGetUniqueFieldValues[rttmas_models.PlateData](rttmas_mongo.PlateDataCollection, "plate_number", searchQuery)

	result := map[string][]interface{}{
		"uids":          allUIDs,
		"plate_numbers": allPlateNumbers,
	}

	c.JSON(http.StatusOK, result)
}
