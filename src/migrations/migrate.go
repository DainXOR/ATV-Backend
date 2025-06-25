package main

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
)

func init() {
	logger.EnvInit()
	configs.DB.EnvInit()
}

func main() {
	//configs.DataBase.AutoMigrate(&models.UserDB{})
	//configs.DB.Migrate(&models.UserDBGorm{})
}
