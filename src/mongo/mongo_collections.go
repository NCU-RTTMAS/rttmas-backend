package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var Database *mongo.Database

var UserDataCollection *mongo.Collection
var PlateDataCollection *mongo.Collection

// Initialize the database and collections
// This is called only upon MongoDB initialization
// and should not be called by external modules
func initializeDatabaseAndCollections(mongoClient *mongo.Client) {
	Database = mongoClient.Database("rttmas")

	UserDataCollection = Database.Collection("user_data")
	PlateDataCollection = Database.Collection("plate_data")
}
