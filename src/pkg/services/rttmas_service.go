package services

import (
	"context"
	"fmt"
	rttmas_database "rttmas-backend/pkg/database"
	rttmas_models "rttmas-backend/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
)

func StoreUserReportToMongoDB(UID string, reportTime int64, report rttmas_models.UserReport) {
	userData, err := rttmas_database.MongoFindSingleByField[rttmas_models.UserData](rttmas_database.UserDataCollection, "uid", UID, nil)

	if err != nil {
		userData = rttmas_models.UserData{}
		userData.ID = rttmas_database.GenerateUUIDv7()
		userData.UID = UID
		userData.Reports = make(map[int64]rttmas_models.UserReport)
		userData.Reports[reportTime] = report
		rttmas_database.MongoCreate(rttmas_database.UserDataCollection, userData)
	} else {
		userData.Reports[reportTime] = report
		rttmas_database.MongoUpdate(rttmas_database.UserDataCollection, userData.ID, userData)
	}
}

func StorePlateReportToMongoDB(plateNumber string, reportTime int64, report rttmas_models.PlateReport) {
	plateData, err := rttmas_database.MongoFindSingleByField[rttmas_models.PlateData](rttmas_database.PlateDataCollection, "plate_number", plateNumber, nil)

	if err != nil {
		plateData = rttmas_models.PlateData{}
		plateData.ID = rttmas_database.GenerateUUIDv7()
		plateData.PlateNumber = plateNumber
		plateData.Reports = make(map[int64][]rttmas_models.PlateReport)
		plateData.Reports[reportTime] = append(plateData.Reports[reportTime], report)
		rttmas_database.MongoCreate(rttmas_database.PlateDataCollection, plateData)
	} else {

		filter := bson.M{"id": plateData.ID}
		updateBson := bson.M{
			"$push": bson.M{
				fmt.Sprintf("reports.%d", reportTime): report,
			},
		}
		rttmas_database.PlateDataCollection.UpdateOne(context.TODO(), filter, updateBson)
		// rttmas_database.MongoUpdate(rttmas_database.PlateDataCollection, plateData.ID, updateBson)
	}
}

func QueryUserPath(targetIdentifier string, startTime int64, endTime int64) *map[int64]rttmas_models.UserReport {
	results := make(map[int64]rttmas_models.UserReport)

	userData, err := rttmas_database.MongoFindSingleByField[rttmas_models.UserData](rttmas_database.UserDataCollection, "uid", targetIdentifier, nil)

	if err != nil {
		return nil
	}

	for reportTime, report := range userData.Reports {
		if reportTime >= startTime && reportTime < endTime {
			results[reportTime] = report
		}
	}

	return &results
}

func QueryPlatePath(targetIdentifier string, startTime int64, endTime int64) *map[int64][]rttmas_models.PlateReport {
	results := make(map[int64][]rttmas_models.PlateReport)

	plateData, err := rttmas_database.MongoFindSingleByField[rttmas_models.PlateData](rttmas_database.PlateDataCollection, "plate_number", targetIdentifier, nil)

	if err != nil {
		return nil
	}

	for reportTime, report := range plateData.Reports {
		if reportTime >= startTime && reportTime < endTime {
			results[reportTime] = report
		}
	}

	return &results
}

func QueryAvailableObjects() {

}
