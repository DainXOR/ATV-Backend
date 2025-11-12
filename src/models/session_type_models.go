package models

type SessionTypeDB struct {
	ID        DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt DBDateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt DBDateTime `json:"deleted_at" bson:"deleted_at"`
}

// SessionTypeCreate represents the request body for creating a new SessionType
type SessionTypeCreate struct {
	Name string `json:"name"`
}

// SessionTypeResponse represents the response body for a SessionType
type SessionTypeResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt DBDateTime `json:"created_at"`
	UpdatedAt DBDateTime `json:"updated_at"`
}

func (u SessionTypeCreate) ToInsert() SessionTypeDB {
	return SessionTypeDB{
		Name:      u.Name,
		CreatedAt: Time.Now(),
		UpdatedAt: Time.Now(),
		DeletedAt: Time.Zero(),
	}
}
func (u SessionTypeCreate) ToUpdate() SessionTypeDB {
	return SessionTypeDB{
		Name:      u.Name,
		UpdatedAt: Time.Now(),
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
func (u SessionTypeDB) IsZero() bool {
	return u == (SessionTypeDB{})
}

func (SessionTypeDB) TableName() string {
	return "session_types"
}

var _ DBModelInterface = (*SessionTypeDB)(nil)
