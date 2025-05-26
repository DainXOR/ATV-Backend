package models

import (
	"dainxor/atv/utils"
	"reflect"
	"strings"
	"time"

	"slices"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// UserDBGorm represents the database model for a user
type UserDBGorm struct {
	gorm.Model              // Embedded gorm.Model struct to include ID, CreatedAt, UpdatedAt, DeletedAt
	IDNumber         string `json:"id_number" gorm:"unique;not null"`
	FirstName        string `json:"first_name" gorm:"not null"`
	LastName         string `json:"last_name" gorm:"not null"`
	PersonalEmail    string `json:"email" gorm:"unique;not null"`
	InstitutionEmail string `json:"institution_email" gorm:"unique;not null"`
	ResidenceAddress string `json:"residence_address" gorm:"not null"`
	Semester         uint   `json:"semester" gorm:"not null"`
	UniversityID     string `json:"university_id" gorm:"not null"`
	PhoneNumber      string `json:"phone_number"`
	// gorm tags are used to specify constraints and properties for the fields
	// Unique forces the field to be unique in the database for each record
	// NotNull forces the field to be not null in the database for each record
	// JSON tags are used to specify the JSON key names for the fields inside de db and in JSON schema
}
type UserDBMongo struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IDNumber         string             `json:"id_number,omitempty" bson:"id_number,omitempty"`
	FirstName        string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	PersonalEmail    string             `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string             `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	ResidenceAddress string             `json:"residence_address,omitempty" bson:"residence_address,omitempty"`
	Semester         uint               `json:"semester,omitempty" bson:"semester,omitempty"`
	UniversityID     string             `json:"university_id,omitempty" bson:"university_id,omitempty"`
	PhoneNumber      string             `json:"phone_number" bson:"phone_number"`
	CreatedAt        time.Time          `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt        time.Time          `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
	DeletedAt        *time.Time         `json:"deleted_at,omitzero" bson:"deleted_at,omitzero"`
}

// UserCreate represents the request body for creating a new user or updating an existing user
// It is used to validate the input data before creating or updating a user in the database
type UserCreate struct {
	IDNumber         string `json:"id_number" gorm:"unique;not null"`
	FirstName        string `json:"first_name" gorm:"not null"`
	LastName         string `json:"last_name" gorm:"not null"`
	PersonalEmail    string `json:"email" gorm:"unique;not null"`
	InstitutionEmail string `json:"institution_email" gorm:"unique;not null"`
	ResidenceAddress string `json:"residence_address" gorm:"not null"`
	Semester         uint   `json:"semester" gorm:"not null"`
	UniversityID     string `json:"university_id" gorm:"not null"`
	PhoneNumber      string `json:"phone_number"`
}

// UserResponse represents the response body for a user
// It is used to format the data returned to the client after a user is created or retrieved
// It includes the ID, created_at, and updated_at fields
type UserResponse struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDNumber         string    `json:"id_number" gorm:"unique;not null"`
	FirstName        string    `json:"first_name" gorm:"not null"`
	LastName         string    `json:"last_name" gorm:"not null"`
	PersonalEmail    string    `json:"email" gorm:"unique;not null"`
	InstitutionEmail string    `json:"institution_email" gorm:"unique;not null"`
	ResidenceAddress string    `json:"residence_address" gorm:"not null"`
	Semester         uint      `json:"semester" gorm:"not null"`
	UniversityID     string    `json:"university_id" gorm:"not null"`
	PhoneNumber      string    `json:"phone_number"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// UserCreateToDB converts a UserCreate struct to a UserDB struct
// This is used to prepare the data for insertion or patch into the database
func (user UserCreate) ToDBGorm() UserDBGorm {
	return UserDBGorm{
		IDNumber:         user.IDNumber,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		PersonalEmail:    user.PersonalEmail,
		InstitutionEmail: user.InstitutionEmail,
		ResidenceAddress: user.ResidenceAddress,
		Semester:         user.Semester,
		UniversityID:     user.UniversityID,
		PhoneNumber:      user.PhoneNumber,
	}
}

// ToPutDB converts a UserCreate struct to a map[string]any
// This is used to prepare the data for updating a user in the database
// It filters out fields that are not needed for the update or should not be zeroed
func (user UserCreate) ToPutDB() map[string]any {
	filter := func(key reflect.StructField, value reflect.Value) bool {
		excludeFields := []string{"id", "created_at", "updated_at", "deleted_at"}
		if slices.Contains(excludeFields, key.Tag.Get("json")) {
			return false
		}

		tags := strings.Split(key.Tag.Get("gorm"), ";")
		if slices.Contains(tags, "not null") && value.IsZero() {
			return false
		}

		return true
	}

	return utils.StructToMap(user, filter)
}

// ToDB converts a UserDB struct to a UserResponse struct
// This is used to prepare the data for returning to the client
func (user UserDBGorm) ToResponse() UserResponse {
	return UserResponse{
		ID:               user.ID,
		IDNumber:         user.IDNumber,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		PersonalEmail:    user.PersonalEmail,
		InstitutionEmail: user.InstitutionEmail,
		ResidenceAddress: user.ResidenceAddress,
		Semester:         user.Semester,
		UniversityID:     user.UniversityID,
		PhoneNumber:      user.PhoneNumber,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}
}

// TableName returns the name of the table in the database for the UserDB struct
// This is used by GORM to determine the table name for the model
func (UserDBGorm) TableName() string {
	return "users"
}
