package models

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"errors"
)

type SessionDB struct {
	ID                  DBID          `json:"_id,omitempty" bson:"_id,omitempty"`
	IDStudent           DBID          `json:"id_student,omitempty" bson:"id_student,omitempty"`
	StudentName         string        `json:"first_name_student,omitempty" bson:"first_name_student,omitempty"`
	StudentSurname      string        `json:"last_name_student,omitempty" bson:"last_name_student,omitempty"`
	IDCompanion         DBID          `json:"id_companion,omitempty" bson:"id_companion,omitempty"`
	CompanionName       string        `json:"first_name_companion,omitempty" bson:"first_name_companion,omitempty"`
	CompanionSurname    string        `json:"last_name_companion,omitempty" bson:"last_name_companion,omitempty"`
	CompanionSpeciality string        `json:"companion_speciality,omitempty" bson:"companion_speciality,omitempty"`
	IDSessionType       DBID          `json:"id_session_type,omitempty" bson:"id_session_type,omitempty"`
	SessionNotes        string        `json:"session_notes,omitempty" bson:"session_notes,omitempty"`
	DeprDate            string        `json:"date,omitempty" bson:"temp_date,omitempty"`
	Date                DBDateTime    `json:"session_date,omitempty" bson:"date,omitempty"`
	Status              sessionStatus `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt           DBDateTime    `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt           DBDateTime    `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt           DBDateTime    `json:"deleted_at" bson:"deleted_at"`
}

// SessionCreate represents the request body for creating a new session or updating an existing one
type SessionCreate struct {
	IDStudent     string `json:"id_student,omitempty"`
	IDCompanion   string `json:"id_companion,omitempty"`
	IDSessionType string `json:"id_session_type,omitempty"`
	SessionNotes  string `json:"session_notes,omitempty"`
	Status        string `json:"status,omitempty"`
	Date          string `json:"date,omitempty"`
}

// SessionResponse represents the response body for a session
type SessionResponse struct {
	ID                  string     `json:"id,omitempty"`
	IDStudent           string     `json:"id_student,omitempty"`
	StudentName         string     `json:"name,omitempty"`
	StudentSurname      string     `json:"surname,omitempty"`
	IDCompanion         string     `json:"id_companion,omitempty"`
	CompanionName       string     `json:"companion_name,omitempty"`
	CompanionSurname    string     `json:"companion_surname,omitempty"`
	CompanionSpeciality string     `json:"companion_speciality,omitempty"`
	IDSessionType       string     `json:"id_session_type,omitempty"`
	SessionNotes        string     `json:"session_notes,omitempty"`
	DeprDate            string     `json:"date,omitempty"`
	Date                DBDateTime `json:"session_date,omitzero"`
	Status              string     `json:"status,omitempty"`
	CreatedAt           DBDateTime `json:"created_at,omitzero"`
	UpdatedAt           DBDateTime `json:"updated_at,omitzero"`
}

type sessionStatus uint8

// To add new session statuses, simply follow this steps:
//
// 1. Define the new status as a constant in the sessionStatus type.
//
// 2. Add the new status to the STATUS map with its corresponding name.
//
// The name you define here will be used in the application to refer to this status.
const (
	STATUS_UNKNOWN sessionStatus = iota + 1
	STATUS_PENDING
	STATUS_COMPLETED
	STATUS_CANCELLED
	STATUS_UNATTENDED
)

var STATUS = map[sessionStatus]string{
	STATUS_UNKNOWN:    "Desconocido",
	STATUS_PENDING:    "Pendiente",
	STATUS_COMPLETED:  "Completado",
	STATUS_CANCELLED:  "Cancelado",
	STATUS_UNATTENDED: "No asistió",
}

func statusName(code sessionStatus) string {
	if name, exists := STATUS[code]; exists {
		return name
	}
	return "Desconocido"
}
func statusCode(name string) sessionStatus {
	for state, stateName := range STATUS {
		if stateName == name {
			return state
		}
	}
	return STATUS_UNKNOWN
}

func (u SessionCreate) ToInsert(extra map[string]string) types.Optional[SessionDB] {
	obj := SessionDB{
		StudentName:         extra["StudentName"],
		StudentSurname:      extra["StudentSurname"],
		CompanionName:       extra["CompanionName"],
		CompanionSurname:    extra["CompanionSurname"],
		CompanionSpeciality: extra["CompanionSpeciality"],
		SessionNotes:        u.SessionNotes,
		DeprDate:            u.Date,
		Status:              statusCode(u.Status),
		CreatedAt:           Time.Now(),
		UpdatedAt:           Time.Now(),
		DeletedAt:           Time.Zero(),
	}

	if !ID.Ensure(u.IDStudent, &obj.IDStudent, "IDStudent") ||
		!ID.Ensure(u.IDCompanion, &obj.IDCompanion, "IDCompanion") ||
		!ID.Ensure(u.IDSessionType, &obj.IDSessionType, "IDSessionType") {
		return types.OptionalEmpty[SessionDB]()
	}
	if date, err := Time.Parse(u.Date, Time.Format()); err == nil {
		obj.Date = date
	} else {
		logger.Warning("Failed to parse session date:", err)
	}

	return types.OptionalOf(obj)
}
func (u SessionCreate) ToUpdate(extra map[string]string) types.Result[SessionDB] {
	obj := SessionDB{
		StudentName:         extra["StudentName"],
		StudentSurname:      extra["StudentSurname"],
		CompanionName:       extra["CompanionName"],
		CompanionSurname:    extra["CompanionSurname"],
		CompanionSpeciality: extra["CompanionSpeciality"],
		SessionNotes:        u.SessionNotes,
		DeprDate:            u.Date,
		Status:              statusCode(u.Status),
		UpdatedAt:           Time.Now(),
	}

	if !ID.OmitEmpty(u.IDStudent, &obj.IDStudent, "IDStudent") ||
		!ID.OmitEmpty(u.IDCompanion, &obj.IDCompanion, "IDCompanion") ||
		!ID.OmitEmpty(u.IDSessionType, &obj.IDSessionType, "IDSessionType") {
		return types.ResultErr[SessionDB](errors.New("Invalid session data"))
	}
	if date, err := Time.Parse(u.Date, Time.Format()); err == nil {
		obj.Date = date
	} else {
		logger.Warning("Failed to parse session date:", err)
	}

	return types.ResultOk(obj)
}
func (u SessionDB) ToResponse() SessionResponse {
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
		DeprDate:            u.DeprDate,
		Date:                u.Date,
		Status:              statusName(u.Status),
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
	}
}
func (u SessionDB) IsEmpty() bool {
	return u == (SessionDB{})
}

func (SessionDB) TableName() string {
	return "sessions"
}

var _ DBModelInterface = (*SessionDB)(nil)
