package models

type SessionTypeDB struct {
	ID        DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt DBDateTime `json:"deleted_at,omitzero" bson:"deleted_at,omitempty"`
}

// SessionTypeCreate represents the request body for creating a new SessionType
type SessionTypeCreate struct {
	Name string `json:"name" gorm:"not null"`
}

// SessionTypeResponse represents the response body for a SessionType
type SessionTypeResponse struct {
	ID        string     `json:"id,omitempty" bson:"id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
}

func (u SessionTypeCreate) ToInsert() SessionTypeDB {
	return SessionTypeDB{
		Name:      u.Name,
		CreatedAt: TimeNow(),
		UpdatedAt: TimeNow(),
		DeletedAt: TimeZero(),
	}
}
func (u SessionTypeCreate) ToUpdate() SessionTypeDB {
	return SessionTypeDB{
		Name:      u.Name,
		UpdatedAt: TimeNow(),
	}
}
func (u SessionTypeDB) ToResponse() SessionTypeResponse {
	return SessionTypeResponse{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func (u SessionTypeDB) IsEmpty() bool {
	return u == (SessionTypeDB{})
}

func (SessionTypeDB) TableName() string {
	return "session_types"
}

var _ DBModelInterface = (*SessionTypeDB)(nil)
