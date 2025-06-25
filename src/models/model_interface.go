package models

import (
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
}

type DBID = bson.ObjectID
type DBDateTime = time.Time

func PrimitiveIDFrom(id any) (primitive.ObjectID, error) {
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
	return BsonIDFrom(id)
}

// If decide to change the time type, you can only change it here
func TimeNow() DBDateTime {
	return time.Now()
}

// If decide to change the time type, you can only change it here
func TimeZero() DBDateTime {
	return time.Time{}
}
