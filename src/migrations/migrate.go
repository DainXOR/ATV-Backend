package main

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"

	"github.com/joho/godotenv"
)

func init() {
	logger.Init()
	logger.Info("Loading configurations")

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")

	}

	logger.EnvInit()
	configs.DB.EnvInit()
}

func main() {
	configs.DataBase.AutoMigrate(&models.UserDB{})
}
