package configs

import (
	"dainxor/atv/models"
	"dainxor/atv/types"
	"errors"
)

type dbError struct{}

var DBErr dbError

func (dbError) NotFound() error {
	return errors.New("document not found")
}
func (dbError) NotModified() error {
	return errors.New("document not modified")
}
func (dbError) NotDeleted() error {
	return errors.New("document not deleted")
}
func (dbError) AlreadyExists() error {
	return errors.New("document already exists")
}
func (dbError) InvalidInput() error {
	return errors.New("invalid input")
}
func (dbError) Internal() error {
	return errors.New("internal server error")
}
func (dbError) NotImplemented() error {
	return errors.New("not implemented")
}

type InterfaceDBAccessor interface {
	Connect(dbName, conectionString string) error
	Disconnect() error
	Migrate(models ...models.DBModelInterface) error

	CreateFilter(filter ...types.SPair[string]) any

	InsertOne(element models.DBModelInterface) types.Result[models.DBID]
	InsertMany(elements ...models.DBModelInterface) types.Result[[]models.DBID]

	FindOne(filter any, model models.DBModelInterface) types.Result[models.DBModelInterface]
	FindMany(filter any, model models.DBModelInterface) types.Result[[]models.DBModelInterface]

	UpdateOne(filter any, update models.DBModelInterface) error
	UpdateMany(filter any, update models.DBModelInterface) error

	PatchOne(filter any, update models.DBModelInterface) error
	PatchMany(filter any, update models.DBModelInterface) error

	SoftDeleteOne(filter any, model models.DBModelInterface) error
	SoftDeleteMany(filter any, model models.DBModelInterface) error

	PermanentDeleteOne(filter any, model models.DBModelInterface) error
	PermanentDeleteMany(filter any, model models.DBModelInterface) error
}
