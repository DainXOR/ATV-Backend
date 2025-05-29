package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type userType struct{}

var User userType

func (userType) GetByID(id string) types.Result[models.UserDBMongo] {
	var user models.UserDBMongo
	configs.DB.Mongo().Collection("users").FindOne(
		configs.DB.Mongo().Context(),
		bson.M{"_id": id},
	).Decode(&user)

	err := types.ErrorNotFound(
		"User not found",
		"User with ID "+id+" not found",
	)

	return types.ResultOf(user, &err, user.ID != primitive.NilObjectID)
}
func (userType) GetByIDGorm(id string) types.Result[models.UserDBGorm] {
	var user models.UserDBGorm
	configs.DataBase.First(&user, id)
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
	configs.DB.Mongo().Collection("users").FindOne(
		configs.DB.Mongo().Context(),
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

	configs.DataBase.Find(&users)

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
	var users []models.UserDBMongo
	configs.DB.Mongo().Collection("users").Find(
		configs.DB.Mongo().Context(),
		bson.M{},
	)

	if len(users) == 0 {
		err := types.ErrorNotFound(
			"No users found",
			"No users found in the database",
		)
		return types.ResultErr[[]models.UserDBMongo](&err)
	}
	return types.ResultOk(users)
}

func (userType) CreateUser(user models.UserCreate) types.Result[models.UserDBGorm] {
	//if User.GetUserByEmail(user.Email).IsOk() {
	//	return types.ResultErr[models.UserDB](models.Error(
	//		types.Http.Conflict(),
	//		"conflict",
	//		"Email is already in use",
	//	))
	//}

	newUser := user.ToDBGorm()

	logger.Debug("Creating user")

	configs.DataBase.Create(&newUser)

	logger.Debug("User id: ", newUser.ID)

	if newUser.ID == 0 {
		err := types.Error(
			types.Http.InternalServerError(),
			"internal",
			"User not created",
		)

		return types.ResultErr[models.UserDBGorm](&err)
	}

	return types.ResultOk(newUser)
}
func (userType) Create(user models.UserCreate) types.Result[models.UserDBMongo] {
	var userDB models.UserDBMongo
	userDB = user.ToDBMongo()

	logger.Debug("Creating user in MongoDB")

	_, err := configs.DB.Mongo().Collection("users").InsertOne(
		configs.DB.Mongo().Context(),
		userDB,
	)

	if err != nil {
		logger.Error("Failed to create user in MongoDB: ", err)
		return types.ResultErr[models.UserDBMongo](err)
	}

	logger.Debug("User created with ID: ", userDB.ID)

	return types.ResultOk(userDB)
}

func (userType) UpdateUser(id string, user models.UserCreate) types.Result[models.UserDBGorm] {
	var userDB models.UserDBGorm
	configs.DataBase.First(&userDB, id)

	if userDB.ID == 0 {
		err := types.ErrorNotFound(
			"User not found",
			"User with ID "+id+" not found",
		)
		return types.ResultErr[models.UserDBGorm](&err)
	}

	configs.DataBase.Model(&userDB).Updates(user.ToPutDB())

	return types.ResultOk(userDB)
}

func (userType) PatchUser(id string, user models.UserCreate) types.Result[models.UserDBGorm] {
	var userDB models.UserDBGorm
	configs.DataBase.First(&userDB, id)

	if userDB.ID == 0 {
		err := types.ErrorNotFound(
			"User not found",
			"User with ID "+id+" not found",
		)
		return types.ResultErr[models.UserDBGorm](&err)
	}

	configs.DataBase.Model(&userDB).Updates(user.ToDBGorm())

	return types.ResultOk(userDB)
}
