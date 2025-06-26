package main

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/middleware"
	"dainxor/atv/routes"
	"strconv"

	"cmp"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	DEFAULT_ROUTE_VERSION = 1       // Default version for the API routes
	DEFAULT_API_VERSION   = "0.1.1" // Default version for the API
)

// Move this number so the deprecations are in sync with the API version

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Warning("Error loading .env file: " + err.Error())
	}

	envVersion, _ := strconv.ParseUint(os.Getenv("ATV_ROUTE_VERSION"), 10, 32)
	programVersion := uint(cmp.Or(envVersion, DEFAULT_ROUTE_VERSION))

	logger.EnvInit()
	logger.SetAppVersion(programVersion)

	configs.DB.EnvInit()
	// configs.DB.Migrate(&models.StudentDBMongo{})
	logger.Info("Env configurations loaded")
	logger.Debug("Starting server")

}

// address returns the server address from the environment variable
func address() string {
	envAddress := os.Getenv("SERVER_ADDRESS")
	return cmp.Or(envAddress, ":8080")
}

func main() {
	router := gin.Default()
	router.Use(middleware.RecoverMiddleware()) // Middleware to recover from panics and log errors
	router.Use(middleware.CORSMiddleware())

	// Root level routes
	routes.MainRoutes(router)

	// Api routes
	routes.InfoRoutes(router, DEFAULT_API_VERSION, logger.AppVersion()) // Routes for information about the API
	routes.TestRoutes(router)                                           // Routes for testing purposes

	// Versioned API routes
	routes.StudentRoutes(router)    // Routes for user management
	routes.UniversityRoutes(router) // Routes for university management
	routes.SpecialityRoutes(router) // Routes for speciality management
	routes.CompanionRoutes(router)  // Routes for companion management
	routes.SessionRoutes(router)    // Routes for session management

	router.Run(address()) // listen and serve on 0.0.0.0:8080 (for windows ":8080")
}
