package dao

import (
	"dainxor/atv/models"
	"dainxor/atv/types"
	"errors"
)

type DAOInterface[C, R any, M models.DBModelInterface[C, R]] interface {
	Create(body C) types.Result[M]
	CreateAll(bodies []C) types.Result[[]M]

	GetByID(id string, filter models.FilterObject) types.Result[M]
	GetAll(filter models.FilterObject) types.Result[[]M]

	UpdateByID(id string, updateBody C, filter models.FilterObject) types.Result[M]
	UpdateAll(updateBody C, filter models.FilterObject) types.Result[[]M]

	PatchByID(id string, createBody C, filter models.FilterObject) types.Result[M]
	PatchAll(updateBody C, filter models.FilterObject) types.Result[[]M]

	DeleteByID(id string, filter models.FilterObject) types.Result[M]
	DeleteAll(filter models.FilterObject) types.Result[[]M]

	DeletePermanentByID(id string, filter models.FilterObject) types.Result[M]
	DeletePermanentAll(filter models.FilterObject) types.Result[[]M]
}

type daoErrorsNS struct{}

var Error daoErrorsNS

var (
	_notImplemented = errors.New("Not implemented")

	_notInserted = errors.New("Error during insert")
)

func (daoErrorsNS) NotImplemented() error {
	return _notImplemented
}

func (daoErrorsNS) Notinserted() error {
	return _notInserted
}
