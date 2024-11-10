package web

import (
	"net/http"
	"rttmas-backend/pkg/rttma_simulation"
	// "rttmas-backend/pkg/utils/logger"

	"github.com/gin-gonic/gin"
)

func HomepageHandler(c *gin.Context) {
	// router := gin.Default()
	// router.GET("/", func(c *gin.Context) {
	c.HTML(http.StatusOK, "show-map.tmpl", gin.H{})
	// })
}

// func VehicleRouteHandler(c *gin.Context) {
// 	// router := gin.Default()
// 	// router.GET("/", func(c *gin.Context) {
// 	Tracks := GetTracksByPlateNumber(c.Param("plate"))
// 	logger.Info(c.Param("plate"))
// 	// c.HTML(http.StatusOK, "vehicle-route.tmpl", gin.H{
// 	// 	"PlateNumber": c.Param("plate"),
// 	// 	"Tracks":      Tracks,
// 	// })
// 	c.JSON(http.StatusOK, gin.H{
// 		"PlateNumber": c.Param("plate"),
// 		"Tracks":      Tracks,
// 	})
// 	// })
// }

// http://localhost:3000/vehicles?plate=ACL-7727&plate=AKA-2910&plate=AYI-1461
// func VehiclesRouteHandler(c *gin.Context) {
// 	// router := gin.Default()
// 	// router.GET("/", func(c *gin.Context) {
// 	// logger.Info(c.QueryArray("plate[]"))
// 	plates := c.QueryArray("plate")
// 	TrackCollections := [][]models.Track_t{}
// 	for _, plate := range plates {
// 		TrackCollections = append(TrackCollections, GetTracksByPlateNumber(plate))
// 	}
// 	logger.Info(c.Param("plate"))
// 	c.HTML(http.StatusOK, "vehicles-route.tmpl", gin.H{
// 		"PlateNumber":      plates,
// 		"TrackCollections": TrackCollections,
// 	})
// 	// })
// }

// func FileHandler(c *gin.Context) {
// 	c.File()
// }

func GetUserHandler(c *gin.Context) {
	result := rttma_simulation.GetAllUsers()
	c.JSON(http.StatusOK, result)
}

func GetUserByUIDHandler(c *gin.Context) {
	result := rttma_simulation.GetUserByUID(c.Param("uid"))
	c.JSON(http.StatusOK, result)
}
func GetVehiclesHandler(c *gin.Context) {
	result := rttma_simulation.GetAllVehicles()
	c.JSON(http.StatusOK, result)
}
func GetSingleVehicleHandler(c *gin.Context) {
	result := rttma_simulation.GetTracksByPlateNumber(c.Param("plate"))
	c.JSON(http.StatusOK, result)
}
