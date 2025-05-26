package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type userType struct{}

var User userType

func (userType) GetUserByID(id string) types.Result[models.UserDBGorm] {
	var user models.UserDBGorm
	configs.DataBase.First(&user, id)

	err := types.ErrorNotFound(
		"User not found",
		"User with ID "+id+" not found",
	)
	return types.ResultOf(user, &err, user.ID != 0)
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

func (userType) GetByID(id string) types.Result[models.UserDBGorm] {
	var user models.UserDBGorm

	configs.DB.Mongo().Collection("users").FindOne(
		configs.DB.Mongo().Context(),
		bson.M{"_id": id},
	).Decode(&user)
	err := types.ErrorNotFound(
		"User not found",
		"User with ID "+id+" not found",
	)

	return types.ResultOf(user, &err, user.ID != 0)
}
