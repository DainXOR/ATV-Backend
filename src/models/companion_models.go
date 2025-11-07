package models

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"errors"
)

type CompanionDB struct {
	ID               DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	NumberID         string     `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email            string     `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string     `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	PhoneNumber      string     `json:"phone_number" bson:"phone_number"`
	IDSpeciality     DBID       `json:"id_speciality" bson:"id_speciality"`
	CreatedAt        DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt        DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt        DBDateTime `json:"deleted_at" bson:"deleted_at"`
}
type CompanionCreate struct {
	NumberID         string `json:"number_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	InstitutionEmail string `json:"institution_email"`
	PhoneNumber      string `json:"phone_number"`
	IDSpeciality     string `json:"id_speciality"`
}
type CompanionResponse struct {
	ID               string     `json:"id"`
	NumberID         string     `json:"number_id"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Email            string     `json:"email"`
	InstitutionEmail string     `json:"institution_email"`
	PhoneNumber      string     `json:"phone_number"`
	IDSpeciality     string     `json:"id_speciality"`
	CreatedAt        DBDateTime `json:"created_at"`
	UpdatedAt        DBDateTime `json:"updated_at"`
}

func (c CompanionCreate) ToInsert() types.Result[CompanionDB] {
	obj := CompanionDB{
		NumberID:         c.NumberID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		InstitutionEmail: c.InstitutionEmail,
		PhoneNumber:      c.PhoneNumber,
		CreatedAt:        Time.Now(),
		UpdatedAt:        Time.Now(),
		DeletedAt:        Time.Zero(),
	}

	if !ID.Ensure(c.IDSpeciality, &obj.IDSpeciality, "IDSpeciality") {
		logger.Lava(types.V("0.2.0"), "Using not standarized error")
		return types.ResultErr[CompanionDB](errors.New("Invalid IDSpeciality"))
	}

	return types.ResultOk(obj)
}
func (c CompanionCreate) ToUpdate() types.Result[CompanionDB] {
	obj := CompanionDB{
		NumberID:         c.NumberID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		InstitutionEmail: c.InstitutionEmail,
		PhoneNumber:      c.PhoneNumber,
		UpdatedAt:        Time.Now(),
	}

	if !ID.OmitEmpty(c.IDSpeciality, &obj.IDSpeciality, "IDSpeciality") {
		return types.ResultErr[CompanionDB](errors.New("Invalid IDSpeciality"))
	}

	return types.ResultOk(obj)
}
func (c CompanionDB) ToResponse() CompanionResponse {
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

func (c CompanionDB) IsEmpty() bool {
	return c == (CompanionDB{})
}

func (CompanionDB) TableName() string {
	return "companions"
}

var _ DBModelInterface = (*CompanionDB)(nil)
