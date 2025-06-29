package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type sessionType struct{}

var Session sessionType

func (sessionType) Create(u models.SessionCreate) types.Result[models.SessionDBMongo] {
	logger.Debug("Creating session with data: ", u)
	studentResult := Student.GetByID(u.IDStudent)
	if studentResult.IsErr() {
		httpErr := studentResult.Error().(*types.HttpError)
		logger.Warning("Failed to get student by ID: ", httpErr)
		return types.ResultErr[models.SessionDBMongo](httpErr)
	}

	companionResult := Companion.GetByID(u.IDCompanion)
	if companionResult.IsErr() {
		httpErr := companionResult.Error().(*types.HttpError)
		logger.Warning("Failed to get companion by ID: ", httpErr)
		return types.ResultErr[models.SessionDBMongo](httpErr)
	}
	student := studentResult.Value()
	companion := companionResult.Value()

	specialityResult := Speciality.GetByID(companion.IDSpeciality.Hex())
	if specialityResult.IsErr() {
		httpErr := specialityResult.Error().(*types.HttpError)
		logger.Warning("Failed to get speciality by ID: ", httpErr)
		return types.ResultErr[models.SessionDBMongo](httpErr)
	}

	extraInfo := make(map[string]string, 5)
	extraInfo["StudentName"] = student.FirstName
	extraInfo["StudentSurname"] = student.LastName
	extraInfo["CompanionName"] = companion.FirstName
	extraInfo["CompanionSurname"] = companion.LastName
	extraInfo["CompanionSpeciality"] = specialityResult.Value().Name

	sessionOptional := u.ToInsert(extraInfo)
	if !sessionOptional.IsPresent() {
		logger.Warning("Failed to create session: Invalid session data")
		httpErr := types.ErrorInternal(
			"Failed to create session",
			"Invalid session data provided",
			"Session data: "+utils.StructToString(u),
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}
	session := sessionOptional.Get()
	result, err := configs.DB.InsertOne(session)

	if err != nil {
		logger.Warning("Error inserting session: ", err)
		return types.ResultErr[models.SessionDBMongo](err)
	}

	session.ID, err = models.DBIDFrom(result.InsertedID)

	if err != nil {
		logger.Error("Error converting inserted ID to PrimitiveID: ", err)
		httpErr := types.ErrorInternal(
			"Failed to create session",
			"Failed to convert inserted ID to PrimitiveID",
			"Error: "+err.Error(),
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	return types.ResultOk(session)
}

func (sessionType) GetByID(id string) types.Result[models.SessionDBMongo] {
	oid, err := models.BsonIDFrom(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	var session models.SessionDBMongo

	err = configs.DB.FindOne(filter, &session)
	if err != nil {
		var httpErr types.HttpError

		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get session by ID: ", err)
			httpErr = types.ErrorNotFound(
				"Session not found",
				"Session with ID "+id+" not found",
			)
		} else {
			logger.Error("Failed to get session by ID: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve session",
				"Decoding error",
				err.Error(),
				"Session ID: "+id,
			)
		}

		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	return types.ResultOk(session)
}
func (sessionType) GetAll() types.Result[[]models.SessionDBMongo] {
	filter := bson.D{{Key: "deleted_at", Value: models.TimeZero()}} // Filter to exclude deleted sessions
	sessions := []models.SessionDBMongo{}

	err := configs.DB.FindAll(filter, &sessions)
	if err != nil {
		logger.Error("Failed to get all sessions from MongoDB:", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve sessions",
			err.Error(),
		)

		return types.ResultErr[[]models.SessionDBMongo](&httpErr)
	}

	logger.Debug("Retrieved", len(sessions), "sessions from MongoDB database")
	return types.ResultOk(sessions)
}
func (sessionType) GetAllByStudentID(id string) types.Result[[]models.SessionDBMongo] {
	oid, err := models.DBIDFrom(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[[]models.SessionDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "id_student", Value: oid}, {Key: "deleted_at", Value: models.TimeZero()}} // Filter to exclude deleted sessions
	sessions := []models.SessionDBMongo{}

	err = configs.DB.FindAll(filter, &sessions)
	if err != nil {
		logger.Error("Failed to get all sessions by student ID from MongoDB:", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve sessions by student ID",
			err.Error(),
		)

		return types.ResultErr[[]models.SessionDBMongo](&httpErr)
	}

	logger.Debug("Retrieved", len(sessions), "sessions for student ID", id, "from MongoDB database")
	return types.ResultOk(sessions)
}
