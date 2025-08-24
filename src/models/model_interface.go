package models

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"

	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Later will see how to use this interface to
// simplify the code on the db interface layer.
//
// To implement this interface, a model must
// define a method TableName that returns the name of the database table
// and a method IsEmpty that checks if the object is empty.
// Yes, that's it, nothing else is required.
type DBModelInterface interface {
	TableName() string
	IsEmpty() bool
	//CreationDate() DBDateTime
	//UpdateDate() DBDateTime
	//DeleteDate() DBDateTime
}

type DBID = bson.ObjectID
type DBDateTime = time.Time

//func (m DBModelBase) CreationDate() DBDateTime {
//	return m.CreatedAt
//}
//
//func (m DBModelBase) UpdateDate() DBDateTime {
//	return m.UpdatedAt
//}
//
//func (m DBModelBase) DeleteDate() DBDateTime {
//	return m.DeletedAt
//}

// Deprecated: Use models.ID.ToPrimitive() instead.
func PrimitiveIDFrom(id any) (primitive.ObjectID, error) {
	logger.Deprecate(types.V("0.1.0"), types.V("0.1.3"), "Use models.ID.ToPrimitive() instead")
	return ID.ToPrimitive(id)
}

// Deprecated: Use models.ID.ToBson() instead.
func BsonIDFrom(id any) (bson.ObjectID, error) {
	logger.Deprecate(types.V("0.1.0"), types.V("0.1.3"), "Use models.ID.ToBson() instead")
	return ID.ToBson(id)
}

// Deprecated: Use models.ID.ToDB() instead.
func DBIDFrom(id any) (DBID, error) {
	logger.Deprecate(types.V("0.1.0"), types.V("0.1.3"), "Use models.ID.ToDB() instead")
	return ID.ToDB(id)
}

// Deprecated: Use models.ID.OmitEmpty() instead.
func OmitEmptyID(id string, result *DBID, idName string) bool {
	logger.Deprecate(types.V("0.1.0"), types.V("0.1.3"), "Use models.ID.OmitEmpty() instead")
	return ID.OmitEmpty(id, result, idName)
}

// Deprecated: Use models.ID.Ensure() instead.
func EnsureID(id string, result *DBID, idName string) bool {
	logger.Deprecate(types.V("0.1.0"), types.V("0.1.3"), "Use models.ID.Ensure() instead")
	return ID.Ensure(id, result, idName)
}

type iID struct{}

var ID iID

func (iID) ToPrimitive(id any) (primitive.ObjectID, error) {
	switch v := id.(type) {
	case string:
		return primitive.ObjectIDFromHex(v)
	case primitive.ObjectID:
		return v, nil
	case bson.ObjectID:
		return primitive.ObjectIDFromHex(v.Hex())
	default:
		return primitive.NilObjectID, fmt.Errorf("unsupported type for ToPrimitive: %T", id)
	}
}
func (iID) ToBson(id any) (bson.ObjectID, error) {
	switch v := id.(type) {
	case string:
		return bson.ObjectIDFromHex(v)
	case primitive.ObjectID:
		return bson.ObjectIDFromHex(v.Hex())
	case bson.ObjectID:
		return v, nil
	default:
		return bson.NilObjectID, fmt.Errorf("unsupported type for ToBson: %T", id)
	}
}

// Change this if you decide to change the ID type in the database
func (iID) ToDB(id any) (DBID, error) {
	return ID.ToBson(id)
}

func (iID) Ensure(id string, result *DBID, idName string) bool {
	if id == "" {
		logger.Warning("Missing required field:", idName)
		return false
	}

	idObj, err := ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert", idName, "to ObjectID:", err)
		return false
	}

	*result = idObj
	return true
}
func (iID) OmitEmpty(id string, result *DBID, idName string) bool {
	if id == "" {
		return true
	}

	idObj, err := ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert", idName, "to ObjectID:", err)
		return false
	}

	*result = idObj
	return true
}

// Deprecated: Use models.Time.Now() instead.
func TimeNow() DBDateTime {
	logger.Deprecate(types.V("0.1.0"), types.V("0.1.3"), "Use models.Time.Now() instead")
	return Time.Now()
}

// Deprecated: Use models.Time.Zero() instead.
func TimeZero() DBDateTime {
	logger.Deprecate(types.V("0.1.0"), types.V("0.1.3"), "Use models.Time.Zero() instead")
	return Time.Zero()
}

type iTime struct{}

// Unified time handling for models.
// If you want to change the time type used, first update the Time struct "DBDateTime".
// Then update the methods in the iTime struct accordingly.
var Time iTime

func (iTime) Now() DBDateTime {
	return time.Now()
}

// This is not necessarily the unix epoch or a 0 date or any other starting point, but rather default "empty" value used in the application.
func (iTime) Zero() DBDateTime {
	return time.Time{}
}

type iFilters struct{}

var Filter iFilters

type FilterType = bson.E

func (iFilters) ID(id bson.ObjectID) FilterType {
	return FilterType{Key: "_id", Value: id} // Filter by ID
}
func (iFilters) IDOf(idName string, id bson.ObjectID) FilterType {
	return FilterType{Key: "id_" + idName, Value: id} // Filter by ID with custom field name
}
func (iFilters) NotDeleted() FilterType {
	return FilterType{Key: "deleted_at", Value: Time.Zero()} // Filter to exclude deleted records
}
func (iFilters) Deleted() FilterType {
	return FilterType{Key: "deleted_at", Value: bson.M{"$ne": Time.Zero()}} // Filter to include deleted records
}

type iUpdate struct{}

var Update iUpdate

type UpdateType = bson.M

func (iUpdate) Delete() UpdateType {
	return UpdateType{"$set": bson.M{"deleted_at": Time.Now()}} // Soft delete
}

func InterfaceTo[T DBModelInterface](a DBModelInterface) T {
	return a.(T)
}
