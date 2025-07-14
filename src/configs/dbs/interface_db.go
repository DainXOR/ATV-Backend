package db

import (
	"dainxor/atv/models"
	"dainxor/atv/types"
	"errors"
)

var (
	// Default errors for database operations
	ErrNotFound      = errors.New("document not found")
	ErrNotModified   = errors.New("document not modified")
	ErrNotDeleted    = errors.New("document not deleted")
	ErrAlreadyExists = errors.New("document already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInternal      = errors.New("internal server error")
)

type InterfaceDB interface {
	CreateFilter(filter []types.SPair[string]) any
	CreateUpdator(update models.DBModelInterface) any

	CreateOne(element models.DBModelInterface) types.Result[models.DBModelInterface]
	CreateMany(elements any) types.Result[any]

	GetOne(filter any, result models.DBModelInterface) types.Result[models.DBModelInterface]
	GetAll(filter any, result any) types.Result[any]

	UpdateOne(filter any, update models.DBModelInterface, result models.DBModelInterface) types.Result[models.DBModelInterface]

	PatchOne(filter any, update models.DBModelInterface, result models.DBModelInterface) types.Result[models.DBModelInterface]

	DeleteOne(filter any, result models.DBModelInterface) types.Result[models.DBModelInterface]
	DeleteAll(filter any, result any) types.Result[any]
}
