package models

type FormQuestionTypeDB struct {
	ID        DBID       `bson:"_id,omitempty"`
	Name      string     `bson:"name,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt DBDateTime `json:"deleted_at" bson:"deleted_at"`
}
type FormQuestionTypeCreate struct {
	Name string `json:"name"`
}
type FormQuestionTypeResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt DBDateTime `json:"created_at"`
	UpdatedAt DBDateTime `json:"updated_at"`
}

func (t FormQuestionTypeCreate) ToInsert() FormQuestionTypeDB {
	return FormQuestionTypeDB{
		Name:      t.Name,
		CreatedAt: Time.Now(),
		UpdatedAt: Time.Now(),
		DeletedAt: Time.Zero(),
	}
}
func (t FormQuestionTypeDB) ToUpdate() FormQuestionTypeDB {
	return FormQuestionTypeDB{
		Name:      t.Name,
		UpdatedAt: Time.Now(),
	}
}
func (t FormQuestionTypeDB) ToResponse() FormQuestionTypeResponse {
	return FormQuestionTypeResponse{
		ID:        t.ID.Hex(),
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func (t FormQuestionTypeDB) IsEmpty() bool {
	return t == (FormQuestionTypeDB{})
}

func (t FormQuestionTypeDB) TableName() string {
	return "form_question_types"
}

var _ DBModelInterface = (*FormQuestionTypeDB)(nil)
