package mqtt

import (
	"encoding/json"
	cfg "rttmas-backend/config"
	"strings"

	rttmas_binding "rttmas-backend/pkg/binding"
	"rttmas-backend/pkg/utils/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Default publish handler
// This function is called whenever a MQTT message is received
// Example topic: uplink/user-report/<UID>
func messageHandler(client mqtt.Client, msg mqtt.Message) {

	if cfg.GetConfigValue("GO_ENV") == "development" {
		logger.Info("Message received from %s:\n%s", msg.Topic(), msg.Payload())
	} else {
		logger.Info("Message received from %s", msg.Topic())
	}

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
	case "user-report":
		HandleUserLocationReport(payload, reporterUID)
	case "plate-report":
		HandlePlateReport(payload, reporterUID)
	}
}

func HandleUserLocationReport(payload map[string]interface{}, reporterUID string) {
	reportTime := int64(payload["report_time"].(float64))
	latitude := payload["lat"].(float64)
	longitude := payload["lon"].(float64)
	// heading := payload["heading"].(float64)
	// speed := payload["speed"].(float64)

	rttmas_binding.RTTMAS_OnUserLocationReport(reportTime, latitude, longitude, reporterUID)
}

func HandlePlateReport(payload map[string]interface{}, reporterUID string) {
	reportTime := int64(payload["report_time"].(float64))
	latitude := payload["lat"].(float64)
	longitude := payload["lon"].(float64)
	reportedPlate := payload["reported_plate"].(string)

	rttmas_binding.RTTMAS_OnPlateReport(reportTime, latitude, longitude, reportedPlate, reporterUID)
}

package mqtt

import (
	"encoding/json"
	cfg "rttmas-backend/config"
	"strings"

	rttmas_binding "rttmas-backend/pkg/binding"
	"rttmas-backend/pkg/utils/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Default publish handler
// This function is called whenever a MQTT message is received
// Example topic: uplink/user-report/<UID>
func messageHandler(client mqtt.Client, msg mqtt.Message) {

	if cfg.GetConfigValue("GO_ENV") == "development" {
		logger.Info("Message received from %s:\n%s", msg.Topic(), msg.Payload())
	} else {
		logger.Info("Message received from %s", msg.Topic())
	}

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
	case "user-report":
		HandleUserLocationReport(payload, reporterUID)
	case "plate-report":
		HandlePlateReport(payload, reporterUID)
	}
}

func HandleUserLocationReport(payload map[string]interface{}, reporterUID string) {
	reportTime := int64(payload["report_time"].(float64))
	latitude := payload["lat"].(float64)
	longitude := payload["lon"].(float64)
	// heading := payload["heading"].(float64)
	// speed := payload["speed"].(float64)

	rttmas_binding.RTTMAS_OnUserLocationReport(reportTime, latitude, longitude, reporterUID)
}

func HandlePlateReport(payload map[string]interface{}, reporterUID string) {
	reportTime := int64(payload["report_time"].(float64))
	latitude := payload["lat"].(float64)
	longitude := payload["lon"].(float64)
	reportedPlate := payload["reported_plate"].(string)

	rttmas_binding.RTTMAS_OnPlateReport(reportTime, latitude, longitude, reportedPlate, reporterUID)
}
