package dao

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
)

type contactReasonType struct{}

var ContactReason contactReasonType

func (contactReasonType) Create(u models.ContactReasonCreate) types.Result[models.ContactReasonDB] {
	contactReasonDB := u.ToInsert()
	resultID := configs.DB.InsertOne(contactReasonDB)

	if resultID.IsErr() {
		logger.Warning("Error inserting contact reason: ", resultID.Error())
		return types.ResultErr[models.ContactReasonDB](resultID.Error())
	}

	contactReasonDB.ID = resultID.Value()
	return types.ResultOk(contactReasonDB)
}

func (contactReasonType) GetByID(id string, filter models.FilterObject) types.Result[models.ContactReasonDB] {
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"contact reason ID: "+id,
		)
		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultContactReason := configs.DB.FindOne(filter, models.ContactReasonDB{})
	if resultContactReason.IsErr() {
		logger.Warning("Failed to get contact reason by ID: ", resultContactReason.Error())
		var httpErr types.HttpError

		switch resultContactReason.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Contact reason not found",
				"Contact reason with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve contact reason",
				resultContactReason.Error().Error(),
				"Contact reason ID: "+id,
			)
		}

		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	return types.ResultOk(resultContactReason.Value().(models.ContactReasonDB))
}
func (contactReasonType) GetAll(filter models.FilterObject) types.Result[[]models.ContactReasonDB] {
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultContactReasons := configs.DB.FindAll(filter, models.ContactReasonDB{})
	if resultContactReasons.IsErr() {
		logger.Warning("Failed to get all contact reasons from database:", resultContactReasons.Error())
		var httpErr types.HttpError

		switch resultContactReasons.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Contact reasons not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve contact reasons",
				resultContactReasons.Error().Error(),
			)

		}

		return types.ResultErr[[]models.ContactReasonDB](&httpErr)
	}

	contactReasons := utils.Map(resultContactReasons.Value(), models.InterfaceTo[models.ContactReasonDB])
	logger.Debug("Retrieved", len(contactReasons), "contact reasons from database")
	return types.ResultOk(contactReasons)
}

func (contactReasonType) UpdateByID(id string, contactReason models.ContactReasonCreate, filter models.FilterObject) types.Result[models.ContactReasonDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Contact reason ID: "+id,
		)
		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	contactReasonUpdate := contactReason.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.UpdateOne(filter, contactReasonUpdate)
	if err != nil {
		logger.Warning("Failed to update contact reason in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Contact reason not found",
				"Contact reason with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to update contact reason",
				err.Error(),
				"Contact reason ID: "+id,
			)
		}
		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	resultContactReasonDB := configs.DB.FindOne(filter, models.ContactReasonDB{})
	if resultContactReasonDB.IsErr() {
		logger.Warning("Failed to retrieve updated contact reason: ", resultContactReasonDB.Error())
		var httpErr types.HttpError

		switch resultContactReasonDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"contact reason not found",
				"contact reason with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated contact reason",
				resultContactReasonDB.Error().Error(),
				"Contact reason ID: "+id,
			)
		}
		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	return types.ResultOk(resultContactReasonDB.Value().(models.ContactReasonDB))
}

func (contactReasonType) PatchByID(id string, contactReason models.ContactReasonCreate, filter models.FilterObject) types.Result[models.ContactReasonDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Contact reason ID: "+id,
		)
		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	contactReasonPatch := contactReason.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.PatchOne(filter, contactReasonPatch)
	if err != nil {
		logger.Warning("Failed to patch contact reason in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Contact reason not found",
				"Contact reason with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to patch contact reason",
				err.Error(),
				"Contact reason ID: "+id,
			)
		}
		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	resultContactReasonDB := configs.DB.FindOne(filter, models.ContactReasonDB{})
	if resultContactReasonDB.IsErr() {
		logger.Warning("Failed to retrieve updated contact reason: ", resultContactReasonDB.Error())
		var httpErr types.HttpError

		switch resultContactReasonDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"contact reason not found",
				"Contact reason with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated contact reason",
				resultContactReasonDB.Error().Error(),
				"Contact reason ID: "+id,
			)
		}
		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	return types.ResultOk(resultContactReasonDB.Value().(models.ContactReasonDB))
}

func (contactReasonType) DeleteByID(id string, filter models.FilterObject) types.Result[models.ContactReasonDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Contact reason ID: "+id,
		)
		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	resultContactReason := configs.DB.FindOne(filter, models.ContactReasonDB{})
	if resultContactReason.IsErr() {
		logger.Warning("Failed to retrieve contact reason: ", resultContactReason.Error())
		var httpErr types.HttpError

		switch resultContactReason.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Contact reason not found",
				"Contact reason with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve contact reason",
				resultContactReason.Error().Error(),
				"Contact reason ID: "+id,
			)
		}

		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	deletedContactReason := resultContactReason.Value().(models.ContactReasonDB)
	err = configs.DB.SoftDeleteOne(filter, deletedContactReason)
	if err != nil {
		logger.Warning("Failed to delete contact reason in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Contact reason not found",
				"Contact reason with ID "+id+" not found",
			)

		case configs.DBErr.NotModified():
			httpErr = types.Error(
				types.Http.C400().PreconditionFailed(),
				"Contact reason was already marked as deleted",
				"Contact reason ID: "+id,
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to delete contact reason",
				err.Error(),
				"Contact reason ID: "+id,
			)
		}

		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	return types.ResultOk(deletedContactReason)
}
func (contactReasonType) DeletePermanentByID(id string, filter models.FilterObject) types.Result[models.ContactReasonDB] {
	logger.Warning("Permanently deleting contact reason by ID: ", id)
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Contact reason ID: "+id,
		)
		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())
	resultContactReason := configs.DB.FindOne(filter, models.ContactReasonDB{})
	if resultContactReason.IsErr() {
		logger.Warning("Failed to find contact reason for permanent deletion: ", resultContactReason.Error())
		var httpErr types.HttpError

		switch resultContactReason.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Contact reason not found",
				"Contact reason with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find contact reason for permanent deletion",
				resultContactReason.Error().Error(),
				"Contact reason ID: "+id,
			)
		}

		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	err = configs.DB.PermanentDeleteOne(filter, models.ContactReasonDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete contact reason in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Contact reason not found",
				"Contact reason with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete contact reason",
				err.Error(),
				"Contact reason ID: "+id,
			)
		}

		return types.ResultErr[models.ContactReasonDB](&httpErr)
	}

	return types.ResultOk(resultContactReason.Value().(models.ContactReasonDB))
}
func (contactReasonType) DeletePermanentAll(filter models.FilterObject) types.Result[[]models.ContactReasonDB] {
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())

	resultContactReasons := configs.DB.FindAll(filter, models.ContactReasonDB{})
	if resultContactReasons.IsErr() {
		logger.Warning("Failed to find contact reasons for permanent deletion: ", resultContactReasons.Error())
		var httpErr types.HttpError

		switch resultContactReasons.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Contact reasons not found",
				"No contact reasons marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find contact reasons for permanent deletion",
				resultContactReasons.Error().Error(),
			)
		}

		return types.ResultErr[[]models.ContactReasonDB](&httpErr)
	}

	err := configs.DB.PermanentDeleteMany(filter, models.ContactReasonDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete all contact reasons in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Contact reasons not found",
				"No contact reasons marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete all contact reasons",
				err.Error(),
			)
		}
		return types.ResultErr[[]models.ContactReasonDB](&httpErr)
	}

	contactReasons := utils.Map(resultContactReasons.Value(), models.InterfaceTo[models.ContactReasonDB])
	return types.ResultOk(contactReasons)
}
