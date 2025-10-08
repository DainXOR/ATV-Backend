package models

type ContactReasonDB struct {
	ID        DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt DBDateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt DBDateTime `json:"deleted_at" bson:"deleted_at"`
}

type ContactReasonCreate struct {
	Name string `json:"name"`
}

type ContactReasonResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt DBDateTime `json:"created_at"`
	UpdatedAt DBDateTime `json:"updated_at"`
}

func (p ContactReasonCreate) ToInsert() ContactReasonDB {
	return ContactReasonDB{
		Name:      p.Name,
		CreatedAt: Time.Now(),
		UpdatedAt: Time.Now(),
		DeletedAt: Time.Zero(),
	}
}
func (p ContactReasonCreate) ToUpdate() ContactReasonDB {
	return ContactReasonDB{
		Name:      p.Name,
		UpdatedAt: Time.Now(),
	}
}

func (p ContactReasonDB) ToResponse() ContactReasonResponse {
	return ContactReasonResponse{
		ID:        p.ID.Hex(),
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (p ContactReasonDB) IsEmpty() bool {
	return p.ID == (DBID{})
}

func (ContactReasonDB) TableName() string {
	return "contact_reasons"
}

var _ DBModelInterface = (*ContactReasonDB)(nil)
