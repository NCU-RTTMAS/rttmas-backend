package analysis

import (
	// "context"
	"context"
	// "encoding/json"
	"fmt"
	rttmas_db "rttmas-backend/pkg/database"
	"rttmas-backend/pkg/utils/logger"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func StartAnalysisModule() {
	// ctx := context.Background()
	ticker := time.NewTicker(1000 * time.Millisecond)
	logger.Info("Analysis Module Started")
	defer ticker.Stop()
	count := 0
	for {
		select {
		case <-ticker.C:
			count += 1
			// Execute your Redis command here
			// GetRedis().ScanType(context.Background(), 0, "user_location_report:*", 0, "zset").Result()
			users, err := getAllSortedSetKeys(context.Background(), rttmas_db.GetRedis(), "user_location_report:*")
			if err != nil {
				logger.Info(err)
			}
			for _, userWithPrefix := range users {
				user, _ := strings.CutPrefix(userWithPrefix, "user_location_report:")

				current_time := rttmas_db.GetRedis().ZRevRange(context.Background(), userWithPrefix, 0, 1)
				// logger.Info(current_time.)
				if len(current_time.Val()) != 2 {
					continue
				}
				// logger.Info(current_time.Val())

				result, err := rttmas_db.RedisExecuteLuaScript("calculate_heading", []string{userWithPrefix}, current_time.Val()[0], current_time.Val()[1])
				if err != nil {
					logger.Info(err)
				}
				// logger.Info(result)
				type foo struct {
					Heading  interface{}
					Velocity interface{} // float64
				}
				// var bar foo
				var bar foo
				bar.Heading = result

				result, err = rttmas_db.RedisExecuteLuaScript("get_velocity", []string{userWithPrefix}, 10)
				if err != nil && err.Error() != redis.Nil.Error() {
					logger.Info(err)
					continue
				}
				bar.Velocity = result

				// bazz, err := json.Marshal(bar)
				// if err != nil {
				// 	logger.Error(err)
				// }
				if result != nil && result.(int64) > 60 {
					logger.Warning(user, " is speeding")
				}

				_, err = rttmas_db.GetRedis().JSONSet(context.Background(), fmt.Sprintf("basic_info:%s", user), "$.Velocity", bar.Velocity).Result()
				if err != nil {
					logger.Error(err)
				}
				_, err = rttmas_db.GetRedis().JSONSet(context.Background(), fmt.Sprintf("basic_info:%s", user), "$.Heading", bar.Heading).Result()
				if err != nil {
					logger.Error(err)
				}
				rttmas_db.GetRedis().Expire(context.Background(), fmt.Sprintf("basic_info:%s", user), 10*time.Second)
				result, err = rttmas_db.RedisExecuteLuaScript("search_nearby_cars", []string{"basic_info:*"}, userWithPrefix)
				if err != nil && err.Error() != redis.Nil.Error() {
					logger.Info(err)
					continue
				}
				// logger.Info(result)
			}

			// err := GetRedis().Set(ctx, "key", count, 0).Err()
			if err != nil {
				logger.Error("Error executing Redis command:", err)
			}
		}
	}

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
