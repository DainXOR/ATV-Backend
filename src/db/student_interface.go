package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type studentType struct{}

var Student studentType

func (studentType) Create(student models.StudentCreate) types.Result[models.StudentDBMongo] {
	studentDB := student.ToInsert()
	if studentDB.IsEmpty() {
		logger.Error("Error converting student to DB model")
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid student data",
			"Student data: "+student.IDUniversity,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	result, err := configs.DB.InsertOne(studentDB)

	if err != nil {
		logger.Error("Failed to create student in MongoDB: ", err)
		return types.ResultErr[models.StudentDBMongo](err)
	}

	studentDB.ID, err = models.ID.ToDB(result.InsertedID)

	if err != nil {
		logger.Error("Failed to convert inserted ID to ObjectID: ", err)
		httpErr := types.ErrorInternal(
			"Failed to create student",
			"Failed to convert inserted ID to ObjectID",
			"Error: "+err.Error(),
		)

		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	return types.ResultOk(studentDB)
}

func (studentType) GetByID(id string) types.Result[models.StudentDBMongo] {
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	var student models.StudentDBMongo

	err = configs.DB.FindOne(filter, &student)
	if err != nil {
		var httpErr types.HttpError

		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get student by ID: ", err)
			httpErr = types.ErrorNotFound(
				"Student not found",
				"Student with ID "+id+" not found",
			)
		} else {
			logger.Error("Failed to get student by ID: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve student",
				"Decoding error",
				err.Error(),
				"Student ID: "+id,
			)
		}

		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	return types.ResultOk(student)
}
func (studentType) GetByNumberID(idNumber string) types.Result[models.StudentDBMongo] {
	filter := bson.D{{Key: "id_number", Value: idNumber}}
	var student models.StudentDBMongo

	err := configs.DB.FindOne(filter, &student)
	if err != nil {
		var httpErr types.HttpError
		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get student by ID number: ", err)
			httpErr = types.ErrorNotFound(
				"Student not found",
				"Student with ID number "+idNumber+" not found",
			)
		} else {
			logger.Error("Failed to get student by ID number: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve student",
				"Decoding error",
				err.Error(),
				"Student ID number: "+idNumber,
			)
		}

		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	return types.ResultOk(student)
}
func (studentType) GetByEmail(email string) types.Result[models.StudentDBMongo] {
	filter := bson.D{{Key: "email", Value: email}}
	var student models.StudentDBMongo

	err := configs.DB.FindOne(filter, &student)
	if err != nil {
		var httpErr types.HttpError

		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get student by email: ", err)
			httpErr = types.ErrorNotFound(
				"Student not found",
				"Student with email "+email+" not found",
			)
		} else {
			logger.Error("Failed to get student by email: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve student",
				"Decoding error",
				err.Error(),
				"Student email: "+email,
			)
		}

		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	return types.ResultOk(student)
}
func (studentType) GetAll() types.Result[[]models.StudentDBMongo] {
	filter := bson.D{{Key: "deleted_at", Value: models.Time.Zero()}} // Filter to exclude deleted students
	students := []models.StudentDBMongo{}

	err := configs.DB.FindAll(filter, &students)
	if err != nil {
		logger.Error("Failed to get all students from MongoDB:", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve students",
			err.Error(),
		)

		return types.ResultErr[[]models.StudentDBMongo](&httpErr)
	}

	logger.Debug("Retrieved", len(students), "students from MongoDB database")
	return types.ResultOk(students)
}

func (studentType) UpdateByID(id string, student models.StudentCreate) types.Result[models.StudentDBMongo] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	studentDB := student.ToUpdate()
	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: studentDB}}

	result := configs.DB.PatchOne(filter, update, &studentDB)

	if result.IsErr() {
		err := result.Error()
		logger.Error("Failed to update student in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update student",
			err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	if result.Value().MatchedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Student not found",
			"Student with ID "+id+" not found",
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	if result.Value().ModifiedCount == 0 {
		logger.Info("No changes made to student with ID: ", id)
		httpErr := types.Error(
			types.Http.C300().NotModified(),
			"No changes made",
			"Student with ID "+id+" was not modified",
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	return types.ResultOk(studentDB)
}

func (studentType) PatchByID(id string, student models.StudentCreate) types.Result[models.StudentDBMongo] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	studentDB := student.ToUpdate()
	if studentDB == (models.StudentDBMongo{}) {
		logger.Error("Error converting student to DB model")
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid student data",
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: studentDB}}

	result := configs.DB.PatchOne(filter, update, &studentDB)

	if result.IsErr() {
		err := result.Error()
		logger.Error("Failed to update student in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update student",
			err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	if result.Value().MatchedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Student not found",
			"Student with ID "+id+" not found",
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	if result.Value().ModifiedCount == 0 {
		logger.Info("No changes made to student with ID: ", id)
		logger.Lava("0.1.2", "Send a more proper code for no changes made")
		httpErr := types.Error(
			types.Http.C200().Accepted(),
			"No changes made",
			"Student with ID "+id+" was not modified",
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	return types.ResultOk(studentDB)
}

func (studentType) DeleteByID(id string) types.Result[models.StudentDBMongo] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: models.Time.Now()}}}}

	var deletedStudent models.StudentDBMongo
	result := configs.DB.UpdateOne(filter, update, &deletedStudent)
	if result.IsErr() {
		logger.Error("Failed to delete student in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to delete student",
			err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	if result.Value().MatchedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Student not found",
			"Student with ID "+id+" not found",
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	err = configs.DB.FindOne(filter, &deletedStudent)
	if err != nil {
		logger.Error("Failed to retrieve deleted student: ", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve deleted student",
			err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	return types.ResultOk(deletedStudent)
}
func (studentType) DeletePermanentByID(id string) types.Result[models.StudentDBMongo] {
	logger.Warning("Permanently deleting student by ID: ", id)
	oid, err := models.BsonIDFrom(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}, {Key: "deleted_at", Value: bson.M{"$ne": time.Time{}}}} // Ensure the student is marked as deleted
	ctx, cancel := configs.DB.Context()
	defer cancel()

	var student models.StudentDBMongo
	err = configs.DB.FindOne(filter, &student)
	if err != nil {
		logger.Debug("Failed to find student for permanent deletion: ", err)

		if err == mongo.ErrNoDocuments {
			httpErr := types.ErrorNotFound(
				"Student not found",
				"Student with ID "+id+" not found or not marked as deleted",
			)
			return types.ResultErr[models.StudentDBMongo](&httpErr)
		}

		httpErr := types.ErrorInternal(
			"Failed to find student for permanent deletion",
			err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	result, err := configs.DB.From(models.StudentDBMongo{}).DeleteOne(ctx, filter)
	if err != nil {
		logger.Debug("Failed to permanently delete student in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to permanently delete student",
			err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	if result.DeletedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Student not found.",
			"Student with ID "+id+" not found.",
			"Ensure the student is marked as deleted before permanent deletion.",
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	return types.ResultOk(student)
}
func (studentType) DeletePermanentAll() types.Result[[]models.StudentDBMongo] {
	filter := bson.D{{Key: "deleted_at", Value: bson.M{"$ne": nil}}}
	ctx, cancel := configs.DB.Context()
	defer cancel()

	result, err := configs.DB.From(models.StudentDBMongo{}).DeleteMany(ctx, filter)
	if err != nil {
		logger.Error("Failed to permanently delete all students in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to permanently delete all students",
			err.Error(),
		)
		return types.ResultErr[[]models.StudentDBMongo](&httpErr)
	}

	if result.DeletedCount == 0 {
		httpErr := types.ErrorNotFound(
			"No deleted students found",
			"No students marked as deleted found",
		)
		return types.ResultErr[[]models.StudentDBMongo](&httpErr)
	}

	return types.ResultOk([]models.StudentDBMongo{})
}
