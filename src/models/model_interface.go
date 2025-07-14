package models

import (
	"dainxor/atv/logger"

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
// Yes, that's it, nothing else is required.
type DBModelInterface interface {
	TableName() string
	IsEmpty() bool
}

type DBID = bson.ObjectID
type DBDateTime = time.Time

func PrimitiveIDFrom(id any) (primitive.ObjectID, error) {
	logger.Deprecate("0.1.0", "0.1.3", "Use models.ID.ToPrimitive() instead")
	switch v := id.(type) {
	case string:
		return primitive.ObjectIDFromHex(v)
	case primitive.ObjectID:
		return v, nil
	case bson.ObjectID:
		return primitive.ObjectIDFromHex(v.Hex())
	default:
		return primitive.NilObjectID, fmt.Errorf("unsupported type for PrimitiveIDFrom: %T", id)
	}
}
func BsonIDFrom(id any) (bson.ObjectID, error) {
	logger.Deprecate("0.1.0", "0.1.3", "Use models.ID.ToBson() instead")
	switch v := id.(type) {
	case string:
		return bson.ObjectIDFromHex(v)
	case primitive.ObjectID:
		return bson.ObjectIDFromHex(v.Hex())
	case bson.ObjectID:
		return v, nil
	default:
		return bson.NilObjectID, fmt.Errorf("unsupported type for BSONIDFrom: %T", id)
	}
}

// Change this if you decide to change the ID type in the database
func DBIDFrom(id any) (DBID, error) {
	logger.Deprecate("0.1.0", "0.1.3", "Use models.ID.ToDBID() instead")
	return ID.ToBson(id)
}

func OmitEmptyID(id string, result *DBID, idName string) bool {
	logger.Deprecate("0.1.0", "0.1.3", "Use models.ID.OmitEmpty() instead")
	if id != "" {
		return ID.Ensure(id, result, idName)
	}
	return true
}
func EnsureID(id string, result *DBID, idName string) bool {
	logger.Deprecate("0.1.0", "0.1.3", "Use models.ID.Ensure() instead")
	if id == "" {
		logger.Warning("Missing required field:", idName)
		return false
	}

	idObj, err := DBIDFrom(id)
	if err != nil {
		logger.Warning("Failed to convert", idName, "to ObjectID:", err)
		return false
	}

	*result = idObj
	return true
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
		return primitive.NilObjectID, fmt.Errorf("unsupported type for PrimitiveIDFrom: %T", id)
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
		return bson.NilObjectID, fmt.Errorf("unsupported type for BSONIDFrom: %T", id)
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
	idObj, err := ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert", idName, "to ObjectID:", err)
		return false
	}

	*result = idObj
	return true
}

// If decide to change the time type, you can only change it here
func TimeNow() DBDateTime {
	logger.Deprecate("0.1.0", "0.1.3", "Use models.Time.Now() directly instead")
	return time.Now()
}

// If decide to change the time type, you can only change it here
func TimeZero() DBDateTime {
	logger.Deprecate("0.1.0", "0.1.3", "Use models.Time.Zero() directly instead")
	return time.Time{}
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

func (iFilters) ID(id bson.ObjectID) bson.E {
	return bson.E{Key: "_id", Value: id} // Filter by ID
}
func (iFilters) IDOf(idName string, id bson.ObjectID) bson.E {
	return bson.E{Key: "id_" + idName, Value: id} // Filter by ID with custom field name
}
func (iFilters) NotDeleted() bson.E {
	return bson.E{Key: "deleted_at", Value: Time.Zero()} // Filter to exclude deleted records
}
func (iFilters) Deleted() bson.E {
	return bson.E{Key: "deleted_at", Value: bson.M{"$ne": Time.Zero()}} // Filter to include deleted records
}
