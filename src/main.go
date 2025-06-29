package main

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/middleware"
	"dainxor/atv/routes"

	"cmp"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Move this number so the deprecations are in sync with the API version

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Warning("Error loading .env file: " + err.Error())
	}

	configs.App.EnvInit() // Initialize application configurations

	logger.EnvInit()
	logger.SetAppVersion(configs.App.Version())

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
	router.Use(middleware.RecoverMiddleware()) // Middleware to recover from panics and logs a small trace
	router.Use(middleware.CORSMiddleware())
	//router.Use(middleware.TokenMiddleware())

	// Root level routes
	routes.MainRoutes(router)

	// Api routes
	routes.InfoRoutes(router, configs.App.ApiVersion(), configs.App.RoutesVersion()) // Routes for information about the API
	routes.TestRoutes(router)                                                        // Routes for testing purposes

	// Versioned API routes
	routes.StudentRoutes(router)
	routes.UniversityRoutes(router)
	routes.SpecialityRoutes(router)
	routes.CompanionRoutes(router)
	routes.SessionTypeRoutes(router)
	routes.SessionRoutes(router)

	router.Run(address()) // listen and serve on 0.0.0.0:8080 (for windows ":8080")
}
