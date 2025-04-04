package main

import (
	"context"
	"os"

	rttmas_analysis "rttmas-backend/analysis"
	rttmas_binding "rttmas-backend/binding"
	rttmas_cfg "rttmas-backend/config"
	"rttmas-backend/cron"
	rttmas_fcm "rttmas-backend/fcm"
	rttmas_models "rttmas-backend/models"
	rttmas_mongo "rttmas-backend/mongo"
	rttmas_mqtt "rttmas-backend/mqtt"
	rttmas_redis "rttmas-backend/redis"

	rttmas_simulation "rttmas-backend/simulation"
	rttmas_web "rttmas-backend/web"

	"rttmas-backend/utils/logger"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

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
	redis := rttmas_redis.GetRedis()
	rttmas_mongo.GetMongoClient()
	rttmas_redis.InitLuaScripts()
	rttmas_mqtt.GetMqttClient()
	// rttmas_mqtt.Init()
	cron.Init()

	redis.FlushAllAsync(context.Background())
}

func initializeMQTT() {
	rttmas_mqtt.GetMqttClient()
}

func initializeWebserver() {
	// Initialize Gin web engine
	webEngine := rttmas_web.GetGinEngine()
	logger.Info("Starting web server on port " + os.Getenv("API_PORT"))

	// Initialize web routes

	// Create server with timeout
	// srv := &http.Server{
	// 	Addr:    ":" + rttmas_cfg.GetConfigValue("API_PORT"),
	// 	Handler: webEngine,
	// 	// Set timeout due CWE-400 - Potential Slowloris Attack
	// 	ReadHeaderTimeout: 5 * time.Second,
	// }

	// // Start server listening
	// if rttmas_cfg.GetConfigValue("USE_TLS") == "true" {
	// 	certfilePath := rttmas_cfg.GetConfigValue("TLS_CERTFILE_PATH")
	// 	keyfilePath := rttmas_cfg.GetConfigValue("TLS_KEYFILE_PATH")
	// 	if err := srv.ListenAndServeTLS(certfilePath, keyfilePath); err != nil {
	// 		logger.Fatal("Failed to start server: %v", err)
	// 	}
	// } else {
	// 	if err := srv.ListenAndServe(); err != nil {
	// 		logger.Fatal("Failed to start server: %v", err)
	// 	}
	// }
	g.Go(func() error {
		if os.Getenv("TLS_ENABLED") == "true" {
			return webEngine.RunTLS(":"+os.Getenv("API_PORT"), os.Getenv("TLS_CERT_DIR")+"/fullchain.pem", os.Getenv("TLS_CERT_DIR")+"/privkey.pem")
		}
		return webEngine.Run(":" + os.Getenv("API_PORT"))
	})
	if err := g.Wait(); err != nil {
		logger.Fatal(err)
	}

}

func initializeRTTMAS() {
	rttmas_binding.RTTMAS_InitializeBindingModule()
}

func initializeFCM() {
	rttmas_fcm.InitializeFCM()
}

func testFunction() {
	r, _ := rttmas_mongo.MongoGetUniqueFieldValues[rttmas_models.UserData](rttmas_mongo.UserDataCollection, "uid", "u__5")
	logger.Info(r)
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
	testFunction()

	initializeRTTMAS()

	if rttmas_cfg.GetConfigValueAsBool("RTTMAS_SIM_ANALYSIS") {
		go rttmas_analysis.StartAnalysisModule()
	}
	if rttmas_cfg.GetConfigValueAsBool("RTTMAS_SIM_BINDING") {
		rttmas_simulation.AnalysisExperiment()
	}

	if rttmas_cfg.GetConfigValueAsBool("RTTMAS_ENABLE_WEBSERVER") {
		initializeWebserver()
	}
}
