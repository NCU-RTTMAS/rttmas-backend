package mqtt

import (
	"encoding/json"
	"strings"

	rttmas_binding "rttmas-backend/pkg/binding"
	rttmas_models "rttmas-backend/pkg/models"
	rttmas_service "rttmas-backend/pkg/services"
	"rttmas-backend/pkg/utils/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Default publish handler
// This function is called whenever a MQTT message is received
// Example topic: uplink/user-report/<UID>
func messageHandler(client mqtt.Client, msg mqtt.Message) {

	// if cfg.GetConfigValue("GO_ENV") == "development" {
	// 	logger.Info("Message received from %s:\n%s", msg.Topic(), msg.Payload())
	// } else {
	// 	logger.Info("Message received from %s", msg.Topic())
	// }

	topicParts := strings.Split(msg.Topic(), "/")

	if topicParts[0] != "uplink" {
		return
	}

	var payload map[string]interface{}
	json.Unmarshal(msg.Payload(), &payload)
	logger.Info(payload)

	reportType := topicParts[1]
	reporterUID := topicParts[2]

	switch reportType {
	case "report":
		HandleReport(payload, reporterUID)
	}
}

func HandleReport(payload map[string]interface{}, reporterUID string) {
	reportTime := int64(payload["report_time"].(float64))
	latitude := payload["lat"].(float64)
	longitude := payload["lon"].(float64)
	heading := payload["heading"].(float64)
	speed := payload["speed"].(float64)

	rttmas_binding.RTTMAS_OnUserLocationReport(reportTime, latitude, longitude, reporterUID)

	var mongoReportRecord rttmas_models.UserReport
	mongoReportRecord.Latitude = latitude
	mongoReportRecord.Longitude = longitude
	mongoReportRecord.Speed = speed
	mongoReportRecord.Heading = heading

	rttmas_service.StoreUserReportToMongoDB(reporterUID, reportTime, mongoReportRecord)
	logger.Info("SAVING")
}

// func HandlePlateReport(payload map[string]interface{}, reporterUID string) {
// 	reportTime := int64(payload["report_time"].(float64))
// 	latitude := payload["lat"].(float64)
// 	longitude := payload["lon"].(float64)
// 	reportedPlate := payload["reported_plate"].(string)

// 	rttmas_binding.RTTMAS_OnPlateReport(reportTime, latitude, longitude, reportedPlate, reporterUID)
// }
