package database

import (
	"context"
	// "fmt"

	// "go.mongodb.org/mongo-driver/bson"
	// "rttma-backend/pkg/models"
	"rttmas-backend/config"
	"rttmas-backend/pkg/utils/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoCli *mongo.Client

type mongoCollection struct {
	*mongo.Collection
}
type Collections struct {
	Vehicles                *mongo.Collection
	UserReports             *mongo.Collection
	PlateRecognitionReports *mongo.Collection
	VehicleTrueLocations    *mongo.Collection
	UserLocationReports     *mongo.Collection
	// WIP: other collections
}

var RTTMA_Collections Collections

// var RTTMA_Database = mongoCli.Database("rttma")
var RTTMA_Database *mongo.Database

func initDatatables() {
	RTTMA_Database = mongoCli.Database("rttma")

	RTTMA_Collections.Vehicles = RTTMA_Database.Collection("vehicles")
	RTTMA_Collections.PlateRecognitionReports = RTTMA_Database.Collection("plate-recognition-reports")
	RTTMA_Collections.UserReports = RTTMA_Database.Collection("user-reports")
	RTTMA_Collections.VehicleTrueLocations = RTTMA_Database.Collection("vehicle-true-locations")
	RTTMA_Collections.UserLocationReports = RTTMA_Database.Collection("user-location-reports")

}

func initEngine() {
	var err error
	dbString := config.GetConfigValue("MONGODB_URI")
	// dbString := fmt.Sprintf("mongodb://%s:%s@%s:%s", databaseInfo.Username, databaseInfo.Password, databaseInfo.Server.Host, databaseInfo.Server.Port)
	clientOptions := options.Client().ApplyURI(dbString)

	mongoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Fatal(err)
	}

	err = mongoCli.Ping(context.TODO(), nil)
	if err != nil {
		logger.Fatal(err)
	}

}

func GetMongo() *mongo.Client {
	if mongoCli == nil {
		initEngine()
		initDatatables()
	}
	return mongoCli
}

// func (coll *mongoCollection) StoreVehicle(v models.Vehicle_t) {
// 	logger.Info(coll)
// 	result, err := coll.Collection.InsertOne(context.Background(), v)
// 	if err != nil {
// 		logger.Error(err)
// 	}
// 	logger.Info(result)
// }

// func initDatabase() {

// 	if result, err := GetClient().Database("mlt").Collection("users").CountDocuments(context.Background(), bson.M{}); err != nil {
// 		logger.Fatal("database initalization failed")
// 	} else if result == 0 {
// 		GetClient().Database("mlt").Collection("users").InsertOne(context.Background(), models.DefaultUser)
// 	}

// 	if result, err := GetClient().Database("mlt").Collection("global_config").CountDocuments(context.Background(), bson.M{}); err != nil {
// 		logger.Fatal("database initalization failed")
// 	} else if result == 0 {
// 		GetClient().Database("mlt").Collection("global_config").InsertOne(context.Background(), models.SMTPConfigTemplate)
// 		GetClient().Database("mlt").Collection("global_config").InsertOne(context.Background(), models.PeriodicPollingPeriodTemplate)
// 		GetClient().Database("mlt").Collection("global_config").InsertOne(context.Background(), models.SystemTimezoneTemplate)
// 	}
// }

/*
 */
