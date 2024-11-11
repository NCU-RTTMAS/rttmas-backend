package web

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	cfg "rttmas-backend/config"
	"rttmas-backend/pkg/utils/logger"
)

var ginEngine *gin.Engine
var once sync.Once

func GetGinEngine() *gin.Engine {
	once.Do(func() {
		gin.SetMode(cfg.GetConfigValue("GIN_MODE"))

		ginEngine = gin.Default()

		// Add security headers to protect the API
		ginEngine.Use(func(c *gin.Context) {
			if c.Request.Host != cfg.GetConfigValue("FULL_API_URL") {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
				return
			}
			c.Header("X-Frame-Options", "DENY")
			c.Header("Access-Control-Allow-Origin", "*")
			// c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
			c.Header("X-XSS-Protection", "1; mode=block")
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			c.Header("Referrer-Policy", "strict-origin")
			c.Header("X-Content-Type-Options", "nosniff")
			c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
			c.Next()
		})
		ginEngine.GET("/users", GetUserHandler)
		ginEngine.GET("/user/:uid", GetUserByUIDHandler)
		ginEngine.GET("/vehicle/:plate", GetSingleVehicleHandler)
		ginEngine.GET("/vehicles", GetVehiclesHandler)

		logger.Info("Gin web server initialization complete.")
	})
	return ginEngine
}
