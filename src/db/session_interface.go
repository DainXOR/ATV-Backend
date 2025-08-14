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

	sessionOptional := utils.Transform(getExtraInfo(u), func(res types.Result[map[string]string]) types.Optional[models.SessionDBMongo] {
		if res.IsErr() {
			return types.OptionalEmpty[models.SessionDBMongo]()
		}
		return u.ToInsert(res.Value())
	})

	if sessionOptional.IsEmpty() {
		logger.Warning("Failed to create session: Invalid session data")
		httpErr := types.ErrorInternal(
			"Failed to create session",
			"Invalid session data provided",
			"Session data: "+utils.StructToString(u),
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}
	session := sessionOptional.Get()
	logger.Debug("Session object to insert: ", session)
	result, err := configs.DB.InsertOne(session)

	if err != nil {
		logger.Warning("Error inserting session: ", err)
		return types.ResultErr[models.SessionDBMongo](err)
	}

	session.ID, err = models.ID.ToDB(result.InsertedID)

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
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
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
	filter := bson.D{models.Filter.NotDeleted()} // Filter to exclude deleted sessions
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
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[[]models.SessionDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "id_student", Value: oid}, models.Filter.NotDeleted()} // Filter to exclude deleted sessions
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

func (sessionType) UpdateByID(id string, session models.SessionCreate) types.Result[models.SessionDBMongo] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	sessionData := utils.Transform(getExtraInfo(session), func(res types.Result[map[string]string]) types.Result[models.SessionDBMongo] {
		if res.IsErr() {
			return types.ResultErr[models.SessionDBMongo](res.Error())
		}
		return session.ToUpdate(res.Value())
	})
	if sessionData.IsErr() {
		logger.Warning("Failed to update session:", sessionData.Error())
		httpErr := types.ErrorInternal(
			"Failed to update session",
			"Invalid session data",
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	sessionDB := sessionData.Value()
	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: sessionDB}}

	result := configs.DB.UpdateOne(filter, update, &sessionDB)

	if result.IsErr() {
		err := result.Error()
		logger.Error("Failed to update session in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update session",
			err.Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	if result.Value().MatchedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Session not found",
			"Session with ID "+id+" not found",
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	if result.Value().ModifiedCount == 0 {
		logger.Info("No changes made to session with ID: ", id)
		httpErr := types.Error(
			types.Http.C300().NotModified(),
			"No changes made",
			"Session with ID "+id+" was not modified",
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	return types.ResultOk(sessionDB)
}

func (sessionType) PatchByID(id string, session models.SessionCreate) types.Result[models.SessionDBMongo] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	sessionData := utils.Transform(getExtraInfoAllowEmpty(session),
		func(res types.Result[map[string]string]) types.Result[models.SessionDBMongo] {
			if res.IsErr() {
				return types.ResultErr[models.SessionDBMongo](res.Error())
			}
			return session.ToUpdate(res.Value())
		},
	)
	if sessionData.IsErr() {
		logger.Warning("Failed to update session:", sessionData.Error())
		httpErr := types.ErrorInternal(
			"Failed to update session",
			"Invalid session data",
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	sessionDB := sessionData.Value()
	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: sessionDB}}

	result := configs.DB.PatchOne(filter, update, &sessionDB)

	if result.IsErr() {
		err := result.Error()
		logger.Error("Failed to patch session in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to patch session",
			err.Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	if result.Value().MatchedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Session not found",
			"Session with ID "+id+" not found",
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	if result.Value().ModifiedCount == 0 {
		logger.Info("No changes made to session with ID: ", id)
		httpErr := types.Error(
			types.Http.C300().NotModified(),
			"No changes made",
			"Session with ID "+id+" was not modified",
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	return types.ResultOk(sessionDB)
}

func (sessionType) DeleteByID(id string) types.Result[models.SessionDBMongo] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: models.Time.Now()}}}}

	var deletedSession models.SessionDBMongo
	result := configs.DB.UpdateOne(filter, update, &deletedSession)
	if result.IsErr() {
		logger.Error("Failed to delete session in MongoDB: ", result.Error())
		httpErr := types.ErrorInternal(
			"Failed to delete session",
			result.Error().Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	if result.Value().MatchedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Session not found",
			"Session with ID "+id+" not found",
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	err = configs.DB.FindOne(filter, &deletedSession)
	if err != nil {
		logger.Error("Failed to retrieve deleted session: ", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve deleted session",
			err.Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDBMongo](&httpErr)
	}

	return types.ResultOk(deletedSession)
}

func getExtraInfo(session models.SessionCreate) types.Result[map[string]string] {
	studentResult := Student.GetByID(session.IDStudent)
	if studentResult.IsErr() {
		httpErr := studentResult.Error().(*types.HttpError)
		logger.Warning("Failed to get student by ID: ", httpErr)
		return types.ResultErr[map[string]string](httpErr)
	}

	companionResult := Companion.GetByID(session.IDCompanion)
	if companionResult.IsErr() {
		httpErr := companionResult.Error().(*types.HttpError)
		logger.Warning("Failed to get companion by ID: ", httpErr)
		return types.ResultErr[map[string]string](httpErr)
	}
	student := studentResult.Value()
	companion := companionResult.Value()

	specialityResult := Speciality.GetByID(companion.IDSpeciality.Hex())
	if specialityResult.IsErr() {
		httpErr := specialityResult.Error().(*types.HttpError)
		logger.Warning("Failed to get speciality by ID: ", httpErr)
		return types.ResultErr[map[string]string](httpErr)
	}

	extraInfo := make(map[string]string, 5)
	extraInfo["StudentName"] = student.FirstName
	extraInfo["StudentSurname"] = student.LastName
	extraInfo["CompanionName"] = companion.FirstName
	extraInfo["CompanionSurname"] = companion.LastName
	extraInfo["CompanionSpeciality"] = specialityResult.Value().Name

	return types.ResultOk(extraInfo)
}

func getExtraInfoAllowEmpty(session models.SessionCreate) types.Result[map[string]string] {
	var student models.StudentDBMongo
	if session.IDStudent != "" {
		studentResult := Student.GetByID(session.IDStudent)
		if studentResult.IsErr() {
			httpErr := studentResult.Error().(*types.HttpError)
			logger.Warning("Failed to get student by ID: ", httpErr)
			return types.ResultErr[map[string]string](httpErr)
		}

		student = studentResult.Value()
	}

	var companion models.CompanionDBMongo
	var speciality models.SpecialityDBMongo
	if session.IDCompanion != "" {
		companionResult := Companion.GetByID(session.IDCompanion)
		if companionResult.IsErr() {
			httpErr := companionResult.Error().(*types.HttpError)
			logger.Warning("Failed to get companion by ID: ", httpErr)
			return types.ResultErr[map[string]string](httpErr)
		}

		companion = companionResult.Value()

		specialityResult := Speciality.GetByID(companion.IDSpeciality.Hex())
		if specialityResult.IsErr() {
			httpErr := specialityResult.Error().(*types.HttpError)
			logger.Warning("Failed to get speciality by ID: ", httpErr)
			return types.ResultErr[map[string]string](httpErr)
		}

		speciality = specialityResult.Value()
	}

	extraInfo := make(map[string]string, 5)
	extraInfo["StudentName"] = student.FirstName
	extraInfo["StudentSurname"] = student.LastName
	extraInfo["CompanionName"] = companion.FirstName
	extraInfo["CompanionSurname"] = companion.LastName
	extraInfo["CompanionSpeciality"] = speciality.Name

	return types.ResultOk(extraInfo)
}
