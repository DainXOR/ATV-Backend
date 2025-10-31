package service

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Overload("../.env", "../.env.test")

	logger.ReloadEnv()
	configs.App.ReloadEnv()
	logger.SetVersion(configs.App.ApiVersion())

	if err := configs.DB.ReloadConnection(); err != nil {
		logger.Fatal("Failed to reload DB config from environment:", err)
	}

	m.Run()
}

func TestCompanionService_Create(t *testing.T) {
}
func TestCompanionService_GetByID(t *testing.T)    {}
func TestCompanionService_GetAll(t *testing.T)     {}
func TestCompanionService_UpdateByID(t *testing.T) {}
func TestCompanionService_PatchByID(t *testing.T)  {}
func TestCompanionService_DeleteByID(t *testing.T) {}
