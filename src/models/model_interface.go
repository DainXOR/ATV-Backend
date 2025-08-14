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
	SetID(id any) error
	GetID() DBID
}
type DBModelBase struct {
	ID DBID `json:"_id,omitempty" bson:"_id,omitempty"`
}

func (m *DBModelBase) SetID(id any) error {
	var err error
	m.ID, err = ID.ToDB(id)
	return err
}
func (m *DBModelBase) GetID() string {
	return m.ID.Hex()
}

type DBID = bson.ObjectID
type DBDateTime = time.Time

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
	return time.Now()
}

// Deprecated: Use models.Time.Zero() instead.
func TimeZero() DBDateTime {
	logger.Deprecate(types.V("0.1.0"), types.V("0.1.3"), "Use models.Time.Zero() instead")
	return Time.Zero()
}

type iTime struct{}

var Time iTime

// If decide to change the time type, you can only change it here
func (iTime) Now() DBDateTime {
	return time.Now()
} // If decide to change the time type, you can only change it here
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
