package models

import "go.mongodb.org/mongo-driver/v2/bson"

type testHelperNS struct{}

var Test testHelperNS

// GenerateObjectID returns a DBID for testing purposes. It is unexported and only for use in model tests.
func (testHelperNS) GenerateObjectID() bson.ObjectID {
	return bson.NewObjectID()
}
