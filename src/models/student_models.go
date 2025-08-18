package models

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"errors"
	"time"
)

// UserDBGorm represents the database model for a user

type StudentDB struct {
	ID               DBID   `json:"_id,omitempty" bson:"_id,omitempty"`
	NumberID         string `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string `json:"last_name,omitempty" bson:"last_name,omitempty"`
	PersonalEmail    string `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	ResidenceAddress string `json:"residence_address,omitempty" bson:"residence_address,omitempty"`
	Semester         uint   `json:"semester,omitempty" bson:"semester,omitempty"`
	IDUniversity     DBID   `json:"id_university,omitempty" bson:"id_university,omitempty"`
	PhoneNumber      string `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	DBModelBase
}

// StudentCreate represents the request body for creating a new user or updating an existing user
// It is used to validate the input data before creating or updating a user in the database
type StudentCreate struct {
	NumberID         string `json:"number_id" gorm:"unique;not null"`
	FirstName        string `json:"first_name" gorm:"not null"`
	LastName         string `json:"last_name" gorm:"not null"`
	PersonalEmail    string `json:"email" gorm:"unique;not null"`
	InstitutionEmail string `json:"institution_email" gorm:"unique;not null"`
	ResidenceAddress string `json:"residence_address" gorm:"not null"`
	Semester         uint   `json:"semester" gorm:"not null"`
	IDUniversity     string `json:"id_university" gorm:"not null"`
	PhoneNumber      string `json:"phone_number"`
}

// StudentResponse represents the response body for a user
// It is used to format the data returned to the client after a user is created or retrieved
// It includes the ID, created_at, and updated_at fields
type StudentResponse struct {
	ID               string    `json:"id" gorm:"primaryKey;autoIncrement"`
	NumberID         string    `json:"number_id" gorm:"unique;not null"`
	FirstName        string    `json:"first_name" gorm:"not null"`
	LastName         string    `json:"last_name" gorm:"not null"`
	PersonalEmail    string    `json:"email" gorm:"unique;not null"`
	InstitutionEmail string    `json:"institution_email" gorm:"unique;not null"`
	ResidenceAddress string    `json:"residence_address" gorm:"not null"`
	Semester         uint      `json:"semester" gorm:"not null"`
	IDUniversity     string    `json:"id_university" gorm:"not null"`
	PhoneNumber      string    `json:"phone_number"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ToInsert and ToUpdate converts a UserCreate struct to a UserDBMongo struct
// This is used to prepare the data for insertion into the MongoDB database
func (user StudentCreate) ToInsert() StudentDB {
	idu, err := ID.ToDB(user.IDUniversity)

	if err != nil {
		logger.Warning("Failed to convert IDUniversity to primitive.ObjectID:", err)
		return StudentDB{} // Return an empty struct if conversion fails
	}

	return StudentDB{
		NumberID:         user.NumberID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		PersonalEmail:    user.PersonalEmail,
		InstitutionEmail: user.InstitutionEmail,
		ResidenceAddress: user.ResidenceAddress,
		Semester:         user.Semester,
		IDUniversity:     idu,
		PhoneNumber:      user.PhoneNumber,
		DBModelBase: DBModelBase{
			CreatedAt: Time.Now(),
			UpdatedAt: Time.Now(),
			DeletedAt: Time.Zero(),
		},
	}
}

// This is used to prepare the data for patch into the MongoDB database
func (user StudentCreate) ToUpdate() types.Result[StudentDB] {
	obj := StudentDB{
		NumberID:         user.NumberID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		PersonalEmail:    user.PersonalEmail,
		InstitutionEmail: user.InstitutionEmail,
		ResidenceAddress: user.ResidenceAddress,
		Semester:         user.Semester,
		PhoneNumber:      user.PhoneNumber,
		DBModelBase: DBModelBase{
			UpdatedAt: Time.Now(),
		},
	}

	if !ID.OmitEmpty(user.IDUniversity, &obj.IDUniversity, "IDUniversity") {
		return types.ResultErr[StudentDB](errors.New("IDUniversity is required"))
	}

	return types.ResultOk(obj)
}

// ToDB converts a UserDB struct to a UserResponse struct
// This is used to prepare the data for returning to the client
func (user StudentDB) ToResponse() StudentResponse {
	return StudentResponse{
		ID:               user.ID.Hex(),
		NumberID:         user.NumberID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		PersonalEmail:    user.PersonalEmail,
		InstitutionEmail: user.InstitutionEmail,
		ResidenceAddress: user.ResidenceAddress,
		Semester:         user.Semester,
		IDUniversity:     user.IDUniversity.Hex(),
		PhoneNumber:      user.PhoneNumber,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}
}
func (user StudentDB) IsEmpty() bool {
	return user == (StudentDB{})
}

// TableName returns the name of the table in the database for the UserDB struct
// This is used by GORM to determine the table name for the model
func (StudentDB) TableName() string {
	return "students"
}

// Explicitly checking if the structs implement the DBModelInterface
// This will error in compile time if the structs do not implement the interface
var _ DBModelInterface = (*StudentDB)(nil)
