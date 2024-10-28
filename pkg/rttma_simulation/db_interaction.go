package rttma_simulation

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"

	"rttmas-backend/pkg/database"
	// "rttma-backend/pkg/utils"
	"rttmas-backend/pkg/utils/logger"
)

func GetAllUsers() interface{} {
	filter := bson.A{
		bson.D{{"$group", bson.D{{"_id", "$uid"}}}},
		bson.D{{"$sort", bson.D{{"_id", 1}}}},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"uid", "$_id"},
				},
			},
		},
	}
	var v []interface{}
	cursor, err := database.RTTMA_Collections.UserLocationReports.Aggregate(context.Background(), filter)
	if err != nil {
		logger.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			logger.Fatal(err)
		}
		v = append(v, result)
	}
	return v
}

func GetUserByUID(uid string) interface{} {
	filter := bson.A{

		bson.D{{"$match", bson.D{{"uid", uid}}}},
		bson.D{{"$sort", bson.D{{"timestep", 1}}}},
	}
	var v []interface{}
	cursor, err := database.RTTMA_Collections.UserLocationReports.Aggregate(context.Background(), filter)
	if err != nil {
		logger.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			logger.Fatal(err)
		}
		v = append(v, result)
	}
	return v
}
func GetAllVehicleIDs() interface{} {
	filter := bson.A{
		bson.D{{"$group", bson.D{{"_id", "$vid"}}}},
		bson.D{{"$sort", bson.D{{"_id", 1}}}},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"vid", "$_id"},
				},
			},
		},
	}
	var v []interface{}
	cursor, err := database.RTTMA_Collections.VehicleTrueLocations.Aggregate(context.Background(), filter)
	if err != nil {
		logger.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			logger.Fatal(err)
		}
		v = append(v, result)
	}
	return v
}
func StorePlateRecognitionReport(prr PlateRecognitionReport) error {
	_, err := database.RTTMA_Collections.PlateRecognitionReports.InsertOne(context.Background(), prr)
	// database.GetRedis().SAdd(context.Background(), fmt.Sprintf("plate_recognition_report:%d", prr.Timestep), utils.JsonToString(prr))
	database.GetRedis().GeoAdd(context.Background(), fmt.Sprintf("plate_recognition_report:%d", prr.Timestep), &redis.GeoLocation{
		Name:      prr.PlateNumberSeen,
		Latitude:  prr.Lat,
		Longitude: prr.Lon,
	})
	database.GetRedis().Expire(context.Background(), fmt.Sprintf("plate_recognition_report:%d", prr.Timestep), 10*time.Second)

	if err != nil {
		return fmt.Errorf("failed to insert PlateRecognitionReport: %v", err)
	}
	return nil
}

func StoreUserLocationReport(ulr UserLocationReport) error {
	_, err := database.RTTMA_Collections.UserLocationReports.InsertOne(context.Background(), ulr)
	database.GetRedis().GeoAdd(context.Background(), fmt.Sprintf("user_location_report:%s", ulr.UID), &redis.GeoLocation{
		Name:      fmt.Sprintf("%d", ulr.Timestep),
		Latitude:  ulr.Lat,
		Longitude: ulr.Lon,
	})
	exists, err := database.GetRedis().Exists(context.Background(), fmt.Sprintf("basic_info:%s", ulr.UID)).Result()
	if err != nil {
		// handle error
	}
	if exists == 0 {
		initialData := `{}`
		_, err = database.GetRedis().JSONSet(context.Background(), fmt.Sprintf("basic_info:%s", ulr.UID), "$", initialData).Result()
		if err != nil {
			// handle error
		}
	}

	database.GetRedis().JSONSet(context.Background(), fmt.Sprintf("basic_info:%s", ulr.UID), "$.LatestTimestep", fmt.Sprintf("%d", ulr.Timestep))

	database.GetRedis().Expire(context.Background(), fmt.Sprintf("user_location_report:%s", ulr.UID), 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to insert UserLocationReport: %v", err)
	}
	return nil
}

func StoreVehicleTrueLocation(vtl VehicleTrueLocation) error {
	_, err := database.RTTMA_Collections.VehicleTrueLocations.InsertOne(context.Background(), vtl)
	database.GetRedis().GeoAdd(context.Background(), fmt.Sprintf("vehicle_true_location:%s", vtl.VID), &redis.GeoLocation{
		Name:      fmt.Sprintf("%d", vtl.Timestep),
		Latitude:  vtl.Lat,
		Longitude: vtl.Lon,
	})
	database.GetRedis().Expire(context.Background(), fmt.Sprintf("vehicle_true_location:%s", vtl.VID), 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to insert VehicleTrueLocation: %v", err)
	}
	return nil
}
