package dao

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
)

func genericGetAllOf[M models.DBModelInterface](model M) func(models.FilterObject) types.Result[[]M] {
	return func(filter models.FilterObject) types.Result[[]M] {
		filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

		resultObjects := configs.DB.FindAll(filter, model)
		if resultObjects.IsErr() {
			logger.Warning("Failed to get all objects: ", resultObjects.Error())
			return types.ResultErr[[]M](resultObjects.Error())
		}

		objectsDB := utils.Map(resultObjects.Value(), models.InterfaceTo[M])
		logger.Debug("Retrieved", len(objectsDB), "objects from db")
		return types.ResultOk(objectsDB)
	}
}
