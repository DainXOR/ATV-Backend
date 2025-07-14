package main

import (
	"cmp"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	//"github.com/joho/godotenv"

	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/middleware"
	"dainxor/atv/routes"
)

//var envErr = godotenv.Load()

func init() {
	//if envErr != nil {
	//	logger.Warning("Error loading .env file: " + envErr.Error())
	//}
	//if logger.UsingDefaults() {
	//	logger.ReloadEnv()
	//	configs.ReloadAppEnv()
	//	configs.ReloadDBEnv()
	//}

	logger.SetAppVersion(configs.App.ApiVersion())

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
	defer configs.DB.Close()

	router := gin.Default()
	router.Use(middleware.RecoverMiddleware()) // Middleware to recover from panics and logs a small trace
	router.Use(middleware.CORSMiddleware())
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
}
