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

func (userType) GetByID(id string) types.Result[models.UserDBMongo] {
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
func (userType) GetByEmail(email string) types.Result[models.UserDBMongo] {
	var user models.UserDBMongo

	ctx, cancel := configs.DB.Mongo().Context()
	defer cancel()

	configs.DB.Mongo().In("users").FindOne(
		ctx,
		bson.M{"email": email},
	).Decode(&user)

	err := types.ErrorNotFound(
		"User not found",
		"User with email "+email+" not found",
	)
	return types.ResultOf(user, &err, user.ID != primitive.NilObjectID)
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
}
func (userType) Create(user models.UserCreate) types.Result[models.UserDBMongo] {
	var userDB models.UserDBMongo
	userDB = user.ToDBMongo()

	ctx, cancel := configs.DB.Mongo().Context()
	defer cancel()

	result, err := configs.DB.Mongo().In(userDB.TableName()).InsertOne(
		ctx,
		userDB,
	)

	if err != nil {
		logger.Error("Failed to create user in MongoDB: ", err)
		return types.ResultErr[models.UserDBMongo](err)
	}

	logger.Debug("User created with ID: ", result.InsertedID)
	logger.Debug("Test: ", result.InsertedID.(bson.ObjectID).Hex())
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
	//logger.Debug("User created with ID: ", userDB.ID.String())

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
