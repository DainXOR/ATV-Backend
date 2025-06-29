package models

import (
	"dainxor/atv/types"
)

type SessionDBMongo struct {
	ID                  DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	IDStudent           DBID       `json:"id_student,omitempty" bson:"id_student,omitempty"`
	StudentName         string     `json:"first_name_student,omitempty" bson:"first_name_student,omitempty"`
	StudentSurname      string     `json:"last_name_student,omitempty" bson:"last_name_student,omitempty"`
	IDCompanion         DBID       `json:"id_companion,omitempty" bson:"id_companion,omitempty"`
	CompanionName       string     `json:"first_name_companion,omitempty" bson:"first_name_companion,omitempty"`
	CompanionSurname    string     `json:"last_name_companion,omitempty" bson:"last_name_companion,omitempty"`
	CompanionSpeciality string     `json:"companion_speciality,omitempty" bson:"companion_speciality,omitempty"`
	IDSessionType       DBID       `json:"id_session_type,omitempty" bson:"id_session_type,omitempty"`
	SessionNotes        string     `json:"session_notes,omitempty" bson:"session_notes,omitempty"`
	Date                string     `json:"date,omitempty" bson:"date,omitempty"`
	CreatedAt           DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt           DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt           DBDateTime `json:"deleted_at,omitzero" bson:"deleted_at,omitempty"`
}

// SessionCreate represents the request body for creating a new session
type SessionCreate struct {
	IDStudent     string `json:"id_student,omitempty" bson:"id_student,omitempty"`
	IDCompanion   string `json:"id_companion,omitempty" bson:"id_companion,omitempty"`
	IDSessionType string `json:"id_session_type,omitempty" bson:"id_session_type,omitempty"`
	SessionNotes  string `json:"session_notes,omitempty" bson:"session_notes,omitempty"`
	Date          string `json:"date,omitempty" bson:"date,omitempty"`
}

// SessionResponse represents the response body for a session
type SessionResponse struct {
	ID                  string     `json:"id,omitempty" bson:"id,omitempty"`
	IDStudent           string     `json:"id_student,omitempty" bson:"id_student,omitempty"`
	StudentName         string     `json:"name,omitempty" bson:"name,omitempty"`
	StudentSurname      string     `json:"surname,omitempty" bson:"surname,omitempty"`
	IDCompanion         string     `json:"id_companion,omitempty" bson:"id_companion,omitempty"`
	CompanionName       string     `json:"companion_name,omitempty" bson:"companion_name,omitempty"`
	CompanionSurname    string     `json:"companion_surname,omitempty" bson:"companion_surname,omitempty"`
	CompanionSpeciality string     `json:"companion_speciality,omitempty" bson:"companion_speciality,omitempty"`
	IDSessionType       string     `json:"id_session_type,omitempty" bson:"id_session_type,omitempty"`
	SessionNotes        string     `json:"session_notes,omitempty" bson:"session_notes,omitempty"`
	Date                string     `json:"date,omitempty" bson:"date,omitempty"`
	CreatedAt           DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt           DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
}

func (u SessionCreate) ToInsert(extra map[string]string) types.Optional[SessionDBMongo] {
	obj := SessionDBMongo{
		StudentName:         extra["StudentName"],
		StudentSurname:      extra["StudentSurname"],
		CompanionName:       extra["CompanionName"],
		CompanionSurname:    extra["CompanionSurname"],
		CompanionSpeciality: extra["CompanionSpeciality"],
		SessionNotes:        u.SessionNotes,
		Date:                u.Date,
		CreatedAt:           TimeNow(),
		UpdatedAt:           TimeNow(),
		DeletedAt:           TimeZero(),
	}

	if !EnsureID(u.IDStudent, &obj.IDStudent, "IDStudent") ||
		!EnsureID(u.IDCompanion, &obj.IDCompanion, "IDCompanion") ||
		!EnsureID(u.IDSessionType, &obj.IDSessionType, "IDSessionType") {
		return types.OptionalEmpty[SessionDBMongo]()
	}

	return types.OptionalOf(obj)
}
func (u SessionCreate) ToUpdate(extra map[string]string) types.Optional[SessionDBMongo] {
	obj := SessionDBMongo{
		StudentName:         extra["StudentName"],
		StudentSurname:      extra["StudentSurname"],
		CompanionName:       extra["CompanionName"],
		CompanionSurname:    extra["CompanionSurname"],
		CompanionSpeciality: extra["CompanionSpeciality"],
		SessionNotes:        u.SessionNotes,
		Date:                u.Date,
		UpdatedAt:           TimeNow(),
	}

	if !OmitEmptyID(u.IDStudent, &obj.IDStudent, "IDStudent") ||
		!OmitEmptyID(u.IDCompanion, &obj.IDCompanion, "IDCompanion") ||
		!OmitEmptyID(u.IDSessionType, &obj.IDSessionType, "IDSessionType") {
		return types.OptionalEmpty[SessionDBMongo]()
	}

	return types.OptionalOf(obj)
}
func (u SessionDBMongo) ToResponse() SessionResponse {
	return SessionResponse{
		ID:                  u.ID.Hex(),
		IDStudent:           u.IDStudent.Hex(),
		StudentName:         u.StudentName,
		StudentSurname:      u.CompanionName,
		IDCompanion:         u.IDCompanion.Hex(),
		CompanionName:       u.CompanionName,
		CompanionSurname:    u.CompanionSurname,
		CompanionSpeciality: u.CompanionSpeciality,
		IDSessionType:       u.IDSessionType.Hex(),
		SessionNotes:        u.SessionNotes,
		Date:                u.Date,
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
	}
}

func (SessionDBMongo) TableName() string {
	return "sessions"
}

var _ DBModelInterface = (*SessionDBMongo)(nil)
