package models

import (
	"time"

	"gorm.io/gorm"
)

type UserDB struct {
	gorm.Model
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

func (user UserCreate) ToDB() UserDB {
	return UserDB{
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

func (user UserDB) ToResponse() UserResponse {
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

func (UserDB) TableName() string {
	return "users"
}
