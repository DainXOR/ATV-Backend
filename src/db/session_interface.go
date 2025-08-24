package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type sessionType struct{}

var Session sessionType

func (sessionType) Create(u models.SessionCreate) types.Result[models.SessionDB] {
	logger.Debug("Creating session with data: ", u)

	var session models.SessionDB
	{
		res := getExtraInfo(u)
		if res.IsErr() {
			logger.Warning("Failed to create session: Invalid session data")
			httpErr := types.ErrorInternal(
				"Failed to create session",
				"Invalid session data provided",
				"Session data: "+utils.StructToString(u),
			)
			return types.ResultErr[models.SessionDB](&httpErr)
		}

		sessionOptional := u.ToInsert(res.Value())
		if sessionOptional.IsEmpty() {
			logger.Warning("Failed to create session: Invalid session data")
			httpErr := types.ErrorInternal(
				"Failed to create session",
				"Invalid session data provided",
				"Session data: "+utils.StructToString(u),
			)
			return types.ResultErr[models.SessionDB](&httpErr)
		}

		session = sessionOptional.Get()
	}

	logger.Debug("Session object to insert: ", session)
	resultID := configs.DB.InsertOne(session)

	if resultID.IsErr() {
		logger.Warning("Error inserting session: ", resultID.Error())
		return types.ResultErr[models.SessionDB](resultID.Error())
	}

	session.ID = resultID.Value()
	return types.ResultOk(session)
}

func (sessionType) GetByID(id string) types.Result[models.SessionDB] {
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	sessionResult := configs.DB.FindOne(filter, models.SessionDB{})
	if sessionResult.IsErr() {
		logger.Warning("Failed to get session by ID: ", sessionResult.Error())
		var httpErr types.HttpError

		switch sessionResult.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Session not found",
				"Session with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve session",
				sessionResult.Error().Error(),
				"Session ID: "+id,
			)
		}

		return types.ResultErr[models.SessionDB](&httpErr)
	}

	return types.ResultOk(sessionResult.Value().(models.SessionDB))
}
func (sessionType) GetAll() types.Result[[]models.SessionDB] {
	filter := bson.D{models.Filter.NotDeleted()} // Filter to exclude deleted sessions
	session := models.SessionDB{}

	sessionsResult := configs.DB.FindAll(filter, &session)
	if sessionsResult.IsErr() {
		logger.Warning("Failed to get all sessions from database:", sessionsResult.Error())
		var httpErr types.HttpError

		switch sessionsResult.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Sessions not found",
				"No sessions found matching the criteria",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve sessions",
				sessionsResult.Error().Error(),
			)
		}
		return types.ResultErr[[]models.SessionDB](&httpErr)
	}

	sessions := utils.Map(sessionsResult.Value(), models.InterfaceTo[models.SessionDB])
	logger.Debug("Retrieved", len(sessions), "sessions from database")
	return types.ResultOk(sessions)
}
func (sessionType) GetAllByStudentID(id string) types.Result[[]models.SessionDB] {
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[[]models.SessionDB](&httpErr)
	}

	filter := bson.D{{Key: "id_student", Value: oid}, models.Filter.NotDeleted()} // Filter to exclude deleted sessions
	sessionResult := configs.DB.FindAll(filter, models.SessionDB{})
	if sessionResult.IsErr() {
		logger.Warning("Failed to get all sessions by student ID from database:", sessionResult.Error())
		var httpErr types.HttpError

		switch sessionResult.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Sessions not found",
				"No sessions found for student ID: "+id,
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve sessions by student ID",
				sessionResult.Error().Error(),
			)
		}

		return types.ResultErr[[]models.SessionDB](&httpErr)
	}

	sessions := utils.Map(sessionResult.Value(), models.InterfaceTo[models.SessionDB])
	logger.Debug("Retrieved", len(sessions), "sessions for student ID", id, "from database")
	return types.ResultOk(sessions)
}

func (sessionType) UpdateByID(id string, session models.SessionCreate) types.Result[models.SessionDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDB](&httpErr)
	}

	var sessionDB models.SessionDB
	{
		if res := getExtraInfo(session); res.IsErr() {
			logger.Warning("Failed to update session:", res.Error())
			httpErr := types.ErrorInternal(
				"Failed to update session",
				"Invalid session data",
				res.Error().Error(),
			)
			return types.ResultErr[models.SessionDB](&httpErr)
		} else if sessionResult := session.ToUpdate(res.Value()); sessionResult.IsErr() {
			logger.Warning("Failed to update session:", sessionResult.Error())
			httpErr := types.ErrorInternal(
				"Failed to update session",
				"Invalid session data",
				sessionResult.Error().Error(),
			)
			return types.ResultErr[models.SessionDB](&httpErr)
		} else {
			sessionDB = sessionResult.Value()
		}
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	err = configs.DB.UpdateOne(filter, sessionDB)
	if err != nil {
		logger.Error("Failed to update session in database: ", err)
		return types.ResultErr[models.SessionDB](err)
	}

	sessionResult := configs.DB.FindOne(filter, sessionDB)
	if sessionResult.IsErr() {
		logger.Error("Failed to retrieve updated session from database: ", sessionResult.Error())
		return types.ResultErr[models.SessionDB](sessionResult.Error())
	}

	return types.ResultOk(sessionResult.Value().(models.SessionDB))
}

func (sessionType) PatchByID(id string, session models.SessionCreate) types.Result[models.SessionDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDB](&httpErr)
	}

	var sessionDB models.SessionDB
	{
		if res := getExtraInfoAllowEmpty(session); res.IsErr() {
			logger.Warning("Failed to patch session:", res.Error())
			httpErr := types.ErrorInternal(
				"Failed to patch session",
				"Invalid session data",
				res.Error().Error(),
			)
			return types.ResultErr[models.SessionDB](&httpErr)
		} else if sessionResult := session.ToUpdate(res.Value()); sessionResult.IsErr() {
			logger.Warning("Failed to patch session:", sessionResult.Error())
			httpErr := types.ErrorInternal(
				"Failed to patch session",
				"Invalid session data",
				sessionResult.Error().Error(),
			)
			return types.ResultErr[models.SessionDB](&httpErr)
		} else {
			sessionDB = sessionResult.Value()
		}
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	err = configs.DB.PatchOne(filter, sessionDB)
	if err != nil {
		logger.Error("Failed to patch session in database: ", err)
		return types.ResultErr[models.SessionDB](err)
	}

	sessionResult := configs.DB.FindOne(filter, sessionDB)
	if sessionResult.IsErr() {
		logger.Error("Failed to retrieve patched session from database: ", sessionResult.Error())
		return types.ResultErr[models.SessionDB](sessionResult.Error())
	}

	return types.ResultOk(sessionResult.Value().(models.SessionDB))
}

func (sessionType) DeleteByID(id string) types.Result[models.SessionDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Session ID: "+id,
		)
		return types.ResultErr[models.SessionDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	err = configs.DB.SoftDeleteOne(filter, models.SessionDB{})
	if err != nil {
		logger.Error("Failed to delete session in database: ", err)
		return types.ResultErr[models.SessionDB](err)
	}

	sessionResult := configs.DB.FindOne(filter, models.SessionDB{})
	if sessionResult.IsErr() {
		logger.Error("Failed to retrieve updated session from database: ", sessionResult.Error())
		return types.ResultErr[models.SessionDB](sessionResult.Error())
	}

	return types.ResultOk(sessionResult.Value().(models.SessionDB))
}

func getExtraInfo(session models.SessionCreate) types.Result[map[string]string] {
	studentResult := Student.GetByID(session.IDStudent)
	if studentResult.IsErr() {
		logger.Warning("Failed to get student by ID: ", studentResult.Error())
		return types.ResultErr[map[string]string](studentResult.Error())
	}

	companionResult := Companion.GetByID(session.IDCompanion)
	if companionResult.IsErr() {
		logger.Warning("Failed to get companion by ID: ", companionResult.Error())
		return types.ResultErr[map[string]string](companionResult.Error())
	}
	student := studentResult.Value()
	companion := companionResult.Value()

	specialityResult := Speciality.GetByID(companion.IDSpeciality.Hex())
	if specialityResult.IsErr() {
		logger.Warning("Failed to get speciality by ID: ", specialityResult.Error())
		return types.ResultErr[map[string]string](specialityResult.Error())
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
	var student models.StudentDB
	if session.IDStudent != "" {
		studentResult := Student.GetByID(session.IDStudent)
		if studentResult.IsErr() {
			logger.Warning("Failed to get student by ID: ", studentResult.Error())
			return types.ResultErr[map[string]string](studentResult.Error())
		}

		student = studentResult.Value()
	}

	var companion models.CompanionDB
	var speciality models.SpecialityDB
	if session.IDCompanion != "" {
		companionResult := Companion.GetByID(session.IDCompanion)
		if companionResult.IsErr() {
			logger.Warning("Failed to get companion by ID: ", companionResult.Error())
			return types.ResultErr[map[string]string](companionResult.Error())
		}

		companion = companionResult.Value()

		specialityResult := Speciality.GetByID(companion.IDSpeciality.Hex())
		if specialityResult.IsErr() {
			logger.Warning("Failed to get speciality by ID: ", specialityResult.Error())
			return types.ResultErr[map[string]string](specialityResult.Error())
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
