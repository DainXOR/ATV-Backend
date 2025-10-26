package service

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	var _ = godotenv.Overload("../.env", "../.env.test")
	logger.SetVersion(configs.App.ApiVersion())
	if configs.DB.Use(configs.GetMongoAccessor()).Start() != nil {
		logger.Fatal("Failed to connect to the database")
	}
	// configs.DB.Migrate(&models.StudentDBMongo{})
	logger.Info("Env configurations loaded")
	logger.Debug("Starting server")
}

func TestCompanionService_Create(t *testing.T)     {}
func TestCompanionService_GetByID(t *testing.T)    {}
func TestCompanionService_GetAll(t *testing.T)     {}
func TestCompanionService_UpdateByID(t *testing.T) {}
func TestCompanionService_PatchByID(t *testing.T)  {}
func TestCompanionService_DeleteByID(t *testing.T) {}
