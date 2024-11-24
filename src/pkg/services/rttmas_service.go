package services

import (
	rttmas_database "rttmas-backend/pkg/database"
	rttmas_models "rttmas-backend/pkg/models"
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

func QueryObjectPath(objectType int64, targetIdentifier string, startTime int64, endTime int64) *map[int64]rttmas_models.UserReport {
	results := make(map[int64]rttmas_models.UserReport)

	if objectType == 1 {
		userData, err := rttmas_database.MongoFindSingleByField[rttmas_models.UserData](rttmas_database.UserDataCollection, "uid", targetIdentifier, nil)

		if err != nil {
			return nil
		}

		for reportTime, report := range userData.Reports {
			if reportTime >= startTime && reportTime < endTime {
				results[reportTime] = report
			}
		}
	}

	return &results
}
