package main

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
)

func init() {
	logger.EnvInit()
	configs.DB.EnvInit()
}

func main() {
	//configs.DataBase.AutoMigrate(&models.UserDB{})
	configs.DB.Migrate(&models.UserDBGorm{})
}
