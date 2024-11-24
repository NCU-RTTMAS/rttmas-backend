package database

import (
	"context"
	"fmt"
	"rttmas-backend/pkg/utils/logger"
	"sync"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisClient *redis.Client

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "sOmE_sEcUrE_pAsS", // no password set
		DB:       0,                  // use default DB
	})
}

func GetRedis() *redis.Client {
	if redisClient == nil {
		initRedis()
	}
	return redisClient
}

func EnableKeyExpiredNotification() {
	// this is telling redis to publish events since it's off by default.
	_, err := GetRedis().Do(ctx, "CONFIG", "SET", "notify-keyspace-events", "KEA").Result()
	if err == nil {
		// this is telling redis to subscribe to events published in the keyevent channel,
		// specifically for expired events
		pubsub := GetRedis().PSubscribe(ctx, "__keyevent@0__:expired")
		wg := &sync.WaitGroup{}
		go func(redis.PubSub) {
			for {
				_, err := pubsub.ReceiveMessage(ctx)
				if err != nil {
					logger.Info(fmt.Sprintf("[PubSub] Error message : %v\n", err.Error()))
					break
				}
				// Handle the logic of received messages below
				// logger.Info(fmt.Sprintf("[PubSub] Object Expired : %s", msg.Payload))
			}
		}(*pubsub)
		wg.Wait()
	} else {
		fmt.Printf("Unable to set keyspace events : %v\n", err.Error())
	}
}
