package web

import (
	// "net/http"
	"embed"
	"io/fs"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	rttmas_api "rttmas-backend/api"
	cfg "rttmas-backend/config"
	"rttmas-backend/utils/logger"
	"rttmas-backend/web/socketio"
)

var ginEngine *gin.Engine
var once sync.Once

//go:embed dist/*
var embeddedFiles embed.FS

func GetGinEngine() *gin.Engine {
	once.Do(func() {
		// gin.SetMode(cfg.GetConfigValue("GIN_MODE"))
		gin.SetMode("debug")

		ginEngine = gin.Default()

		// Add security headers to protect the API

		ginEngine.Use(func(c *gin.Context) {
			if c.Request.Host != cfg.GetConfigValue("FULL_API_URL") {
				// logger.Info(c.Request.Host)
				// c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
				return
			}
			c.Header("X-Frame-Options", "DENY")
			c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
			// c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
			c.Header("X-XSS-Protection", "1; mode=block")
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			c.Header("Referrer-Policy", "strict-origin")
			c.Header("X-Content-Type-Options", "nosniff")
			c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
			c.Next()
		})
		// ginEngine.GET("/users", GetUserHandler)
		// ginEngine.GET("/user/:uid", GetUserByUIDHandler)
		// ginEngine.GET("/vehicle/:plate", GetSingleVehicleHandler)
		// ginEngine.GET("/vehicles", GetVehiclesHandler)
		setAPIRoutes(ginEngine)
		ginEngine.GET("/socket.io/*any", gin.WrapH(socketio.GetServerInstance()))
		ginEngine.POST("/socket.io/*any", gin.WrapH(socketio.GetServerInstance()))
		logger.Info("Gin web server initialization complete.")
		// ginEngine.Use(SPAMiddleware("/", embeddedFiles))
		SPAHandler().Register(ginEngine)
		// ginEngine.NoRoute(func(c *gin.Context) {
		// 	distFS, err := DistFS()
		// 	if err != nil {
		// 		logger.Error("Error accessing embedded files:", err)
		// 		c.Status(http.StatusInternalServerError)
		// 		return
		// 	}
		// 	filePath := c.Request.URL.Path
		// 	_, err = distFS.Open(filePath)
		// 	if err != nil {
		// 		// If the requested file is not found, serve the main page (index.html)
		// 		filePath = "index.html"
		// 	}
		// 	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		// 	c.Header("Pragma", "no-cache")
		// 	c.Header("Expires", "0")
		// 	http.StripPrefix("/", http.FileServer(distFS)).ServeHTTP(c.Writer, c.Request)
		// })
	})
	return ginEngine
}
func DistFS() (http.FileSystem, error) {
	subFS, err := fs.Sub(embeddedFiles, "dist")
	if err != nil {
		return nil, err
	}
	return http.FS(subFS), nil
}

func setAPIRoutes(ginEngine *gin.Engine) {
	// ginEngine.GET("/users", GetUserHandler)
	// ginEngine.GET("/user/:uid", GetUserByUIDHandler)
	// ginEngine.GET("/vehicle/:plate", GetSingleVehicleHandler)
	// ginEngine.GET("/vehicles", GetVehiclesHandler)

	apis := ginEngine.Group("/api/v1")
	apis.GET("/velocity", rttmas_api.GetAverageSpeed)
	apis.GET("/playback", rttmas_api.QueryObjectPath)
	apis.GET("/playback/objects", rttmas_api.QueryAvailableObjects)
}
