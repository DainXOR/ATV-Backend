package models

import (
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

type DBID = primitive.ObjectID
type DBDateTime = primitive.DateTime

func IDFrom(id any) (DBID, error) {
	switch v := id.(type) {
	case string:
		return primitive.ObjectIDFromHex(v)
	case DBID:
		return v, nil
	case bson.ObjectID:
		return primitive.ObjectIDFromHex(v.Hex())
	default:
		return primitive.NilObjectID, nil
	}
}
func TimeNow() DBDateTime {
	return primitive.NewDateTimeFromTime(time.Now())
}
