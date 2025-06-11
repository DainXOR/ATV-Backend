package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type studentType struct{}

var Student studentType

// Mongo
func (studentType) CreateMongo(user models.StudentCreate) types.Result[models.StudentDBMongo] {
	userDB := user.ToDBMongo()
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
		logger.Error("Failed to get all students from MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve students",
			err.Error(),
		)

		return types.ResultErr[[]models.StudentDBMongo](&httpErr)
	}

	users := types.Map(usersR, models.StudentDBMongoReceiver.ToDB)
	logger.Debug("Retrieved ", len(users), " users from MongoDB database")
	return types.ResultOk(users)
}

func (studentType) UpdateMongo(id string, user models.StudentCreate) types.Result[models.StudentDBMongo] {
	oid, err := primitive.ObjectIDFromHex(id)
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
	update := bson.D{{Key: "$set", Value: user.ToDBMongo()}}

	ctx, cancel := configs.DB.Mongo().Context()
	defer cancel()

	result, err := configs.DB.Mongo().From(models.StudentDBMongo{}).UpdateOne(ctx, filter, update)

	if err != nil {
		logger.Error("Failed to update student in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update student",
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

	return Student.GetByIDMongo(id)
}

func (studentType) PatchMongo(id string, student models.StudentCreate) types.Result[models.StudentDBMongo] {
	oid, err := primitive.ObjectIDFromHex(id)
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
	studentDB := student.ToDBMongo()
	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: studentDB}}

	ctx, cancel := configs.DB.Mongo().Context()
	defer cancel()

	result, err := configs.DB.Mongo().From(models.StudentDBMongo{}).UpdateOne(ctx, filter, update)

	if err != nil {
		logger.Error("Failed to update student in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update student",
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

	return Student.GetByIDMongo(id)
}

// GORM
func (studentType) CreateGorm(user models.StudentCreate) types.Result[models.UserDBGorm] {
	//if User.GetUserByEmail(user.Email).IsOk() {
	//	return types.ResultErr[models.UserDB](models.Error(
	//		types.Http.Conflict(),
	//		"conflict",
	//		"Email is already in use",
	//	))
	//}

	newUser := user.ToDBGorm()

	if res := Student.GetByIDNumberGorm(newUser.IDNumber); res.IsOk() {
		logger.Error("User with ID number already exists: ", newUser.IDNumber)
		err := types.Error(
			types.Http.Conflict(),
			"User already exists",
			"User with ID number "+newUser.IDNumber+" already exists",
		)

		return types.ResultErr[models.UserDBGorm](&err)
	} else {
		if res.Error().(*types.HttpError).Code == types.Http.NotFound() {
			logger.Debug("Creating user")

			configs.DB.Gorm().DB().Create(&newUser)

			logger.Debug("User id: ", newUser.ID)

			if newUser.ID == 0 {
				err := types.ErrorInternal(
					"Failed to create user",
				)

				return types.ResultErr[models.UserDBGorm](&err)
			}

			return types.ResultOk(newUser)
		} else {
			logger.Error("Failed to create user: ", res.Error())
			return types.ResultErr[models.UserDBGorm](res.Error())
		}
	}
}

func (studentType) GetByIDGorm(id string) types.Result[models.UserDBGorm] {
	var user models.UserDBGorm
	configs.DB.Gorm().DB().First(&user, id)
	if user.ID == 0 {
		err := types.ErrorNotFound(
			"User not found",
			"User with ID "+id+" not found",
		)
		return types.ResultErr[models.UserDBGorm](&err)
	}
	return types.ResultOk(user)
}
func (studentType) GetByIDNumberGorm(idNumber string) types.Result[models.UserDBGorm] {
	var user models.UserDBGorm
	configs.DB.Gorm().DB().Where("id_number = ?", idNumber).First(&user)
	if user.ID == 0 {
		err := types.ErrorNotFound(
			"User not found",
			"User with ID number "+idNumber+" not found",
		)
		return types.ResultErr[models.UserDBGorm](&err)
	}
	return types.ResultOk(user)
}
func (studentType) GetAllGorm() types.Result[[]models.UserDBGorm] {
	var users []models.UserDBGorm

	configs.DB.Gorm().DB().Find(&users)

	if len(users) == 0 {
		err := types.ErrorNotFound(
			"No users found",
			"No users found in the database",
		)
		return types.ResultErr[[]models.UserDBGorm](&err)
	}

	logger.Debug("Retrieved ", len(users), " users from GORM database")
	return types.ResultOk(users)
}

func (studentType) UpdateGorm(id string, user models.StudentCreate) types.Result[models.UserDBGorm] {
	var userDB models.UserDBGorm
	configs.DB.Gorm().DB().First(&userDB, id)

	if userDB.ID == 0 {
		err := types.ErrorNotFound(
			"User not found",
			"User with ID "+id+" not found",
		)
		return types.ResultErr[models.UserDBGorm](&err)
	}

	configs.DB.Gorm().DB().Model(&userDB).Updates(user.ToPutDBGorm())

	return types.ResultOk(userDB)
}

func (studentType) PatchGorm(id string, user models.StudentCreate) types.Result[models.UserDBGorm] {
	var userDB models.UserDBGorm
	configs.DB.Gorm().DB().First(&userDB, id)

	if userDB.ID == 0 {
		err := types.ErrorNotFound(
			"User not found",
			"User with ID "+id+" not found",
		)
		return types.ResultErr[models.UserDBGorm](&err)
	}

	configs.DB.Gorm().DB().Model(&userDB).Updates(user.ToDBGorm())

	return types.ResultOk(userDB)
}
