package dao

import (
	"dainxor/atv/models"
	"dainxor/atv/types"
)

type DAOInterface[M models.DBModelInterface, C any] interface {
	Create(m C) types.Result[M]
	// Missing CreateAll

	GetByID(id string, filter models.FilterObject) types.Result[M]
	GetAll(filter models.FilterObject) types.Result[[]M]

	UpdateByID(id string, createBody C, filter models.FilterObject) types.Result[M]
	// Missing UpdateAll

	PatchByID(id string, createBody C, filter models.FilterObject) types.Result[M]
	// Missing PatchAll

	DeleteByID(id string, filter models.FilterObject) types.Result[M]
	// Missing DeleteAll

	// Missing DeletePermanentByID
	// Missing DeletePermanentAll
}
