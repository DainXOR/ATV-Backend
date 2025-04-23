package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
)

type userType struct{}

var User userType

func (userType) GetUserByID(id string) types.Result[models.UserDB] {
	var user models.UserDB
	configs.DataBase.First(&user, id)

	err := types.ErrorNotFound(
		"User not found",
		"User with ID "+id+" not found",
	)
	return types.ResultOf(user, &err, user.ID != 0)
}

func (userType) CreateUser(user models.UserCreate) types.Result[models.UserDB] {
	//if User.GetUserByEmail(user.Email).IsOk() {
	//	return types.ResultErr[models.UserDB](models.Error(
	//		types.Http.Conflict(),
	//		"conflict",
	//		"Email is already in use",
	//	))
	//}

	newUser := user.ToDB()

	logger.Debug("Creating user")

	configs.DataBase.Create(&newUser)

	logger.Debug("User id: ", newUser.ID)

	if newUser.ID == 0 {
		err := types.Error(
			types.Http.InternalServerError(),
			"internal",
			"User not created",
		)

		return types.ResultErr[models.UserDB](&err)
	}

	return types.ResultOk(newUser)
}
