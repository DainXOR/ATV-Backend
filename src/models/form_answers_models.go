package models

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"errors"
)

type answer[ID any] struct {
	IDQuestion ID     `json:"id_question" bson:"id_question,omitempty"`
	Answer     string `json:"answer" bson:"answer,omitempty"`
}
type FormAnswerDB struct {
	ID        DBID           `json:"id" bson:"_id,omitempty"`
	IDForm    DBID           `json:"id_form" bson:"id_form,omitempty"`
	Answers   []answer[DBID] `json:"answers" bson:"answers,omitempty"`
	CreatedAt DBDateTime     `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt DBDateTime     `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt DBDateTime     `json:"deleted_at" bson:"deleted_at"`
}
type FormAnswerCreate struct {
	IDForm  string           `json:"id_form"`
	Answers []answer[string] `json:"answers"`
}
type FormAnswerResponse struct {
	ID        DBID             `json:"id"`
	IDForm    DBID             `json:"id_form"`
	Answers   []answer[string] `json:"answers"`
	CreatedAt DBDateTime       `json:"created_at,omitzero"`
	UpdatedAt DBDateTime       `json:"updated_at,omitzero"`
}

func (o FormAnswerCreate) ToInsert() types.Result[FormAnswerDB] {
	obj := FormAnswerDB{
		CreatedAt: Time.Now(),
		UpdatedAt: Time.Now(),
		DeletedAt: Time.Zero(),
	}

	if !ID.Ensure(o.IDForm, &obj.IDForm, "IDForm") {
		logger.Lava(types.V("0.2.1"), "Using not standarized error")
		return types.ResultErr[FormAnswerDB](errors.New("Invalid IDQuestionType"))
	}

	for i, answer := range o.Answers {
		questionID := answer.IDQuestion

		if !ID.Ensure(questionID, &obj.Answers[i].IDQuestion, "IDForm") {
			logger.Lava(types.V("0.2.1"), "Using not standarized error")
			return types.ResultErr[FormAnswerDB](errors.New("Invalid IDQuestionType"))
		}
	}

	return types.ResultOk(obj)
}
