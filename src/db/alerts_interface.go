package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
)

type alertType struct{}

var Alert alertType

func (alertType) Create(u models.AlertCreate) types.Result[models.AlertDB] {
	resultAlertDB := u.ToInsert()
	if resultAlertDB.IsErr() {
		logger.Error("Error converting alert for insertion: ", resultAlertDB.Error())
		return types.ResultErr[models.AlertDB](resultAlertDB.Error())
	}
	alertDB := resultAlertDB.Value()
	resultID := configs.DB.InsertOne(alertDB)

	if resultID.IsErr() {
		logger.Error("Error inserting alert: ", resultID.Error())
		return types.ResultErr[models.AlertDB](resultID.Error())
	}

	alertDB.ID = resultID.Value()
	return types.ResultOk(alertDB)
}

func (alertType) GetByID(id string, filter models.FilterObject) types.Result[models.AlertDB] {
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Alert ID: "+id,
		)
		return types.ResultErr[models.AlertDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultAlert := configs.DB.FindOne(filter, models.AlertDB{})
	if resultAlert.IsErr() {
		logger.Warning("Failed to get alert by ID: ", resultAlert.Error())
		var httpErr types.HttpError

		switch resultAlert.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Alert not found",
				"Alert with ID "+id+" not found",
			)

		default:
			logger.Error("Failed to get alert by ID: ", resultAlert.Error())
			httpErr = types.ErrorInternal(
				"Failed to retrieve alert",
				resultAlert.Error().Error(),
				"Alert ID: "+id,
			)
		}

		return types.ResultErr[models.AlertDB](&httpErr)
	}

	return types.ResultOk(resultAlert.Value().(models.AlertDB))
}
func (alertType) GetAll(filter models.FilterObject) types.Result[[]models.AlertDB] {
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultSpecialities := configs.DB.FindAll(filter, models.AlertDB{})
	if resultSpecialities.IsErr() {
		logger.Warning("Failed to get all specialities from database:", resultSpecialities.Error())
		var httpErr types.HttpError

		switch resultSpecialities.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Specialities not found",
				"No specialities found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve specialities",
				resultSpecialities.Error().Error(),
			)
		}
		return types.ResultErr[[]models.AlertDB](&httpErr)
	}

	specialities := utils.Map(resultSpecialities.Value(), models.InterfaceTo[models.AlertDB])
	logger.Debug("Retrieved", len(specialities), "specialities from database")
	return types.ResultOk(specialities)
}

func (alertType) UpdateByID(id string, alert models.AlertCreate, filter models.FilterObject) types.Result[models.AlertDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Alert ID: "+id,
		)
		return types.ResultErr[models.AlertDB](&httpErr)
	}

	resultAlertUpdate := alert.ToUpdate()
	if resultAlertUpdate.IsErr() {
		logger.Warning("Error converting alert to DB model: ", resultAlertUpdate.Error())
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid alert data",
			resultAlertUpdate.Error().Error(),
			"Alert ID: "+id,
		)
		return types.ResultErr[models.AlertDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.UpdateOne(filter, resultAlertUpdate.Value())
	if err != nil {
		logger.Warning("Failed to update alert in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Alert not found",
				"Alert with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to update alert",
				err.Error(),
				"Alert ID: "+id,
			)
		}
		return types.ResultErr[models.AlertDB](&httpErr)
	}

	resultAlertDB := configs.DB.FindOne(filter, models.AlertDB{})
	if resultAlertDB.IsErr() {
		logger.Warning("Failed to retrieve updated alert: ", resultAlertDB.Error())
		var httpErr types.HttpError

		switch resultAlertDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Alert not found",
				"Alert with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated alert",
				resultAlertDB.Error().Error(),
				"Alert ID: "+id,
			)
		}
		return types.ResultErr[models.AlertDB](&httpErr)
	}

	return types.ResultOk(resultAlertDB.Value().(models.AlertDB))
}

func (alertType) PatchByID(id string, alert models.AlertCreate, filter models.FilterObject) types.Result[models.AlertDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Alert ID: "+id,
		)
		return types.ResultErr[models.AlertDB](&httpErr)
	}

	logger.Debug("Alert patch:", alert)
	resultAlertPatch := alert.ToUpdate()
	if resultAlertPatch.IsErr() {
		logger.Warning("Error converting alert to DB model: ", resultAlertPatch.Error())
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid alert data",
			resultAlertPatch.Error().Error(),
			"Alert ID: "+id,
		)
		return types.ResultErr[models.AlertDB](&httpErr)
	}
	logger.Debug("Alert db:", resultAlertPatch.Value())

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.PatchOne(filter, resultAlertPatch.Value())
	if err != nil {
		logger.Warning("Failed to patch alert in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Alert not found",
				"Alert with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to patch alert",
				err.Error(),
				"Alert ID: "+id,
			)
		}
		return types.ResultErr[models.AlertDB](&httpErr)
	}

	resultAlertDB := configs.DB.FindOne(filter, models.AlertDB{})
	if resultAlertDB.IsErr() {
		logger.Warning("Failed to retrieve updated alert: ", resultAlertDB.Error())
		var httpErr types.HttpError

		switch resultAlertDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Alert not found",
				"Alert with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated alert",
				resultAlertDB.Error().Error(),
				"Alert ID: "+id,
			)
		}
		return types.ResultErr[models.AlertDB](&httpErr)
	}

	return types.ResultOk(resultAlertDB.Value().(models.AlertDB))
}

func (alertType) DeleteByID(id string, filter models.FilterObject) types.Result[models.AlertDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Alert ID: "+id,
		)
		return types.ResultErr[models.AlertDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	resultAlert := configs.DB.FindOne(filter, models.AlertDB{})
	if resultAlert.IsErr() {
		logger.Warning("Failed to retrieve alert: ", resultAlert.Error())
		var httpErr types.HttpError

		switch resultAlert.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Alert not found",
				"Alert with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve alert",
				resultAlert.Error().Error(),
				"Alert ID: "+id,
			)
		}

		return types.ResultErr[models.AlertDB](&httpErr)
	}

	deletedAlert := resultAlert.Value().(models.AlertDB)
	err = configs.DB.SoftDeleteOne(filter, deletedAlert)
	if err != nil {
		logger.Warning("Failed to delete alert in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Alert not found",
				"Alert with ID "+id+" not found",
			)

		case configs.DBErr.NotModified():
			httpErr = types.Error(
				types.Http.C400().PreconditionFailed(),
				"Alert was already marked as deleted",
				"Alert ID: "+id,
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to delete alert",
				err.Error(),
				"Alert ID: "+id,
			)
		}

		return types.ResultErr[models.AlertDB](&httpErr)
	}

	return types.ResultOk(deletedAlert)
}
func (alertType) DeletePermanentByID(id string, filter models.FilterObject) types.Result[models.AlertDB] {
	logger.Warning("Permanently deleting alert by ID: ", id)
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Alert ID: "+id,
		)
		return types.ResultErr[models.AlertDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())
	resultAlert := configs.DB.FindOne(filter, models.AlertDB{})
	if resultAlert.IsErr() {
		logger.Warning("Failed to find alert for permanent deletion: ", resultAlert.Error())
		var httpErr types.HttpError

		switch resultAlert.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Alert not found",
				"Alert with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find alert for permanent deletion",
				resultAlert.Error().Error(),
				"Alert ID: "+id,
			)
		}

		return types.ResultErr[models.AlertDB](&httpErr)
	}

	err = configs.DB.PermanentDeleteOne(filter, models.AlertDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete alert in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Alert not found",
				"Alert with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete alert",
				err.Error(),
				"Alert ID: "+id,
			)
		}

		return types.ResultErr[models.AlertDB](&httpErr)
	}

	return types.ResultOk(resultAlert.Value().(models.AlertDB))
}
func (alertType) DeletePermanentAll(filter models.FilterObject) types.Result[[]models.AlertDB] {
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())

	resultAlerts := configs.DB.FindAll(filter, models.AlertDB{})
	if resultAlerts.IsErr() {
		logger.Warning("Failed to find alert for permanent deletion: ", resultAlerts.Error())
		var httpErr types.HttpError

		switch resultAlerts.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Alerts not found",
				"No alerts marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find alert for permanent deletion",
				resultAlerts.Error().Error(),
			)
		}

		return types.ResultErr[[]models.AlertDB](&httpErr)
	}

	err := configs.DB.PermanentDeleteMany(filter, models.AlertDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete all alerts in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Alerts not found",
				"No alerts marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete all alerts",
				err.Error(),
			)
		}
		return types.ResultErr[[]models.AlertDB](&httpErr)
	}

	alerts := utils.Map(resultAlerts.Value(), models.InterfaceTo[models.AlertDB])
	return types.ResultOk(alerts)
}
