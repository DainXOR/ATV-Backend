package models

import (
	"dainxor/atv/logger"
	"time"
)

// UserDBGorm represents the database model for a user

type StudentDBMongo struct {
	ID               DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	NumberID         string     `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
	PersonalEmail    string     `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string     `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	ResidenceAddress string     `json:"residence_address,omitempty" bson:"residence_address,omitempty"`
	Semester         uint       `json:"semester,omitempty" bson:"semester,omitempty"`
	IDUniversity     DBID       `json:"id_university,omitempty" bson:"id_university,omitempty"`
	PhoneNumber      string     `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	CreatedAt        DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt        DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt        DBDateTime `json:"deleted_at" bson:"deleted_at"`
}

type StudentDBMongoReceiver struct {
	ID               any        `json:"_id,omitempty" bson:"_id,omitempty"`
	NumberID         string     `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
	PersonalEmail    string     `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string     `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	ResidenceAddress string     `json:"residence_address,omitempty" bson:"residence_address,omitempty"`
	Semester         uint       `json:"semester,omitempty" bson:"semester,omitempty"`
	IDUniversity     any        `json:"id_university,omitempty" bson:"id_university,omitempty"`
	PhoneNumber      string     `json:"phone_number" bson:"phone_number"`
	CreatedAt        DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt        DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
	DeletedAt        DBDateTime `json:"deleted_at,omitzero" bson:"deleted_at,omitzero"`
}

func (StudentDBMongo) Receiver() StudentDBMongoReceiver {
	return StudentDBMongoReceiver{}
}
func (StudentDBMongo) ReceiverList() []StudentDBMongoReceiver {
	s := make([]StudentDBMongoReceiver, 1)
	s[0] = StudentDBMongo{}.Receiver()
	return s
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

// ToDB converts a UserCreate struct to a UserDBMongo struct
// This is used to prepare the data for insertion or patch into the MongoDB database
func (user StudentCreate) ToInsert() StudentDBMongo {
	idu, err := DBIDFrom(user.IDUniversity)

	if err != nil {
		logger.Warning("Failed to convert IDUniversity to primitive.ObjectID:", err)
		return StudentDBMongo{} // Return an empty struct if conversion fails
	}

	return StudentDBMongo{
		NumberID:         user.NumberID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		PersonalEmail:    user.PersonalEmail,
		InstitutionEmail: user.InstitutionEmail,
		ResidenceAddress: user.ResidenceAddress,
		Semester:         user.Semester,
		IDUniversity:     idu,
		PhoneNumber:      user.PhoneNumber,
		CreatedAt:        TimeNow(),
		UpdatedAt:        TimeNow(),
		DeletedAt:        TimeZero(), // DeletedAt is nil by default, indicating the user is not deleted
	}
}
func (user StudentCreate) ToUpdate() StudentDBMongo {
	idu, err := DBIDFrom(user.IDUniversity)

	if err != nil {
		logger.Warning("Failed to convert IDUniversity to primitive.ObjectID:", err)
		return StudentDBMongo{} // Return an empty struct if conversion fails
	}

	return StudentDBMongo{
		NumberID:         user.NumberID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		PersonalEmail:    user.PersonalEmail,
		InstitutionEmail: user.InstitutionEmail,
		ResidenceAddress: user.ResidenceAddress,
		Semester:         user.Semester,
		IDUniversity:     idu,
		PhoneNumber:      user.PhoneNumber,
		UpdatedAt:        TimeNow(),
	}
}
func (user StudentDBMongoReceiver) ToDB() StudentDBMongo {
	id, err1 := DBIDFrom(user.ID)
	idu, err2 := DBIDFrom(user.IDUniversity)

	if err1 != nil {
		logger.Warning("Failed to convert ID to primitive.ObjectID:", err1)
		return StudentDBMongo{} // Return an empty struct if conversion fails
	}

	if err2 != nil {
		logger.Warning("Failed to convert IDUniversity to primitive.ObjectID:", err2)
		return StudentDBMongo{} // Return an empty struct if conversion fails
	}

	return StudentDBMongo{
		ID:               id,
		NumberID:         user.NumberID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		PersonalEmail:    user.PersonalEmail,
		InstitutionEmail: user.InstitutionEmail,
		ResidenceAddress: user.ResidenceAddress,
		Semester:         user.Semester,
		IDUniversity:     idu,
		PhoneNumber:      user.PhoneNumber,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}

}

// ToDB converts a UserDB struct to a UserResponse struct
// This is used to prepare the data for returning to the client

func (user StudentDBMongo) ToResponse() StudentResponse {
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

// TableName returns the name of the table in the database for the UserDB struct
// This is used by GORM to determine the table name for the model
func (StudentDBMongo) TableName() string {
	return "students"
}
func (StudentDBMongoReceiver) TableName() string {
	return StudentDBMongo{}.TableName()
}

// Explicitly checking if the structs implement the DBModelInterface
// This will error in compile time if the structs do not implement the interface
var _ DBModelInterface = (*StudentDBMongo)(nil)
var _ DBModelInterface = (*StudentDBMongoReceiver)(nil)
