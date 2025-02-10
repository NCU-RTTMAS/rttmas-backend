package cron

import (
	"fmt"

	"github.com/robfig/cron/v3"

	// "omc-backend/pkg/models"
	// "omc-backend/pkg/services/database"
	// "omc-backend/pkg/services/messaging/mqtt"
	// logger "omc-backend/pkg/utils/logger"
	"rttmas-backend/pkg/database"
	logger "rttmas-backend/pkg/utils/logger"
)

var cr *cron.Cron

func Init() {
	cr = GetCronInstance()
	MapChunkAutoPrune(cr)
}

func GetCronInstance() *cron.Cron {
	if cr == nil {
		cr = cron.New()
		cr.Start()
	}
	return cr
}

// validates crontab strings
func ValidateCron(cronString string) error {
	_, cronerr := cron.ParseStandard(cronString)
	if cronerr != nil {
		logger.Info()
		return cronerr
	}
	return nil
}

// // create cronjob service by given cron string and device IDs
func MapChunkAutoPrune(cr *cron.Cron) {
	cronString := "*/15 * * * *"
	id, err := cr.AddFunc(cronString, func() {
		database.DeleteKeysWithPrefix("map_chunks")
		logger.Info(fmt.Sprintf("Pruning Map Chunks"))
	})
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info(fmt.Sprintf("Periodic Pruning in crontab of %s, next execution time :%s", cronString, cr.Entry(id).Next))
	}

}

// // Global config updater
// func UpdateCronConfig() {
// 	var result models.SystemConfig
// 	database.SystemConfigDb.FindOne(context.Background(), bson.D{{"key", "PERIODIC_POLLING_PERIOD"}}).Decode(&result)
// 	RemoveCron(cr)
// 	AddAutoPolling(cr, fmt.Sprintf("0 */%s * * *", result.Value))
// }

// remove cronjob service
func RemoveCron(cr *cron.Cron) {
	// logger.Info(cr.Entries()[0])
	logger.Info(cr.Entries())
	if len(cr.Entries()) != 0 {
		cr.Remove(cr.Entries()[0].ID)
	}

}
