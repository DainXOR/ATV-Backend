package main

import (
	"cmp"
	"dainxor/atv/configs"
	"dainxor/atv/db"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Warning("Error loading .env file: " + err.Error())
	}

	envVersion, _ := strconv.ParseUint(os.Getenv("ATV_ROUTE_VERSION"), 10, 32)
	programVersion := uint64(cmp.Or(envVersion, 1))

	logger.EnvInit()
	logger.SetAppVersion(programVersion)

	configs.DB.EnvInit()
	// configs.DB.Migrate(&models.StudentDBMongo{})
	logger.Info("Env configurations loaded")
	logger.Debug("Starting server")

}

func TestStudentOperations(t *testing.T) {
	createObj := models.StudentCreate{
		NumberID:         "123456789",
		FirstName:        "John",
		LastName:         "Doe",
		PersonalEmail:    "john.doe@example.com",
		InstitutionEmail: "john.doe@university.edu",
		ResidenceAddress: "123 University St, City, Country",
		Semester:         1,
		IDUniversity:     "685c180f0d2362de34ec5721", // Example ObjectID
		PhoneNumber:      "123-456-7890",
	}

	resultObj := db.Student.Create(createObj)

	if resultObj.IsErr() {
		t.Errorf("Failed to create student: %v", resultObj.Error())
		return
	}

	getResult := db.Student.GetByID(resultObj.Value().ID.Hex())

	patchObg := models.StudentCreate{
		NumberID:    "1234567890",
		FirstName:   "Johnny",
		Semester:    2,
		PhoneNumber: "1234567891",
	}

	patchResult := db.Student.PatchByID(getResult.Value().ID.Hex(), patchObg)
	if patchResult.IsErr() {
		t.Errorf("Failed to patch student: %v", patchResult.Error())
		return
	}

}
