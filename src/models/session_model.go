package models

import (
	"dainxor/atv/types"
	"errors"
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
	Status              status     `json:"status,omitempty" bson:"status,omitempty"`
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
	Status        string `json:"status,omitempty" bson:"status,omitempty"`
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
	Status              string     `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt           DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt           DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
}

type status = uint8

const (
	STATUS_UNKNOWN   status = iota // 0
	STATUS_PENDING                 // 1
	STATUS_COMPLETED               // 2
	STATUS_CANCELLED               // 3
)

var STATUS = map[status]string{
	STATUS_PENDING:   "Pendiente",
	STATUS_COMPLETED: "Completado",
	STATUS_CANCELLED: "Cancelado",
}

func statusName(code status) string {
	if name, exists := STATUS[code]; exists {
		return name
	}
	return "Desconocido"
}
func statusCode(name string) status {
	for state, stateName := range STATUS {
		if stateName == name {
			return state
		}
	}
	return 0
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
		Status:              statusCode(u.Status),
		CreatedAt:           Time.Now(),
		UpdatedAt:           Time.Now(),
		DeletedAt:           Time.Zero(),
	}

	if !ID.Ensure(u.IDStudent, &obj.IDStudent, "IDStudent") ||
		!ID.Ensure(u.IDCompanion, &obj.IDCompanion, "IDCompanion") ||
		!ID.Ensure(u.IDSessionType, &obj.IDSessionType, "IDSessionType") {
		return types.OptionalEmpty[SessionDBMongo]()
	}

	return types.OptionalOf(obj)
}
func (u SessionCreate) ToUpdate(extra map[string]string) types.Result[SessionDBMongo] {
	obj := SessionDBMongo{
		StudentName:         extra["StudentName"],
		StudentSurname:      extra["StudentSurname"],
		CompanionName:       extra["CompanionName"],
		CompanionSurname:    extra["CompanionSurname"],
		CompanionSpeciality: extra["CompanionSpeciality"],
		SessionNotes:        u.SessionNotes,
		Date:                u.Date,
		Status:              statusCode(u.Status),
		UpdatedAt:           TimeNow(),
	}

	if !ID.OmitEmpty(u.IDStudent, &obj.IDStudent, "IDStudent") ||
		!ID.OmitEmpty(u.IDCompanion, &obj.IDCompanion, "IDCompanion") ||
		!ID.OmitEmpty(u.IDSessionType, &obj.IDSessionType, "IDSessionType") {
		return types.ResultErr[SessionDBMongo](errors.New("Invalid session data"))
	}

	return types.ResultOk(obj)
}
func (u SessionDBMongo) ToResponse() SessionResponse {
	return SessionResponse{
		ID:                  u.ID.Hex(),
		IDStudent:           u.IDStudent.Hex(),
		StudentName:         u.StudentName,
		StudentSurname:      u.StudentSurname,
		IDCompanion:         u.IDCompanion.Hex(),
		CompanionName:       u.CompanionName,
		CompanionSurname:    u.CompanionSurname,
		CompanionSpeciality: u.CompanionSpeciality,
		IDSessionType:       u.IDSessionType.Hex(),
		SessionNotes:        u.SessionNotes,
		Date:                u.Date,
		Status:              statusName(u.Status),
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
	}
}
func (u SessionDBMongo) IsEmpty() bool {
	return u == (SessionDBMongo{})
}

func (SessionDBMongo) TableName() string {
	return "sessions"
}

var _ DBModelInterface = (*SessionDBMongo)(nil)
