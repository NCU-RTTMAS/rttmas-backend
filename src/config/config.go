package config

import (
	"fmt"
	"os"
	"strconv"
)

var config_items map[string]string

func InitializeConfig() {
	config_items = make(map[string]string)

	config_items["RTTMAS_ENABLE_FCM"] = os.Getenv("RTTMAS_ENABLE_FCM")
	config_items["RTTMAS_ENABLE_MQTT"] = os.Getenv("RTTMAS_ENABLE_MQTT")
	config_items["RTTMAS_ENABLE_WEBSERVER"] = os.Getenv("RTTMAS_ENABLE_WEBSERVER")
	config_items["RTTMAS_SIM_ANALYSIS"] = os.Getenv("RTTMAS_SIM_ANALYSIS")
	config_items["RTTMAS_SIM_BINDING"] = os.Getenv("RTTMAS_SIM_BINDING")

	config_items["GO_ENV"] = os.Getenv("GO_ENV")
	if config_items["GO_ENV"] == "development" {
		config_items["GIN_MODE"] = "debug"
	} else {
		config_items["GIN_MODE"] = "release"
	}

	config_items["RATE_LIMITER_REQUEST_PER_SECOND"] = os.Getenv("RATE_LIMITER_REQUEST_PER_SECOND")

	config_items["USE_TLS"] = os.Getenv("USE_TLS")
	config_items["TLS_CERTFILE_PATH"] = os.Getenv("TLS_CERTFILE_PATH")
	config_items["TLS_KEYFILE_PATH"] = os.Getenv("TLS_KEYFILE_PATH")

	config_items["SYSTEM_HOSTNAME"] = os.Getenv("SYSTEM_HOSTNAME")

	if config_items["USE_TLS"] == "true" {
		config_items["API_PORT"] = os.Getenv("API_PORT_TLS")
	} else {
		config_items["API_PORT"] = os.Getenv("API_PORT")
	}

	config_items["CORS_ALLOW_ALL_ORIGINS"] = os.Getenv("CORS_ALLOW_ALL_ORIGINS")

	if config_items["API_PORT"] == "80" || config_items["API_PORT"] == "443" {
		config_items["FULL_API_URL"] = config_items["SYSTEM_HOSTNAME"]
	} else {
		config_items["FULL_API_URL"] = fmt.Sprintf("%s:%s", config_items["SYSTEM_HOSTNAME"], config_items["API_PORT"])
	}

	config_items["JWT_SECRET"] = os.Getenv("JWT_SECRET")

	config_items["MQTT_BROKER_URI"] = fmt.Sprintf("%s://%s:%s", os.Getenv("MQTT_SCHEME"), os.Getenv("MQTT_HOST"), os.Getenv("MQTT_PORT"))
	config_items["MQTT_USERNAME"] = os.Getenv("MQTT_USERNAME")
	config_items["MQTT_PASSWORD"] = os.Getenv("MQTT_PASSWORD")
	config_items["MQTT_QOS"] = os.Getenv("MQTT_QOS")
	config_items["MQTT_SELF_CLIENT_ID"] = os.Getenv("MQTT_SELF_CLIENT_ID")

	config_items["AMQP_BROKER_URI"] = fmt.Sprintf("%s:%s", os.Getenv("AMQP_HOST"), os.Getenv("AMQP_PORT"))
	config_items["AMQP_USERNAME"] = os.Getenv("AMQP_USERNAME")
	config_items["AMQP_PASSWORD"] = os.Getenv("AMQP_PASSWORD")

	config_items["MONGODB_URI"] = fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		os.Getenv("MONGODB_USERNAME"),
		os.Getenv("MONGODB_PASSWORD"),
		os.Getenv("MONGODB_HOST"),
		os.Getenv("MONGODB_PORT"),
	)
}

func GetConfigValue(key string) string {
	val, ok := config_items[key]

	if ok {
		return val
	}
	return ""
}

func GetConfigValueAsInt(key string) int {
	stringVal := GetConfigValue(key)

	intVal, err := strconv.Atoi(stringVal)
	if err != nil {
		intVal = -1
	}

	return int(intVal)
}

func GetConfigValueAsBool(key string) bool {
	stringVal := GetConfigValue(key)
	return stringVal == "true"
}
