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

type userType struct{}

var User userType

func (userType) GetByIDMongo(id string) types.Result[models.UserDBMongo] {
	oid, err := bson.ObjectIDFromHex(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"User ID: "+id,
		)
		return types.ResultErr[models.UserDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	var userF models.UserDBMongoReceiver

	err = configs.DB.Mongo().FindOne(filter, &userF)
	if err != nil {
		var httpErr types.HttpError

		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get user by ID: ", err)
			httpErr = types.ErrorNotFound(
				"User not found",
				"User with ID "+id+" not found",
			)
		} else {
			logger.Error("Failed to get user by ID: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve user",
				"Decoding error",
				err.Error(),
				"User ID: "+id,
			)
		}

		return types.ResultErr[models.UserDBMongo](&httpErr)
	}

	return types.ResultOk(userF.ToDB())
}
func (userType) GetByIDGorm(id string) types.Result[models.UserDBGorm] {
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
func (userType) GetByIDNumberGorm(idNumber string) types.Result[models.UserDBGorm] {
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
func (userType) GetByIDNumberMongo(idNumber string) types.Result[models.UserDBMongo] {
	filter := bson.D{{Key: "id_number", Value: idNumber}}
	var user models.UserDBMongoReceiver
	err := configs.DB.Mongo().FindOne(filter, &user)
	if err != nil {
		var httpErr types.HttpError
		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get user by ID number: ", err)
			httpErr = types.ErrorNotFound(
				"User not found",
				"User with ID number "+idNumber+" not found",
			)
		} else {
			logger.Error("Failed to get user by ID number: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve user",
				"Decoding error",
				err.Error(),
				"User ID number: "+idNumber,
			)
		}

		return types.ResultErr[models.UserDBMongo](&httpErr)
	}

	return types.ResultOk(user.ToDB())
}

func (userType) GetByEmailMongo(email string) types.Result[models.UserDBMongo] {
	filter := bson.D{{Key: "email", Value: email}}
	var userF models.UserDBMongoReceiver

	err := configs.DB.Mongo().FindOne(filter, &userF)
	if err != nil {
		var httpErr types.HttpError

		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get user by email: ", err)
			httpErr = types.ErrorNotFound(
				"User not found",
				"User with email "+email+" not found",
			)
		} else {
			logger.Error("Failed to get user by email: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve user",
				"Decoding error",
				err.Error(),
				"User email: "+email,
			)
		}

		return types.ResultErr[models.UserDBMongo](&httpErr)
	}

	return types.ResultOk(userF.ToDB())
}

func (userType) GetAllGorm() types.Result[[]models.UserDBGorm] {
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
func (userType) GetAllMongo() types.Result[[]models.UserDBMongo] {
	filter := bson.D{{Key: "deleted_at", Value: nil}} // Filter to exclude deleted users
	usersR := models.UserDBMongo{}.ReceiverList()

	err := configs.DB.Mongo().FindAll(filter, &usersR)
	if err != nil {
		logger.Error("Failed to get all users from MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve users",
			err.Error(),
		)

		return types.ResultErr[[]models.UserDBMongo](&httpErr)
	}

	users := types.Map(usersR, models.UserDBMongoReceiver.ToDB)
	logger.Debug("Retrieved ", len(users), " users from MongoDB database")
	return types.ResultOk(users)
}

func (userType) CreateGorm(user models.UserCreate) types.Result[models.UserDBGorm] {
	//if User.GetUserByEmail(user.Email).IsOk() {
	//	return types.ResultErr[models.UserDB](models.Error(
	//		types.Http.Conflict(),
	//		"conflict",
	//		"Email is already in use",
	//	))
	//}

	newUser := user.ToDBGorm()

	if res := User.GetByIDNumberGorm(newUser.IDNumber); res.IsOk() {
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
func (userType) CreateMongo(user models.UserCreate) types.Result[models.UserDBMongo] {
	userDB := user.ToDBMongo()
	result, err := configs.DB.Mongo().InsertOne(userDB)

	if err != nil {
		logger.Error("Failed to create user in MongoDB: ", err)
		return types.ResultErr[models.UserDBMongo](err)
	}

	logger.Debug("User created with ID: ", result.InsertedID)
	userDB.ID, err = primitive.ObjectIDFromHex(result.InsertedID.(bson.ObjectID).Hex())

	if err != nil {
		logger.Error("Failed to convert inserted ID to ObjectID: ", err)
		newErr := types.ErrorInternal(
			"Failed to create user",
			"Failed to convert inserted ID to ObjectID",
			"Error: "+err.Error(),
		)

		return types.ResultErr[models.UserDBMongo](&newErr)
	}

	return types.ResultOk(userDB)
}

func (userType) UpdateGorm(id string, user models.UserCreate) types.Result[models.UserDBGorm] {
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
func (userType) UpdateMongo(id string, user models.UserCreate) types.Result[models.UserDBMongo] {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"User ID: "+id,
		)
		return types.ResultErr[models.UserDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: user.ToDBMongo()}}

	ctx, cancel := configs.DB.Mongo().Context()
	defer cancel()

	result, err := configs.DB.Mongo().From(models.UserDBMongo{}).UpdateOne(ctx, filter, update)

	if err != nil {
		logger.Error("Failed to update user in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update user",
			err.Error(),
			"User ID: "+id,
		)
		return types.ResultErr[models.UserDBMongo](&httpErr)
	}

	if result.MatchedCount == 0 {
		httpErr := types.ErrorNotFound(
			"User not found",
			"User with ID "+id+" not found",
		)
		return types.ResultErr[models.UserDBMongo](&httpErr)
	}

	return User.GetByIDMongo(id)
}

func (userType) PatchGorm(id string, user models.UserCreate) types.Result[models.UserDBGorm] {
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
