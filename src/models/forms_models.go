package models

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"errors"
)

type position = uint8
type section = uint8

type questionInfo[ID any] = types.Triplet[position, section, ID]

type FormDB struct {
	ID          DBID                 `json:"id" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name,omitempty"`
	Description string               `json:"description" bson:"description,omitempty"`
	Questions   []questionInfo[DBID] `json:"questions" bson:"questions,omitempty"`
	CreatedAt   DBDateTime           `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt   DBDateTime           `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt   DBDateTime           `json:"deleted_at" bson:"deleted_at"`
}
type FormCreate struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Questions   []questionInfo[string] `json:"questions"`
}
type FormResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Questions   []questionInfo[string] `json:"questions"`
	CreatedAt   DBDateTime             `json:"created_at"`
	UpdatedAt   DBDateTime             `json:"updated_at"`
}

func (q FormCreate) ToInsert() types.Result[FormDB] {
	obj := FormDB{
		Name:        q.Name,
		Description: q.Description,
		Questions:   nil,
		CreatedAt:   Time.Now(),
		UpdatedAt:   Time.Now(),
		DeletedAt:   Time.Zero(),
	}

	l := logger.Lava(types.V("0.2.0"), "Using not standarized error")
	l.LavaStart()
	questions, err := utils.MapE(q.Questions, func(v questionInfo[string]) (questionInfo[DBID], error) {
		var oid DBID

		if !ID.Ensure(v.Third, &oid, "IDQuestion") {
			return types.NewTriplet(v.First, v.Second, oid), errors.New("Invalid IDQuestion: " + v.Third)
		}

		return types.NewTriplet(v.First, v.Second, oid), nil
	})
	l.LavaEnd()

	if err != nil {
		return types.ResultErr[FormDB](err)
	}

	obj.Questions = questions
	return types.ResultOk(obj)
}
func (q FormDB) ToUpdate() FormDB {
	return FormDB{
		Name:        q.Name,
		Description: q.Description,
		Questions:   q.Questions,
		UpdatedAt:   Time.Now(),
	}
}
func (q FormDB) ToResponse() FormResponse {
	return FormResponse{
		ID:          q.ID.Hex(),
		Name:        q.Name,
		Description: q.Description,
		Questions: utils.Map(q.Questions, func(v questionInfo[DBID]) questionInfo[string] {
			return types.NewTriplet(v.First, v.Second, v.Third.Hex())
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
		(len(q.Questions) == 0 ||
			utils.All(q.Questions, func(q questionInfo[DBID]) bool {
				return q == (questionInfo[DBID]{})
			}))

	if comp {
		for i, o := range q.Questions {
			comp = comp && (o == zeroObj.Questions[i])
		}
	}

	return comp
}

func (FormDB) TableName() string {
	return "forms"
}

var _ DBModelInterface = (*FormDB)(nil)
