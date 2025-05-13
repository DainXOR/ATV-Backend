package main

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/middleware"
	"dainxor/atv/models"
	"dainxor/atv/routes"
	"os"

	"github.com/gin-gonic/gin"
	//"github.com/joho/godotenv"
)

func init() {
	logger.Init()

	//err := godotenv.Load()
	//if err != nil {
	//	logger.Error("Error loading .env file")
	//}

	logger.EnvInit()
	configs.DB.EnvInit()
	logger.Info("Env configurations loaded")
	logger.Info("Starting server")

	configs.DB.Migrate(&models.UserDB{})

}

func address() string {
	address, exist := os.LookupEnv("SERVER_ADDRESS")
	if !exist {
		logger.Warning("SERVER_ADDRESS not found, using default")
		address = ":8080"
	}

	return address
}

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	routes.MainRoutes(router)
	routes.InfoRoutes(router)
	routes.TestRoutes(router)

	routes.UserRoutes(router)

	router.Run(address()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
