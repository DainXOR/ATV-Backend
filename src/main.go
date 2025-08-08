package main

import (
	"cmp"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"dainxor/atv/configs"
	"dainxor/atv/logger"
)

//var envErr = godotenv.Load()

func init() {
	logger.SetVersion(configs.App.ApiVersion())
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

	record := logger.NewRecord("Server is starting")
	record.LogLevel = logger.Level.Info()

	/*
		defer configs.DB.Close()

		router := gin.New()
		router.Use(gin.Logger())
		router.Use(gin.Recovery())
		router.Use(middleware.Recovery()) // Middleware to recover from panics and logs a small trace
		router.Use(middleware.CORS())
		//router.Use(middleware.TokenMiddleware())

		// Root level routes
		routes.MainRoutes(router)

		// Api routes
		routes.InfoRoutes(router) // Routes for information about the API
		routes.TestRoutes(router) // Routes for testing purposes

		// Versioned API routes
		routes.StudentRoutes(router)
		routes.UniversityRoutes(router)
		routes.SpecialityRoutes(router)
		routes.CompanionRoutes(router)
		routes.SessionTypeRoutes(router)
		routes.SessionRoutes(router)

		router.Run(address()) // listen and serve on 0.0.0.0:8080 (for windows ":8080")
	*/
}
