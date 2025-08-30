package models

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"errors"
	"net/url"
	"strings"

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

func (iTime) Format() string {
	return "02-01-2006T15:04"
}

func (iTime) Parse(dateStr string, format string) (DBDateTime, error) {
	parsedTime, err := time.Parse(format, dateStr)
	if err != nil {
		return Time.Zero(), err
	}
	return parsedTime, nil
}

type iFilters struct {
}

var Filter iFilters

type FilterPart = bson.E
type FilterObject = bson.D

var timeValues = []string{
	"created",
	"updated",
	"deleted",
}

func (iFilters) dateFormatString() string {
	return "02-01-2006T15:04"
}
func (iFilters) dateFormatLetters() string {
	dateString := strings.Replace(Filter.dateFormatString(), "02", "dd", 1)
	dateString = strings.Replace(dateString, "01", "MM", 1)
	dateString = strings.Replace(dateString, "04", "mm", 1)
	dateString = strings.Replace(dateString, "06", "yy", 1)
	dateString = strings.Replace(dateString, "20yy", "yyyy", 1)

	if strings.Contains(dateString, "15") {
		dateString = strings.Replace(dateString, "15", "HH", 1)
	} else {
		dateString = strings.Replace(dateString, "03", "HH", 1)
	}

	return dateString
}
func (iFilters) ParseDate(dateStr string) (DBDateTime, error) {
	return Time.Parse(dateStr, Filter.dateFormatString())
}

// If you want to change how the filter key-value pairs are created, modify this function.
func (iFilters) parse(name string, values []string) (FilterPart, error) {
	logger.Debugf("Filter for %s: %#v", name, values)

	if len(values) == 1 {
		return FilterPart{Key: name, Value: values[0]}, nil
	}

	dateFirst, err := Filter.ParseDate(values[0])
	if err != nil {
		logger.Warning("Failed to parse date:", err)
		logger.Warningf("Format should be %s (e.g., 29-06-2025T15:32)", Filter.dateFormatLetters())
		return FilterPart{}, err
	}

	dateLast, err := Filter.ParseDate(values[1])
	if err != nil {
		logger.Warning("Failed to parse date:", err)
		logger.Warningf("Format should be %s (e.g., 29-06-2025T15:32)", Filter.dateFormatLetters())
		return FilterPart{}, err
	}

	if dateFirst.Compare(dateLast) > 0 {
		auxDate := dateFirst
		dateFirst = dateLast
		dateLast = auxDate
	}

	logger.Lava(types.V("0.2.0"), "Using mongoDB specific syntax for date ranges.")

	// { "created_at": { "$gte": ISODate("2025-06-29T01:32:45.401+00:00"), "$lte": ISODate("2025-07-29T01:32:45.401+00:00") } }
	// { "created_at": { "$gte": ISODate("2025-06-29T00:00:00Z"), "$lte": ISODate("2025-07-29T23:59:59Z") } }
	if utils.Any(timeValues, func(s string) bool { return name == s }) {
		return FilterPart{Key: name + "_at", Value: bson.M{"$gte": dateFirst, "$lte": dateLast}}, nil

	} else if name == "date" {
		return FilterPart{}, errors.New("unsupported filter")
		// return FilterPart{Key: name, Value: bson.M{"$gte": dateFirst, "$lte": dateLast}}, nil
	}

	return FilterPart{}, errors.New("unsupported filter")
}

// If you want to filter out certain query parameters, modify this function.
func (iFilters) skip(_ string, _ []string) bool {
	return false
}

// Returns a filter object from the given query parameters.
// If you want to change how the filter is created or the type, modify this function.
// Make sure to modify the corresponding aliases used if needed.
func (iFilters) Create(queryParams url.Values) FilterObject {
	var filter FilterObject
	for key, vals := range queryParams {
		if Filter.skip(key, vals) {
			continue
		}

		part, err := Filter.parse(key, vals)
		if err != nil {
			logger.Warning("Failed to parse filter part:", err)
			continue
		}
		filter = append(filter, part)
	}
	return filter
}
func (iFilters) Add(filter FilterObject, name string, value []string) FilterObject {
	part, err := Filter.parse(name, value)
	if err != nil {
		return filter
	}
	filter = append(filter, part)
	return filter
}
func (iFilters) AddPart(filter FilterObject, part FilterPart) FilterObject {
	filter = append(filter, part)
	return filter
}
func (iFilters) Merge(filter1, filter2 FilterObject) FilterObject {
	return append(filter1, filter2...)
}

func (iFilters) Empty() FilterObject {
	return FilterObject{}
}

func (iFilters) Of(name string, value any) FilterPart {
	return FilterPart{Key: name, Value: value} // Generic filter by field name and value
}
func (iFilters) ID(id bson.ObjectID) FilterPart {
	return FilterPart{Key: "_id", Value: id} // Filter by ID
}
func (iFilters) IDOf(idName string, id bson.ObjectID) FilterPart {
	return FilterPart{Key: "id_" + idName, Value: id} // Filter by ID with custom field name
}
func (iFilters) NotDeleted() FilterPart {
	return FilterPart{Key: "deleted_at", Value: Time.Zero()} // Filter to exclude deleted records
}
func (iFilters) Deleted() FilterPart {
	return FilterPart{Key: "deleted_at", Value: bson.M{"$ne": Time.Zero()}} // Filter to include deleted records
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
func ToInterface[T DBModelInterface](a T) DBModelInterface {
	return a
}
