package models

import (
	"dainxor/atv/logger"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CompanionDBMongo struct {
	ID               DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	NumberID         string     `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email            string     `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string     `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	PhoneNumber      string     `json:"phone_number" bson:"phone_number"`
	IDSpeciality     DBID       `json:"id_speciality" bson:"id_speciality"`
	CreatedAt        DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt        DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
	DeletedAt        DBDateTime `json:"deleted_at" bson:"deleted_at"`
}
type CompanionDBMongoReceiver struct {
	ID               any        `json:"_id,omitempty" bson:"_id,omitempty"`
	NumberID         string     `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email            string     `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string     `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	PhoneNumber      string     `json:"phone_number" bson:"phone_number"`
	IDSpeciality     any        `json:"id_speciality" bson:"id_speciality"`
	CreatedAt        DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt        DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
	DeletedAt        DBDateTime `json:"deleted_at" bson:"deleted_at"`
}

func (CompanionDBMongo) Receiver() CompanionDBMongoReceiver {
	return CompanionDBMongoReceiver{}
}
func (CompanionDBMongo) ReceiverList() []CompanionDBMongoReceiver {
	s := make([]CompanionDBMongoReceiver, 1)
	s[0] = CompanionDBMongo{}.Receiver()
	return s
}

type CompanionCreate struct {
	NumberID         string `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email            string `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	PhoneNumber      string `json:"phone_number" bson:"phone_number"`
	IDSpeciality     string `json:"id_speciality" bson:"id_speciality"`
}
type CompanionResponse struct {
	ID               string     `json:"_id,omitempty" bson:"_id,omitempty"`
	NumberID         string     `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email            string     `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string     `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	PhoneNumber      string     `json:"phone_number" bson:"phone_number"`
	IDSpeciality     string     `json:"id_speciality" bson:"id_speciality"`
	CreatedAt        DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt        DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
}

func (c CompanionCreate) ToInsert() CompanionDBMongo {
	idSpeciality, err := bson.ObjectIDFromHex(c.IDSpeciality)
	if err != nil {
		logger.Warning("Failed to convert IDSpeciality to PrimitiveID: ", err)
		return CompanionDBMongo{}
	}

	return CompanionDBMongo{
		NumberID:         c.NumberID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		InstitutionEmail: c.InstitutionEmail,
		PhoneNumber:      c.PhoneNumber,
		IDSpeciality:     idSpeciality,
		CreatedAt:        TimeNow(),
		UpdatedAt:        TimeNow(),
		DeletedAt:        TimeZero(),
	}
}
func (c CompanionCreate) ToUpdate() CompanionDBMongo {
	idSpeciality, err := DBIDFrom(c.IDSpeciality)
	if err != nil {
		logger.Warning("Failed to convert IDSpeciality to PrimitiveID: ", err)
		return CompanionDBMongo{}
	}

	return CompanionDBMongo{
		NumberID:         c.NumberID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		InstitutionEmail: c.InstitutionEmail,
		PhoneNumber:      c.PhoneNumber,
		IDSpeciality:     idSpeciality,
		CreatedAt:        TimeNow(),
	}
}
func (c CompanionDBMongoReceiver) ToDB() CompanionDBMongo {
	id, err2 := DBIDFrom(c.ID)
	idSpeciality, err1 := DBIDFrom(c.IDSpeciality)

	if err2 != nil {
		logger.Warning("Failed to convert ID to PrimitiveID: ", err2)
		return CompanionDBMongo{}

	} else if err1 != nil {
		logger.Warning("Failed to convert IDSpeciality to PrimitiveID: ", err1)
		return CompanionDBMongo{}
	}

	return CompanionDBMongo{
		ID:               id,
		NumberID:         c.NumberID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		InstitutionEmail: c.InstitutionEmail,
		PhoneNumber:      c.PhoneNumber,
		IDSpeciality:     idSpeciality,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

func (c CompanionDBMongo) ToResponse() CompanionResponse {
	return CompanionResponse{
		ID:               c.ID.Hex(),
		NumberID:         c.NumberID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		InstitutionEmail: c.InstitutionEmail,
		PhoneNumber:      c.PhoneNumber,
		IDSpeciality:     c.IDSpeciality.Hex(),
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

func (CompanionDBMongo) TableName() string {
	return "companions"
}
func (CompanionDBMongoReceiver) TableName() string {
	return CompanionDBMongo{}.TableName()
}

var _ DBModelInterface = (*CompanionDBMongo)(nil)
var _ DBModelInterface = (*CompanionDBMongoReceiver)(nil)
