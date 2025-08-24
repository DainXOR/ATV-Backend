package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type studentType struct{}

var Student studentType

func (studentType) Create(student models.StudentCreate) types.Result[models.StudentDB] {
	studentDB := student.ToInsert()
	if studentDB.IsEmpty() {
		logger.Error("Error converting student to DB model")
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid student data",
			"Student data: "+student.IDUniversity,
		)
		return types.ResultErr[models.StudentDB](&httpErr)
	}

	resultID := configs.DB.InsertOne(studentDB)

	if resultID.IsErr() {
		logger.Error("Failed to create student in MongoDB: ", resultID.Error())
		return types.ResultErr[models.StudentDB](resultID.Error())
	}

	studentDB.ID = resultID.Value()
	return types.ResultOk(studentDB)
}

func (studentType) GetByID(id string) types.Result[models.StudentDB] {
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	resultStudent := configs.DB.FindOne(filter, models.StudentDB{})
	if resultStudent.IsErr() {
		logger.Warning("Failed to get student by ID: ", err)
		var httpErr types.HttpError

		switch resultStudent.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Student not found",
				"Student with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve student",
				resultStudent.Error().Error(),
				"Student ID: "+id,
			)
		}

		return types.ResultErr[models.StudentDB](&httpErr)
	}

	return types.ResultOk(resultStudent.Value().(models.StudentDB))
}
func (studentType) GetByNumberID(idNumber string) types.Result[models.StudentDB] {
	filter := bson.D{{Key: "id_number", Value: idNumber}}

	resultStudent := configs.DB.FindOne(filter, models.StudentDB{})
	if resultStudent.IsErr() {
		logger.Warning("Failed to get student by ID number: ", resultStudent.Error())
		var httpErr types.HttpError

		switch resultStudent.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Student not found",
				"Student with ID number "+idNumber+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve student",
				resultStudent.Error().Error(),
				"Student ID number: "+idNumber,
			)
		}

		return types.ResultErr[models.StudentDB](&httpErr)
	}

	return types.ResultOk(resultStudent.Value().(models.StudentDB))
}
func (studentType) GetByEmail(email string) types.Result[models.StudentDB] {
	filter := bson.D{{Key: "email", Value: email}}

	resultStudent := configs.DB.FindOne(filter, models.StudentDB{})
	if resultStudent.IsErr() {
		logger.Warning("Failed to get student by email: ", resultStudent.Error())
		var httpErr types.HttpError

		switch resultStudent.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Student not found",
				"Student with email "+email+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve student",
				resultStudent.Error().Error(),
				"Student email: "+email,
			)
		}

		return types.ResultErr[models.StudentDB](&httpErr)
	}

	return types.ResultOk(resultStudent.Value().(models.StudentDB))
}
func (studentType) GetAll() types.Result[[]models.StudentDB] {
	filter := bson.D{{Key: "deleted_at", Value: models.Time.Zero()}} // Filter to exclude deleted students

	resultStudents := configs.DB.FindAll(filter, models.StudentDB{})
	if resultStudents.IsErr() {
		logger.Warning("Failed to get all students from MongoDB:", resultStudents.Error())
		var httpErr types.HttpError

		switch resultStudents.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"No students found",
				resultStudents.Error().Error(),
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve students",
				resultStudents.Error().Error(),
			)

		}

		return types.ResultErr[[]models.StudentDB](&httpErr)
	}

	students := utils.Map(resultStudents.Value(), models.InterfaceTo[models.StudentDB])
	logger.Debug("Retrieved", len(students), "students from MongoDB database")
	return types.ResultOk(students)
}

func (studentType) UpdateByID(id string, student models.StudentCreate) types.Result[models.StudentDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDB](&httpErr)
	}

	resultStudent := student.ToUpdate()
	if resultStudent.IsErr() {
		logger.Warning("Error converting student to DB model: ", resultStudent.Error())
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid student data",
			resultStudent.Error().Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	err = configs.DB.UpdateOne(filter, resultStudent.Value())
	if err != nil {
		logger.Warning("Failed to update student in MongoDB: ", err)
		return types.ResultErr[models.StudentDB](err)
	}

	studentDB := configs.DB.FindOne(filter, models.StudentDB{})
	return types.ResultOk(studentDB.Value().(models.StudentDB))
}

func (studentType) PatchByID(id string, student models.StudentCreate) types.Result[models.StudentDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDB](&httpErr)
	}

	resultStudent := student.ToUpdate()
	if resultStudent.IsErr() {
		logger.Warning("Error converting student to DB model: ", resultStudent.Error())
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid student data",
			resultStudent.Error().Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	err = configs.DB.PatchOne(filter, resultStudent.Value())
	if err != nil {
		logger.Warning("Failed to patch student in MongoDB: ", err)
		return types.ResultErr[models.StudentDB](err)
	}

	studentDB := configs.DB.FindOne(filter, models.StudentDB{})
	return types.ResultOk(studentDB.Value().(models.StudentDB))
}

func (studentType) DeleteByID(id string) types.Result[models.StudentDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	resultStudent := configs.DB.FindOne(filter, models.StudentDB{})
	if resultStudent.IsErr() {
		logger.Warning("Failed to retrieve student: ", resultStudent.Error())
		var httpErr types.HttpError

		switch resultStudent.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Student not found",
				"Student with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve student",
				resultStudent.Error().Error(),
				"Student ID: "+id,
			)
		}

		return types.ResultErr[models.StudentDB](&httpErr)
	}

	deletedStudent := resultStudent.Value().(models.StudentDB)
	err = configs.DB.SoftDeleteOne(filter, deletedStudent)
	if err != nil {
		logger.Warning("Failed to delete student in MongoDB: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Student not found",
				"Student with ID "+id+" not found",
			)

		case configs.DBErr.NotModified():
			httpErr = types.Error(
				types.Http.C400().PreconditionFailed(),
				"Student was already marked as deleted",
				"Student ID: "+id,
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to delete student",
				err.Error(),
				"Student ID: "+id,
			)
		}

		return types.ResultErr[models.StudentDB](&httpErr)
	}

	return types.ResultOk(deletedStudent)
}
func (studentType) DeletePermanentByID(id string) types.Result[models.StudentDB] {
	logger.Warning("Permanently deleting student by ID: ", id)
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Student ID: "+id,
		)
		return types.ResultErr[models.StudentDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}, models.Filter.Deleted()} // Ensure the student is marked as deleted
	resultStudent := configs.DB.FindOne(filter, models.StudentDB{})
	if resultStudent.IsErr() {
		logger.Warning("Failed to find student for permanent deletion: ", resultStudent.Error())
		var httpErr types.HttpError

		switch resultStudent.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Student not found",
				"Student with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find student for permanent deletion",
				resultStudent.Error().Error(),
				"Student ID: "+id,
			)
		}

		return types.ResultErr[models.StudentDB](&httpErr)
	}

	err = configs.DB.PermanentDeleteOne(filter, models.StudentDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete student in MongoDB: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Student not found",
				"Student with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete student",
				err.Error(),
				"Student ID: "+id,
			)
		}

		return types.ResultErr[models.StudentDB](&httpErr)
	}

	return types.ResultOk(resultStudent.Value().(models.StudentDB))
}
func (studentType) DeletePermanentAll() types.Result[[]models.StudentDB] {
	filter := bson.D{models.Filter.Deleted()}

	resultStudents := configs.DB.FindAll(filter, models.StudentDB{})
	if resultStudents.IsErr() {
		logger.Warning("Failed to find student for permanent deletion: ", resultStudents.Error())
		var httpErr types.HttpError

		switch resultStudents.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Students not found",
				"No students marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find student for permanent deletion",
				resultStudents.Error().Error(),
			)
		}

		return types.ResultErr[[]models.StudentDB](&httpErr)
	}

	err := configs.DB.PermanentDeleteMany(filter, models.StudentDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete all students in MongoDB: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Students not found",
				"No students marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete all students",
				err.Error(),
			)
		}
		return types.ResultErr[[]models.StudentDB](&httpErr)
	}

	Students := utils.Map(resultStudents.Value(), models.InterfaceTo[models.StudentDB])
	return types.ResultOk(Students)
}
