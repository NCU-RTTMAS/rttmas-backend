package analysis

import (
	// "context"
	"context"
	"encoding/json"
	"math"
	"sort"
	"strconv"

	// "encoding/json"
	"fmt"
	rttmas_db "rttmas-backend/pkg/database"
	"rttmas-backend/pkg/mqtt"
	"rttmas-backend/pkg/utils"
	"rttmas-backend/pkg/utils/logger"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type AnalysisResult struct {
	Heading        int64
	Velocity       int64
	LatestTimestep int64
}

func StartAnalysisModule() {
	ctx := context.Background()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	logger.Info("Analysis Module Started")
	redisClient := rttmas_db.GetRedis()

	for {
		select {
		case <-ticker.C:
			err := processAnalysis(ctx, redisClient)
			if err != nil {
				logger.Error("Error during analysis process:", err)
			}
		}
	}
}

func processAnalysis(ctx context.Context, redisClient *redis.Client) error {
	users, err := getAllSortedSetKeys(ctx, redisClient, "plate_locations:*")
	if err != nil {
		return fmt.Errorf("failed to get sorted set keys: %w", err)
	}

	for _, userWithPrefix := range users {
		err := processUser(ctx, redisClient, userWithPrefix)
		if err != nil {
			logger.Error("Error processing user:", err)
		}
	}

	dangerCars, err := FetchTopKeysByVelocity(ctx, redisClient, "basic_info:")
	if err != nil {
		return fmt.Errorf("failed to fetch top danger cars: %w", err)
	}
	if len(dangerCars) != 0 {
		mqtt.PublishMessageToTopic("alarm/speeding", utils.Jsonalize(dangerCars))
		for plate := range dangerCars {
			latest_timestep, err := rttmas_db.GetRedis().JSONGet(context.Background(), "basic_info:"+plate, "$").Result()
			// if err != nil {
			// 	logger.Info(latest_timestep)
			// }

			var vehInfo []AnalysisResult
			err = json.Unmarshal([]byte(latest_timestep), &vehInfo)
			if err != nil {
				logger.Error(err)
			}
			if len(vehInfo) != 0 {
				logger.Info("p_locations:" + strconv.Itoa(int(vehInfo[0].LatestTimestep)))
				result, err := rttmas_db.GetRedis().GeoPos(context.Background(), "p_locations:"+string(vehInfo[0].LatestTimestep), plate).Result()
				if err != nil {
					logger.Error(err)
				}
				logger.Info(result)
			}
			// logger.Info(vehInfo)

			// latest_timestep = latest_timestep.([]AnalysisResult)
			// logger.Info(vehInfo[0])
			// logger.Info("loc:", result)
		}
	}

	return nil
}

func processUser(ctx context.Context, redisClient *redis.Client, userWithPrefix string) error {
	user, _ := strings.CutPrefix(userWithPrefix, "plate_locations:")

	currentTime := redisClient.ZRevRange(ctx, userWithPrefix, 0, 1).Val()
	if len(currentTime) != 2 {
		return nil // Insufficient data, skip
	}

	heading, err := rttmas_db.RedisExecuteLuaScript("calculate_heading", []string{userWithPrefix}, currentTime[0], currentTime[1])
	if err != nil {
		return fmt.Errorf("error calculating heading: %w", err)
	}

	velocity, err := rttmas_db.RedisExecuteLuaScript("get_velocity", []string{userWithPrefix}, 5)
	if err != nil && err.Error() != redis.Nil.Error() {
		return fmt.Errorf("error getting velocity: %w", err)
	}

	if heading != nil && velocity != nil {
		result := AnalysisResult{
			Heading:  heading.(int64),
			Velocity: velocity.(int64),
		}

		if velocity != nil && velocity.(int64) > 40 {
			logger.Warning(user, "is speeding")
		}

		err = updateUserInfo(ctx, redisClient, user, result)
		if err != nil {
			return fmt.Errorf("error updating user info: %w", err)
		}

	}

	return nil
}

func updateUserInfo(ctx context.Context, redisClient *redis.Client, user string, result AnalysisResult) error {
	basicInfoKey := fmt.Sprintf("basic_info:%s", user)

	_, err := redisClient.JSONSet(ctx, basicInfoKey, "$.Velocity", result.Velocity).Result()
	if err != nil {
		return fmt.Errorf("error setting velocity: %w", err)
	}

	_, err = redisClient.JSONSet(ctx, basicInfoKey, "$.Heading", result.Heading).Result()
	if err != nil {
		return fmt.Errorf("error setting heading: %w", err)
	}

	err = redisClient.Expire(ctx, basicInfoKey, 10*time.Second).Err()
	if err != nil {
		return fmt.Errorf("error setting expiration: %w", err)
	}

	return nil
}
func getAllSortedSetKeys(ctx context.Context, rdb *redis.Client, pattern string) ([]string, error) {
	var cursor uint64
	var keys []string
	for {
		var err error
		var batchKeys []string

		// SCAN with count 0 lets Redis decide the batch size, and TYPE zset filters for sorted sets
		batchKeys, cursor, err = rdb.ScanType(ctx, cursor, pattern, 0, "zset").Result()
		if err != nil {
			return nil, err
		}

		// Append the found keys to the result list
		keys = append(keys, batchKeys...)

		// Exit when the cursor is 0, meaning the scan is complete
		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

type BasicInfo struct {
	LatestTimestep int `json:"LatestTimestep"`
	Velocity       int `json:"Velocity"`
	Heading        int `json:"Heading"`
}

func FetchTopKeysByVelocity(ctx context.Context, rdb *redis.Client, keyPrefix string) (map[string]int, error) {
	// Fetch all keys matching the prefix
	keys, err := rdb.Keys(ctx, keyPrefix+"*").Result()
	if err != nil {
		return nil, fmt.Errorf("error fetching keys: %w", err)
	}

	// Map to hold key and velocity
	data := make([]struct {
		Plate    string
		Velocity int
	}, 0)

	// Iterate over keys and fetch their values
	for _, key := range keys {
		val, err := rdb.JSONGet(ctx, key).Result() // Fetch the JSON string for the key
		if err != nil {
			logger.Error(fmt.Sprintf("Error fetching key %s: %v", key, err))
			continue
		}

		var info BasicInfo
		if err := json.Unmarshal([]byte(val), &info); err != nil {
			logger.Error(fmt.Sprintf("Error unmarshalling data for key %s: %v", key, err))
			continue
		}

		// Extract the car plate (strip the prefix from the key)
		plate := key[len(keyPrefix):]

		data = append(data, struct {
			Plate    string
			Velocity int
		}{Plate: plate, Velocity: info.Velocity})
	}

	// Sort data by velocity in descending order
	sort.Slice(data, func(i, j int) bool {
		return data[i].Velocity > data[j].Velocity
	})

	// Calculate the top 10%
	topCount := int(math.Ceil(float64(len(data)) * 0.1))

	// Prepare the result as a slice of maps
	// var topEntries map[string]int
	topEntries := make(map[string]int, topCount)
	// topEntries := make(map[string]int, 0, topCount)
	for _, entry := range data[:topCount] {
		if entry.Velocity < 20 {
			continue
		}
		// topEntries = append(topEntries, map[string]int{entry.Plate: entry.Velocity})
		topEntries[entry.Plate] = entry.Velocity

	}

	return topEntries, nil
}
