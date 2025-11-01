package models

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"errors"
)

type FormQuestionDB struct {
	ID             DBID       `json:"id" bson:"_id,omitempty"`
	Name           string     `json:"name" bson:"name,omitempty"`
	Question       string     `json:"question" bson:"question,omitempty"`
	Options        []string   `json:"options" bson:"options,omitempty"`
	IDQuestionType DBID       `json:"id_question_type" bson:"id_question_type,omitempty"`
	CreatedAt      DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt      DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt      DBDateTime `json:"deleted_at" bson:"deleted_at"`
}
type FormQuestionCreate struct {
	Name           string   `json:"name"`
	Question       string   `json:"question"`
	Options        []string `json:"options"`
	IDQuestionType string   `json:"id_question_type"`
}
type FormQuestionResponse struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Question       string     `json:"question"`
	Options        []string   `json:"options"`
	IDQuestionType string     `json:"id_question_type"`
	CreatedAt      DBDateTime `json:"created_at"`
	UpdatedAt      DBDateTime `json:"updated_at"`
}

func (q FormQuestionCreate) ToInsert() types.Result[FormQuestionDB] {
	obj := FormQuestionDB{
		Name:      q.Name,
		Question:  q.Question,
		Options:   q.Options,
		CreatedAt: Time.Now(),
		UpdatedAt: Time.Now(),
		DeletedAt: Time.Zero(),
	}

	if !ID.Ensure(q.IDQuestionType, &obj.IDQuestionType, "IDQuestionType") {
		logger.Lava(types.V("0.2.1"), "Using not standarized error")
		return types.ResultErr[FormQuestionDB](errors.New("Invalid IDQuestionType"))
	}

	return types.ResultOk(obj)
}
func (q FormQuestionDB) ToUpdate() FormQuestionDB {
	return FormQuestionDB{
		Name:           q.Name,
		Question:       q.Question,
		Options:        q.Options,
		IDQuestionType: q.IDQuestionType,
		UpdatedAt:      Time.Now(),
	}
}
func (q FormQuestionDB) ToResponse() FormQuestionResponse {
	return FormQuestionResponse{
		ID:             q.ID.Hex(),
		Name:           q.Name,
		Question:       q.Question,
		Options:        q.Options,
		IDQuestionType: q.IDQuestionType.Hex(),
		UpdatedAt:      q.UpdatedAt,
		CreatedAt:      q.CreatedAt,
	}
}

func (q FormQuestionDB) IsEmpty() bool {
	zeroObj := FormQuestionDB{}

	cmp := q.ID == zeroObj.ID &&
		q.Name == zeroObj.Name &&
		q.Question == zeroObj.Question &&
		q.IDQuestionType == zeroObj.IDQuestionType &&
		q.CreatedAt.Equal(zeroObj.CreatedAt) &&
		q.UpdatedAt.Equal(zeroObj.UpdatedAt) &&
		q.DeletedAt.Equal(zeroObj.DeletedAt)

	if cmp {
		for i, o := range q.Options {
			cmp = cmp && (o == zeroObj.Options[i])
		}
	}

	return cmp
}

func (FormQuestionDB) TableName() string {
	return "form_questions"
}

var _ DBModelInterface = (*FormQuestionDB)(nil)
