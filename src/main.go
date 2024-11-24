package main

import (
	"context"
	"net/http"
	"time"

	rttmas_cfg "rttmas-backend/config"
	rttmas_binding "rttmas-backend/pkg/binding"
	rttmas_db "rttmas-backend/pkg/database"
	rttmas_fcm "rttmas-backend/pkg/fcm"
	rttmas_mqtt "rttmas-backend/pkg/mqtt"

	rttmas_simulation "rttmas-backend/pkg/simulation"
	rttmas_web "rttmas-backend/pkg/web"

	"rttmas-backend/pkg/utils/logger"

	"github.com/joho/godotenv"
)

func initializeConfig() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			logger.Fatal("Error loading .env file")
		}
	}
	rttmas_cfg.InitializeConfig()
}

func initializeDatabase() {
	redis := rttmas_db.GetRedis()
	rttmas_db.GetMongoClient()
	rttmas_db.InitLuaScripts()

	redis.FlushAllAsync(context.Background())
}

func initializeMQTT() {
	rttmas_mqtt.GetMqttClient()
}

func initializeWebserver() {
	// Initialize Gin web engine
	webEngine := rttmas_web.GetGinEngine()

	// Create server with timeout
	srv := &http.Server{
		Addr:    ":" + rttmas_cfg.GetConfigValue("API_PORT"),
		Handler: webEngine,
		// Set timeout due CWE-400 - Potential Slowloris Attack
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start server listening
	if rttmas_cfg.GetConfigValue("USE_TLS") == "true" {
		certfilePath := rttmas_cfg.GetConfigValue("TLS_CERTFILE_PATH")
		keyfilePath := rttmas_cfg.GetConfigValue("TLS_KEYFILE_PATH")
		if err := srv.ListenAndServeTLS(certfilePath, keyfilePath); err != nil {
			logger.Fatal("Failed to start server: %v", err)
		}
	} else {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal("Failed to start server: %v", err)
		}
	}
}

func initializeRTTMAS() {
	rttmas_binding.RTTMAS_InitializeBindingModule()
}

func initializeFCM() {
	rttmas_fcm.InitializeFCM()
}

func main() {

	initializeConfig()
	initializeDatabase()

	if rttmas_cfg.GetConfigValueAsBool("RTTMAS_ENABLE_MQTT") {
		initializeMQTT()
	}
	if rttmas_cfg.GetConfigValueAsBool("RTTMAS_ENABLE_FCM") {
		initializeFCM()
	}

	initializeRTTMAS()

	if rttmas_cfg.GetConfigValueAsBool("RTTMAS_ENABLE_WEBSERVER") {
		initializeWebserver()
	}

	if rttmas_cfg.GetConfigValueAsBool("RTTMAS_SIM_ANALYSIS") {
		// go rttmas_analysis.StartAnalysisModule()
	}
	if rttmas_cfg.GetConfigValueAsBool("RTTMAS_SIM_BINDING") {
		rttmas_simulation.AnalysisExperiment()
	}
}
