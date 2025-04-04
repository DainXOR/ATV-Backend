package models

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	ID               int    `json:"id" gorm:"primaryKey;autoIncrement"`
	IDNumber         int    `json:"id_number" gorm:"unique;not null"`
	FullName         string `json:"username" gorm:"not null"`
	PersonalEmail    string `json:"email" gorm:"unique;not null"`
	InstitutionEmail string `json:"institution_email" gorm:"unique;not null"`
	ResidenceAddress string `json:"residence_address" gorm:"not null"`
	Semester         int    `json:"semester" gorm:"not null"`
	UniversityID     int    `json:"university_id" gorm:"not null"`
	PhoneNumber      string `json:"phone_number"`
}

type UserCreate struct {
	ID               int    `json:"id" gorm:"primaryKey;autoIncrement"`
	IDNumber         int    `json:"id_number" gorm:"unique;not null"`
	FullName         string `json:"username" gorm:"not null"`
	PersonalEmail    string `json:"email" gorm:"unique;not null"`
	InstitutionEmail string `json:"institution_email" gorm:"unique;not null"`
	ResidenceAddress string `json:"residence_address" gorm:"not null"`
	Semester         int    `json:"semester" gorm:"not null"`
	UniversityID     int    `json:"university_id" gorm:"not null"`
}

type UserResponse struct {
	ID               int       `json:"id" gorm:"primaryKey;autoIncrement"`
	IDNumber         int       `json:"id_number" gorm:"unique;not null"`
	FullName         string    `json:"username" gorm:"not null"`
	PersonalEmail    string    `json:"email" gorm:"unique;not null"`
	InstitutionEmail string    `json:"institution_email" gorm:"unique;not null"`
	ResidenceAddress string    `json:"residence_address" gorm:"not null"`
	Semester         int       `json:"semester" gorm:"not null"`
	UniversityID     int       `json:"university_id" gorm:"not null"`
	PhoneNumber      string    `json:"phone_number"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type UserDB struct {
	gorm.Model
	IDNumber         int    `json:"id_number" gorm:"unique;not null"`
	FullName         string `json:"username" gorm:"not null"`
	PersonalEmail    string `json:"email" gorm:"unique;not null"`
	InstitutionEmail string `json:"institution_email" gorm:"unique;not null"`
	ResidenceAddress string `json:"residence_address" gorm:"not null"`
	Semester         int    `json:"semester" gorm:"not null"`
	UniversityID     int    `json:"university_id" gorm:"not null"`
	PhoneNumber      string `json:"phone_number"`
}
