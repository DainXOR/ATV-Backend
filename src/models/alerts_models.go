package models

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"errors"
)

type AlertDB struct {
	ID              DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	IDPriority      DBID       `json:"id_priority,omitempty" bson:"id_priority,omitempty"`
	IDStudent       DBID       `json:"id_student,omitempty" bson:"id_student,omitempty"`
	IDVulnerability DBID       `json:"id_vulnerability,omitempty" bson:"id_vulnerability,omitempty"`
	Message         string     `json:"message,omitempty" bson:"message,omitempty"`
	CreatedAt       DBDateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       DBDateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt       DBDateTime `json:"deleted_at" bson:"deleted_at"`
}

type AlertCreate struct {
	IDPriority      string `json:"id_priority"`
	IDStudent       string `json:"id_student"`
	IDVulnerability string `json:"id_vulnerability"`
	Message         string `json:"message"`
}

type AlertResponse struct {
	ID              string     `json:"id"`
	IDPriority      string     `json:"id_priority"`
	IDStudent       string     `json:"id_student"`
	IDVulnerability string     `json:"id_vulnerability"`
	Message         string     `json:"message"`
	CreatedAt       DBDateTime `json:"created_at"`
	UpdatedAt       DBDateTime `json:"updated_at"`
}

func (a AlertCreate) ToInsert() types.Result[AlertDB] {
	obj := AlertDB{
		Message:   a.Message,
		CreatedAt: Time.Now(),
		UpdatedAt: Time.Now(),
		DeletedAt: Time.Zero(),
	}

	if !ID.Ensure(a.IDPriority, &obj.IDPriority, "IDPriority") ||
		!ID.Ensure(a.IDStudent, &obj.IDStudent, "IDStudent") ||
		!ID.Ensure(a.IDVulnerability, &obj.IDVulnerability, "IDVulnerability") {

		logger.Lava(types.V("0.2.0"), "Using not standarized error")
		return types.ResultErr[AlertDB](errors.New(""))
	}

	return types.ResultOk(obj)
}
func (a AlertCreate) ToUpdate() types.Result[AlertDB] {
	obj := AlertDB{
		Message:   a.Message,
		UpdatedAt: Time.Now(),
	}

	if !ID.OmitEmpty(a.IDPriority, &obj.IDPriority, "IDPriority") ||
		!ID.OmitEmpty(a.IDStudent, &obj.IDStudent, "IDStudent") ||
		!ID.OmitEmpty(a.IDVulnerability, &obj.IDVulnerability, "IDVulnerability") {

		logger.Lava(types.V("0.2.0"), "Using not standarized error")
		return types.ResultErr[AlertDB](errors.New(""))
	}

	return types.ResultOk(obj)
}

func (a AlertDB) ToResponse() AlertResponse {
	return AlertResponse{
		ID:              a.ID.Hex(),
		IDPriority:      a.IDPriority.Hex(),
		IDStudent:       a.IDStudent.Hex(),
		IDVulnerability: a.IDVulnerability.Hex(),
		Message:         a.Message,
		CreatedAt:       Time.Now(),
		UpdatedAt:       Time.Now(),
	}
}

func (a AlertDB) IsEmpty() bool {
	return a == (AlertDB{})
}

func (AlertDB) TableName() string {
	return "alerts"
}

var _ DBModelInterface[AlertCreate, AlertResponse] = (*AlertDB)(nil)
var _ CreateModelInterface[AlertResponse, AlertDB] = (*AlertCreate)(nil)
