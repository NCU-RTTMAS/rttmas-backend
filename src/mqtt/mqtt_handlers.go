package mqtt

import (
	"encoding/json"
	// "fmt"
	// cfg "rttmas-backend/config"
	"strings"

	// "rttmas-backend/analysis"
	// "rttmas-backend/analysis"
	rttmas_binding "rttmas-backend/binding"
	"rttmas-backend/rttma_simulation"
	"rttmas-backend/statistics"
	"rttmas-backend/utils"
	"rttmas-backend/web/socketio"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Default publish handler
// This function is called whenever a MQTT message is received
// Example topic: uplink/user-report/<UID>
func messageHandler(client mqtt.Client, msg mqtt.Message) {

	// if cfg.GetConfigValue("GO_ENV") == "development" {
	// 	logger.Debug(fmt.Sprintf("Message received from %s:\n%s", string(msg.Topic()), string(msg.Payload())))
	// } else {
	// 	logger.Debug(fmt.Sprintf("Message received from %s", string(msg.Topic())))
	// }

	topicParts := strings.Split(msg.Topic(), "/")

	if topicParts[0] != "traffic" {
		return
	}

	var payload map[string]interface{}
	json.Unmarshal(msg.Payload(), &payload)

	reportType := topicParts[1]
	// reporterUID := topicParts[2]

	switch reportType {
	case "user-report":
		var report rttma_simulation.UserReport
		json.Unmarshal(msg.Payload(), &report)
		rttma_simulation.StoreUserLocationReport(report)
		HandleUserLocationReport(report, report.ReporterUID)
		rttmas_binding.RTTMAS_OnUserLocationReport(int64(report.ReportTime), report.Latitude, report.Longitude, report.ReporterUID)
		socketio.EmitMessage("rttmas", "user-report", utils.Jsonalize(report))
		// parsing plate workflow
		plates := utils.ParseCommaSeparatedString(report.Plates)
		for _, plate := range plates {
			rttma_simulation.StorePlateRecognitionReport(plate, report)
			rttmas_binding.RTTMAS_OnPlateReport(int64(report.ReportTime), report.Latitude, report.Longitude, plate, report.ReporterUID)
			plateBasicInfo := rttma_simulation.GetBasicInfoByPlate(plate)

			// create a struct for plate report
			type PlateReport struct {
				Timestep    int64   `json:"timestep"`
				Latitude    float64 `json:"latitude"`
				Longitude   float64 `json:"longitude"`
				PlateNumber string  `json:"plate_number"`
				SpeedMS     float64 `json:"speed_ms"`
				Heading     float64 `json:"heading"`
			}
			plateReport := PlateReport{
				Timestep:    int64(report.ReportTime),
				Latitude:    report.Latitude,
				Longitude:   report.Longitude,
				PlateNumber: plate,
				SpeedMS:     plateBasicInfo["speed_ms"],
				Heading:     plateBasicInfo["heading"],
			}
			socketio.EmitMessage("rttmas", "plate-report", utils.Jsonalize(plateReport))

		}

		// case "plate-report":
		// 	var prr rttma_simulation.PlateRecognitionReport
		// 	json.Unmarshal(msg.Payload(), &prr)
		// 	rttmas_binding.RTTMAS_OnPlateReport(int64(prr.Timestep), prr.Lat, prr.Lon, prr.PlateNumberSeen, prr.ReporterUID)
		// 	// rttma_simulation.StorePlateRecognitionReport(prr) // WIP
		// 	socketio.EmitMessage("rttmas", "plate-report", utils.Jsonalize(prr))
		// 	// HandlePlateReport(prr)
	}
}

func HandleUserLocationReport(payload rttma_simulation.UserReport, reporterUID string) {
	// reportTime := payload.Timestep
	latitude := payload.Latitude
	longitude := payload.Longitude
	// heading := payload["heading"].(float64)
	// speed := payload["speed"].(float64)

	rttmas_binding.RTTMAS_OnUserLocationReport(int64(payload.ReportTime), latitude, longitude, reporterUID)
	// speed := payload["speed"].(int64)
	// heading := payload["heading"].(int64)
	if payload.SpeedMS != 0 {
		statistics.CollectMapTrafficVectors(int64(payload.ReportTime), payload.Latitude, payload.Longitude, payload.SpeedMS, payload.Heading)
	}
}

func HandlePlateReport(payload rttma_simulation.PlateRecognitionReport) {
	// reportTime := int64(payload["report_time"].(float64))
	// latitude := payload["lat"].(float64)
	// longitude := payload["lon"].(float64)
	// reportedPlate := payload["reported_plate"].(string)
	rttmas_binding.RTTMAS_OnUserLocationReport(int64(payload.Timestep), payload.Lat, payload.Lon, payload.ReporterUID)
}
