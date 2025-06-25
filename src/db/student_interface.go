package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type studentType struct{}

var Student studentType

// Mongo
func (studentType) CreateMongo(user models.StudentCreate) types.Result[models.StudentDBMongo] {
	userDB := user.ToInsert()
	result, err := configs.DB.Mongo().InsertOne(userDB)

	if err != nil {
		logger.Error("Failed to create student in MongoDB: ", err)
		return types.ResultErr[models.StudentDBMongo](err)
	}

	logger.Debug("Student created with ID: ", result.InsertedID)
	userDB.ID, err = primitive.ObjectIDFromHex(result.InsertedID.(bson.ObjectID).Hex())

	if err != nil {
		logger.Error("Failed to convert inserted ID to ObjectID: ", err)
		newErr := types.ErrorInternal(
			"Failed to create student",
			"Failed to convert inserted ID to ObjectID",
			"Error: "+err.Error(),
		)

		return types.ResultErr[models.StudentDBMongo](&newErr)
	}

	return types.ResultOk(userDB)
}

func (studentType) GetByIDMongo(id string) types.Result[models.StudentDBMongo] {
	oid, err := bson.ObjectIDFromHex(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"User ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	var userF models.StudentDBMongoReceiver

	err = configs.DB.Mongo().FindOne(filter, &userF)
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

	return types.ResultOk(userF.ToDB())
}
func (studentType) GetByIDNumberMongo(idNumber string) types.Result[models.StudentDBMongo] {
	filter := bson.D{{Key: "id_number", Value: idNumber}}
	var user models.StudentDBMongoReceiver

	err := configs.DB.Mongo().FindOne(filter, &user)
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

	return types.ResultOk(user.ToDB())
}
func (studentType) GetByEmailMongo(email string) types.Result[models.StudentDBMongo] {
	filter := bson.D{{Key: "email", Value: email}}
	var userF models.StudentDBMongoReceiver

	err := configs.DB.Mongo().FindOne(filter, &userF)
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

	return types.ResultOk(userF.ToDB())
}
func (studentType) GetAllMongo() types.Result[[]models.StudentDBMongo] {
	filter := bson.D{{Key: "deleted_at", Value: nil}} // Filter to exclude deleted students
	usersR := models.StudentDBMongo{}.ReceiverList()

	err := configs.DB.Mongo().FindAll(filter, &usersR)
	if err != nil {
		logger.Error("Failed to get all students from MongoDB:", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve students",
			err.Error(),
		)

		return types.ResultErr[[]models.StudentDBMongo](&httpErr)
	}

	users := utils.Map(usersR, models.StudentDBMongoReceiver.ToDB)
	logger.Debug("Retrieved", len(users), "users from MongoDB database")
	return types.ResultOk(users)
}

func (studentType) UpdateMongo(id string, user models.StudentCreate) types.Result[models.StudentDBMongo] {
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

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: user.ToUpdate()}}
	studentDB := user.ToUpdate().Receiver()

	result := configs.DB.Mongo().PatchOne(filter, update, studentDB)
	// .From(models.StudentDBMongo{}).UpdateOne(ctx, filter, update)

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
		logger.Lava(2, "Send a more proper code for no changes made")
		httpErr := types.Error(
			types.Http.C200().Accepted(),
			"No changes made",
			"Student with ID "+id+" was not modified",
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	return Student.GetByIDMongo(id)
}

func (studentType) PatchMongo(id string, student models.StudentCreate) types.Result[models.StudentDBMongo] {
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
	receiver := studentDB.Receiver()

	result := configs.DB.Mongo().PatchOne(filter, update, &receiver)

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
		logger.Lava(2, "Send a more proper code for no changes made")
		httpErr := types.Error(
			types.Http.C200().Accepted(),
			"No changes made",
			"Student with ID "+id+" was not modified",
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	return types.ResultOk(receiver.ToDB())
}

func (studentType) DeleteByID(id string) types.Result[models.StudentDBMongo] {
	oid, err := bson.ObjectIDFromHex(id)
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
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: time.Now()}}}}
	ctx, cancel := configs.DB.Mongo().Context()
	defer cancel()

	result, err := configs.DB.Mongo().From(models.StudentDBMongo{}).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to delete student in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to delete student",
			err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	if result.MatchedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Student not found",
			"Student with ID "+id+" not found",
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	var deletedStudent models.StudentDBMongoReceiver
	err = configs.DB.Mongo().FindOne(filter, &deletedStudent)
	if err != nil {
		logger.Error("Failed to retrieve deleted student: ", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve deleted student",
			err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDBMongo](&httpErr)
	}

	deletedStudent.DeletedAt = models.TimeNow() // primitive.NewDateTimeFromTime(configs.DB.Mongo().Now())
	return types.ResultOk(deletedStudent.ToDB())
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
	ctx, cancel := configs.DB.Mongo().Context()
	defer cancel()

	var student models.StudentDBMongoReceiver
	err = configs.DB.Mongo().FindOne(filter, &student)
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

	result, err := configs.DB.Mongo().From(models.StudentDBMongo{}).DeleteOne(ctx, filter)
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

	return types.ResultOk(student.ToDB())
}
func (studentType) DeletePermanentAll() types.Result[[]models.StudentDBMongo] {
	filter := bson.D{{Key: "deleted_at", Value: bson.M{"$ne": nil}}}
	ctx, cancel := configs.DB.Mongo().Context()
	defer cancel()

	result, err := configs.DB.Mongo().From(models.StudentDBMongo{}).DeleteMany(ctx, filter)
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
