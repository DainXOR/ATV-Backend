package models

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"errors"
)

//type position = uint8
//type section = uint8

type questionInfo[ID any] struct {
	Position         uint8 `json:"position" bson:"position"`
	Section          uint8 `json:"section" bson:"section"`
	IDParentQuestion ID    `json:"id_parent_question" bson:"id_parent_question"`
	IDQuestion       ID    `json:"id_question" bson:"id_question"`
	Optional         bool  `json:"optional" bson:"optional"`
}

func questionFrom[IID, OID any](other questionInfo[IID], transformFunc func(IID) OID) questionInfo[OID] {
	return questionInfo[OID]{
		IDQuestion:       transformFunc(other.IDQuestion),
		Position:         other.Position,
		Section:          other.Section,
		IDParentQuestion: transformFunc(other.IDParentQuestion),
		Optional:         other.Optional,
	}
}
func questionEmpty[ID any]() questionInfo[ID] {
	return questionInfo[ID]{}
}

// type questionInfo[ID any] = types.Triplet[position, section, ID]

type FormDB struct {
	ID            DBID                 `json:"id" bson:"_id,omitempty"`
	Name          string               `json:"name" bson:"name,omitempty"`
	Description   string               `json:"description" bson:"description,omitempty"`
	QuestionsInfo []questionInfo[DBID] `json:"questions_info" bson:"questions_info,omitempty"`
	CreatedAt     DBDateTime           `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt     DBDateTime           `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt     DBDateTime           `json:"deleted_at" bson:"deleted_at"`
}
type FormCreate struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	QuestionsInfo []questionInfo[string] `json:"questions_info"`
}
type FormResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	QuestionsInfo []questionInfo[string] `json:"questions_info"`
	CreatedAt     DBDateTime             `json:"created_at"`
	UpdatedAt     DBDateTime             `json:"updated_at"`
}

func (q FormCreate) ToInsert() types.Result[FormDB] {
	obj := FormDB{
		Name:          q.Name,
		Description:   q.Description,
		QuestionsInfo: nil,
		CreatedAt:     Time.Now(),
		UpdatedAt:     Time.Now(),
		DeletedAt:     Time.Zero(),
	}

	l := logger.Lava(types.V("0.2.0"), "Using not standarized error")
	l.LavaStart()
	questions, err := utils.MapE(q.QuestionsInfo, func(v questionInfo[string]) (questionInfo[DBID], error) {
		var oid DBID
		var oidParent DBID

		if !ID.Ensure(v.IDQuestion, &oid, "IDQuestion") {
			return questionEmpty[DBID](), errors.New("Invalid IDQuestion: " + v.IDQuestion)
		}
		if !ID.OmitEmpty(v.IDParentQuestion, &oidParent, "IDParentQuestion") {
			return questionEmpty[DBID](), errors.New("Invalid IDQuestion: " + v.IDParentQuestion)
		}

		return questionFrom(v, func(id string) DBID {
			oid, _ := ID.ToDB(id)
			return oid
		}), nil
	})
	l.LavaEnd()

	if err != nil {
		obj.QuestionsInfo = make([]questionInfo[DBID], 0)
		return types.ResultErr[FormDB](err)
	}

	obj.QuestionsInfo = questions
	return types.ResultOk(obj)
}
func (q FormDB) ToUpdate() FormDB {
	return FormDB{
		Name:          q.Name,
		Description:   q.Description,
		QuestionsInfo: q.QuestionsInfo,
		UpdatedAt:     Time.Now(),
	}
}
func (q FormDB) ToResponse() FormResponse {
	return FormResponse{
		ID:          q.ID.Hex(),
		Name:        q.Name,
		Description: q.Description,
		QuestionsInfo: utils.Map(q.QuestionsInfo, func(v questionInfo[DBID]) questionInfo[string] {
			return questionFrom(v, DBID.Hex)
		}),
		UpdatedAt: q.UpdatedAt,
		CreatedAt: q.CreatedAt,
	}
}

func (q FormDB) IsEmpty() bool {
	zeroObj := FormDB{}

	comp := q.ID == zeroObj.ID &&
		q.Name == zeroObj.Name &&
		q.Description == zeroObj.Description &&
		q.CreatedAt.Equal(zeroObj.CreatedAt) &&
		q.UpdatedAt.Equal(zeroObj.UpdatedAt) &&
		q.DeletedAt.Equal(zeroObj.DeletedAt) &&
		(len(q.QuestionsInfo) == 0 ||
			utils.All(q.QuestionsInfo, func(q questionInfo[DBID]) bool {
				return q == (questionInfo[DBID]{})
			}))

	if comp {
		for i, o := range q.QuestionsInfo {
			comp = comp && (o == zeroObj.QuestionsInfo[i])
		}
	}

	return comp
}

func (FormDB) TableName() string {
	return "forms"
}

var _ DBModelInterface = (*FormDB)(nil)
