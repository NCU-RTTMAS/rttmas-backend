package statistics

import (
	// "math"

	"context"
	"encoding/json"
	"fmt"
	rttmas_redis "rttmas-backend/redis"
	"rttmas-backend/utils/logger"

	// "rttmas-backend/utils/logger"

	"github.com/mmcloughlin/geohash"
)

// rttmas_db "rttmas-backend/database"

// "github.com/redis/go-redis/v9"

type Direction int

const (
	N   Direction = iota // North
	NNE                  // North-Northeast
	NE                   // Northeast
	ENE                  // East-Northeast
	E                    // East
	ESE                  // East-Southeast
	SE                   // Southeast
	SSE                  // South-Southeast
	S                    // South
	SSW                  // South-Southwest
	SW                   // Southwest
	WSW                  // West-Southwest
	W                    // West
	WNW                  // West-Northwest
	NW                   // Northwest
	NNW                  // North-Northwest
)

type TrafficVector struct {
	Speed   float64   `json:"speed"`
	Heading Direction `json:"heading"`
}

func classifyDirection(heading float64) Direction {
	index := int((heading+11.25)/22.5) % 16
	return Direction(index)
}

type MapChunk struct {
	// Geohash         string             `json:"geohash"`          // Unique identifier for the map chunk (e.g., geohash of 6 digits)
	// VehicleCount    int                `json:"vehicle_count"`    // Total number of vehicles observed
	VelocityByDirection map[Direction]DirectionInfo `json:"velocity_by_direction"` // Sum of the speeds of all vehicles
	// Vectors         map[Direction]float64 `json:"vehicles"`         // Optional: Map of vehicle IDs to their last observed speed
	LastUpdated int64 `json:"last_updated"` // Last update timestamp (UNIX time)
}
type DirectionInfo struct {
	AverageSpeed float64   `json:"average_speed"`
	Records      []float64 `json:"records"`
}

// var ChunksCollection map[string]*MapChunk

func init() {

}
func CollectMapTrafficVectors(reportTime int64, latitude float64, longitude float64, speed float64, heading float64) {
	hash := geohash.Encode(latitude, longitude)[:7]
	var currentChunk MapChunk

	// Fetch the current value from the database
	val, err := rttmas_redis.GetRedis().JSONGet(context.Background(), fmt.Sprintf("map_chunks:%s", hash), "$").Result()
	if err == nil && val != "" {
		// If value exists, unmarshal it into currentChunk
		var chunks []MapChunk
		err = json.Unmarshal([]byte(val), &chunks)
		if err != nil {
			fmt.Println("Error unmarshalling data:", err)
			return
		}
		if len(chunks) > 0 {
			currentChunk = chunks[0]
		}
	} else {
		// If value does not exist, initialize a new MapChunk
		currentChunk = MapChunk{
			VelocityByDirection: make(map[Direction]DirectionInfo),
			LastUpdated:         reportTime,
		}
	}

	// currentChunk.VelocityByDirection[classifyDirection(heading)] = append(currentChunk.VelocityByDirection[classifyDirection(heading)], speed)
	direction := classifyDirection(heading)
	directionInfo := currentChunk.VelocityByDirection[direction]
	directionInfo.Records = append(directionInfo.Records, speed)
	totalSpeed := 0.0
	for _, s := range directionInfo.Records {
		totalSpeed += s
	}
	directionInfo.AverageSpeed = totalSpeed / float64(len(directionInfo.Records))
	currentChunk.VelocityByDirection[direction] = directionInfo

	currentChunk.LastUpdated = reportTime

	// Store the updated chunk back to the database
	rttmas_redis.GetRedis().JSONSet(context.Background(), fmt.Sprintf("map_chunks:%s", hash), "$", currentChunk)
}

func GetAverageSpeed(geohash string) map[Direction]DirectionInfo {
	result, err := rttmas_redis.GetRedis().JSONGet(context.Background(), fmt.Sprintf("map_chunks:%s", geohash), "$.velocity_by_direction").Result()
	if err != nil || result == "" {
		logger.Error(err)
		return map[Direction]DirectionInfo{}
	}
	speedMap := []map[Direction]DirectionInfo{}
	err = json.Unmarshal([]byte(result), &speedMap)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(speedMap)

	return speedMap[0]
}

func PruneMapChunks() {

}

func storeMapChunks() {

}

func UpdateMapChunk() {

}
