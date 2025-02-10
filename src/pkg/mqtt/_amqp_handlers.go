package mqtt

// import (
// 	"encoding/json"

	rttmas_binding "rttmas-backend/pkg/binding"
	"rttmas-backend/pkg/rttma_simulation"
	"rttmas-backend/pkg/utils"
	"rttmas-backend/pkg/utils/logger"
	"rttmas-backend/pkg/web/socketio"

// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// func CreateAMQPEndpoints() {
// 	err := CreateConsumer("analysis_module_queue", func(msg amqp.Delivery) {
// 		// logger.Info(msg.RoutingKey, "\n", string(msg.Body))
// 		// PublishToTopic("foo", string(msg.Body))
// 		// Add your message processing logic here

		switch msg.RoutingKey {
		case "traffic.plate_recognition_reports":
			var prr rttma_simulation.PlateRecognitionReport
			json.Unmarshal(msg.Body, &prr)
			socketio.EmitMessage("", "plate-report", utils.Jsonalize(prr))
			rttmas_binding.RTTMAS_OnPlateReport(int64(prr.Timestep), prr.Lat, prr.Lon, prr.PlateNumberSeen, prr.ReporterUID)
			rttma_simulation.StorePlateRecognitionReport(prr)
		case "traffic.user_location_reports":
			var ulr rttma_simulation.UserLocationReport
			json.Unmarshal(msg.Body, &ulr)
			rttma_simulation.StoreUserLocationReport(ulr)
			HandleUserLocationReport(ulr, ulr.UID)
			rttmas_binding.RTTMAS_OnUserLocationReport(int64(ulr.Timestep), ulr.Lat, ulr.Lon, ulr.UID)
			socketio.EmitMessage("", "user-report", utils.Jsonalize(ulr))
		case "traffic.vehicle_true_locations":
			var vtl rttma_simulation.VehicleTrueLocation
			json.Unmarshal(msg.Body, &vtl)
			rttma_simulation.StoreVehicleTrueLocation(vtl)
		}
	})
	if err != nil {
		logger.Info(err)
	}

// }
