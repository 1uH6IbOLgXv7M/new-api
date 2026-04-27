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
	// Default to debug mode locally for easier development; set GIN_MODE=release for production.
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
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

	// Validate that port is a reasonable number before starting
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		common.FatalLog(fmt.Sprintf("invalid port number: %s", port))
	}

	common.SysLog(fmt.Sprintf("New API is running on port %s", port))

	if err := server.Run(":" + port); err != nil {
		common.FatalLog("failed to start HTTP server: " + err.Error())
	}
}
