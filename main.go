package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"new-api/common"
	"new-api/middleware"
	"new-api/model"
	"new-api/router"
)

func main() {
	common.SetupLogger()
	common.SysLog("New API starting...")

	// Initialize database
	err := model.InitDB()
	if err != nil {
		common.FatalLog("failed to initialize database: " + err.Error())
	}
	defer func() {
		err := model.CloseDB()
		if err != nil {
			common.SysError("failed to close database: " + err.Error())
		}
	}()

	// Initialize Redis if configured
	err = common.InitRedisClient()
	if err != nil {
		common.SysError("failed to initialize Redis: " + err.Error())
	}

	// Initialize options from database
	model.InitOptionMap()

	// Set Gin mode based on environment.
	// Defaulting to release mode here to keep logs clean; use GIN_MODE=debug locally when needed.
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(middleware.RequestId())
	middleware.SetUpLogger(server)

	// Setup all routes
	router.SetRouter(server)

	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(*common.Port)
	}

	common.SysLog(fmt.Sprintf("New API is running on port %s", port))

	if err := server.Run(":" + port); err != nil {
		common.FatalLog("failed to start HTTP server: " + err.Error())
	}
}
