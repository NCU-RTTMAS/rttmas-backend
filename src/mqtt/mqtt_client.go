package mqtt

import (
	cfg "rttmas-backend/config"
	"rttmas-backend/utils/logger"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client
var once sync.Once

// Singleton MQTT instance getter
// This initializes the MQTT client for the first time
func GetMqttClient() mqtt.Client {
	once.Do(func() {
		opts := mqtt.NewClientOptions()

		opts.AddBroker(cfg.GetConfigValue("MQTT_BROKER_URI"))
		opts.SetUsername(cfg.GetConfigValue("MQTT_USERNAME"))
		opts.SetPassword(cfg.GetConfigValue("MQTT_PASSWORD"))
		opts.SetClientID(cfg.GetConfigValue("MQTT_SELF_CLIENT_ID"))

		opts.OnConnect = handleMqttOnConnect
		opts.OnConnectionLost = handleMqttOnConnectionLost

		opts.SetDefaultPublishHandler(messageHandler)

		mqttClient = mqtt.NewClient(opts)

		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			logger.Error(token.Error())
			return
		}

		logger.Info("MQTT client initialization complete.")
	})

	return mqttClient
}

// Default handler for mqtt connect events
func handleMqttOnConnect(client mqtt.Client) {
	client.Subscribe("traffic/#", byte(cfg.GetConfigValueAsInt("MQTT_QOS")), nil)
}

// Default handler for mqtt connection lost events
func handleMqttOnConnectionLost(client mqtt.Client, err error) {

}

func PublishMessageToTopic(topic string, message string) {
	token := mqttClient.Publish(topic, byte(cfg.GetConfigValueAsInt("MQTT_QOS")), false, message)
	token.Wait()
}
