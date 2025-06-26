package models

import "dainxor/atv/logger"

type SessionDBMongo struct {
	ID                    DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	IDStudent             DBID       `json:"id_student,omitempty" bson:"id_student,omitempty"`
	StudentName           string     `json:"name,omitempty" bson:"name,omitempty"`
	StudentSurname        string     `json:"surname,omitempty" bson:"surname,omitempty"`
	IDCompanion           DBID       `json:"id_companion,omitempty" bson:"id_companion,omitempty"`
	CompanionName         string     `json:"companion_name,omitempty" bson:"companion_name,omitempty"`
	CompanionSurname      string     `json:"companion_surname,omitempty" bson:"companion_surname,omitempty"`
	IDCompanionSpeciality DBID       `json:"id_companion_speciality,omitempty" bson:"id_companion_speciality,omitempty"`
	SessionNotes          string     `json:"session_notes,omitempty" bson:"session_notes,omitempty"`
	Date                  string     `json:"date,omitempty" bson:"date,omitempty"`
	CreatedAt             DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt             DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt             DBDateTime `json:"deleted_at,omitzero" bson:"deleted_at,omitempty"`
}

type SessionDBMongoReceiver struct {
	ID                    any        `json:"_id,omitempty" bson:"_id,omitempty"`
	IDStudent             any        `json:"id_student,omitempty" bson:"id_student,omitempty"`
	StudentName           string     `json:"name,omitempty" bson:"name,omitempty"`
	StudentSurname        string     `json:"surname,omitempty" bson:"surname,omitempty"`
	IDCompanion           any        `json:"id_companion,omitempty" bson:"id_companion,omitempty"`
	CompanionName         string     `json:"companion_name,omitempty" bson:"companion_name,omitempty"`
	CompanionSurname      string     `json:"companion_surname,omitempty" bson:"companion_surname,omitempty"`
	IDCompanionSpeciality any        `json:"id_companion_speciality,omitempty" bson:"id_companion_speciality,omitempty"`
	SessionNotes          string     `json:"session_notes,omitempty" bson:"session_notes,omitempty"`
	Date                  string     `json:"date,omitempty" bson:"date,omitempty"`
	CreatedAt             DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt             DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
	DeletedAt             DBDateTime `json:"deleted_at,omitzero" bson:"deleted_at,omitzero"`
}

func (SessionDBMongo) Receiver() SessionDBMongoReceiver {
	return SessionDBMongoReceiver{}
}
func (SessionDBMongo) ReceiverList() []SessionDBMongoReceiver {
	u := make([]SessionDBMongoReceiver, 1)
	u[0] = SessionDBMongo{}.Receiver()
	return u
}

// SessionCreate represents the request body for creating a new session
type SessionCreate struct {
	IDStudent    string `json:"id_student,omitempty" bson:"id_student,omitempty"`
	IDCompanion  string `json:"id_companion,omitempty" bson:"id_companion,omitempty"`
	SessionNotes string `json:"session_notes,omitempty" bson:"session_notes,omitempty"`
	Date         string `json:"date,omitempty" bson:"date,omitempty"`
}

// SessionResponse represents the response body for a session
type SessionResponse struct {
	ID                    string     `json:"_id,omitempty" bson:"_id,omitempty"`
	IDStudent             string     `json:"id_student,omitempty" bson:"id_student,omitempty"`
	StudentName           string     `json:"name,omitempty" bson:"name,omitempty"`
	StudentSurname        string     `json:"surname,omitempty" bson:"surname,omitempty"`
	IDCompanion           string     `json:"id_companion,omitempty" bson:"id_companion,omitempty"`
	CompanionName         string     `json:"companion_name,omitempty" bson:"companion_name,omitempty"`
	CompanionSurname      string     `json:"companion_surname,omitempty" bson:"companion_surname,omitempty"`
	IDCompanionSpeciality string     `json:"id_companion_speciality,omitempty" bson:"id_companion_speciality,omitempty"`
	SessionNotes          string     `json:"session_notes,omitempty" bson:"session_notes,omitempty"`
	Date                  string     `json:"date,omitempty" bson:"date,omitempty"`
	CreatedAt             DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt             DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
}

func (u SessionCreate) ToInsert(extra any) SessionDBMongo {
	idStudent, _ := DBIDFrom(u.IDStudent)
	idCompanion, _ := DBIDFrom(u.IDCompanion)
	idCompanionSpeciality, _ := DBIDFrom((extra.(map[string]string))["IDCompanionSpeciality"])
	extraInfo, ok := extra.(map[string]string)

	if !ok {
		logger.Error("Expected extra to be a map[string]string, got: ", extra)
		return SessionDBMongo{}
	}

	return SessionDBMongo{
		IDStudent:             idStudent,
		StudentName:           extraInfo["StudentName"],
		StudentSurname:        extraInfo["StudentSurname"],
		IDCompanion:           idCompanion,
		CompanionName:         extraInfo["CompanionName"],
		CompanionSurname:      extraInfo["CompanionSurname"],
		IDCompanionSpeciality: idCompanionSpeciality,
		SessionNotes:          u.SessionNotes,
		Date:                  u.Date,
		CreatedAt:             TimeNow(),
		UpdatedAt:             TimeNow(),
	}
}
func (u SessionCreate) ToUpdate(extra map[string]string) SessionDBMongo {
	idStudent, _ := DBIDFrom(u.IDStudent)
	idCompanion, _ := DBIDFrom(u.IDCompanion)
	idCompanionSpeciality, _ := DBIDFrom(extra["IDCompanionSpeciality"])

	return SessionDBMongo{
		IDStudent:             idStudent,
		StudentName:           extra["StudentName"],
		StudentSurname:        extra["StudentSurname"],
		IDCompanion:           idCompanion,
		CompanionName:         extra["CompanionName"],
		CompanionSurname:      extra["CompanionSurname"],
		IDCompanionSpeciality: idCompanionSpeciality,
		SessionNotes:          u.SessionNotes,
		Date:                  u.Date,
		UpdatedAt:             TimeNow(),
	}
}

func (u SessionDBMongoReceiver) ToDB() SessionDBMongo {
	id, _ := DBIDFrom(u.ID)
	idStudent, _ := DBIDFrom(u.IDStudent)
	idCompanion, _ := DBIDFrom(u.IDCompanion)
	idCompanionSpeciality, _ := DBIDFrom(u.IDCompanionSpeciality)

	return SessionDBMongo{
		ID:                    id,
		IDStudent:             idStudent,
		StudentName:           u.StudentName,
		StudentSurname:        u.CompanionName,
		IDCompanion:           idCompanion,
		CompanionName:         u.CompanionName,
		CompanionSurname:      u.CompanionSurname,
		IDCompanionSpeciality: idCompanionSpeciality,
		SessionNotes:          u.SessionNotes,
		Date:                  u.Date,
		CreatedAt:             u.CreatedAt,
		UpdatedAt:             u.UpdatedAt,
		DeletedAt:             u.DeletedAt,
	}
}
func (u SessionDBMongo) ToResponse() SessionResponse {
	return SessionResponse{
		ID:                    u.ID.Hex(),
		IDStudent:             u.IDStudent.Hex(),
		StudentName:           u.StudentName,
		StudentSurname:        u.CompanionName,
		IDCompanion:           u.IDCompanion.Hex(),
		CompanionName:         u.CompanionName,
		CompanionSurname:      u.CompanionSurname,
		IDCompanionSpeciality: u.IDCompanionSpeciality.Hex(),
		SessionNotes:          u.SessionNotes,
		Date:                  u.Date,
		CreatedAt:             u.CreatedAt,
		UpdatedAt:             u.UpdatedAt,
	}
}

func (SessionDBMongo) TableName() string {
	return "sessions"
}
func (SessionDBMongoReceiver) TableName() string {
	return SessionDBMongo{}.TableName()
}

var _ DBModelInterface = (*SessionDBMongo)(nil)
var _ DBModelInterface = (*SessionDBMongoReceiver)(nil)
