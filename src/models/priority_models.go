package models

type PriorityDB struct {
	ID               DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Name             string     `json:"name,omitempty" bson:"name,omitempty"`
	Level            uint8      `json:"level,omitempty" bson:"level,omitempty"`
	SessionsPerMonth uint8      `json:"sessions_per_month" bson:"sessions_per_month"`
	CreatedAt        DBDateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt        DBDateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt        DBDateTime `json:"deleted_at" bson:"deleted_at"`
}

type PriorityCreate struct {
	Name             string `json:"name"`
	Level            uint8  `json:"level"`
	SessionsPerMonth uint8  `json:"sessions_per_month"`
}

type PriorityResponse struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Level            uint8      `json:"level"`
	SessionsPerMonth uint8      `json:"sessions_per_month"`
	CreatedAt        DBDateTime `json:"created_at"`
	UpdatedAt        DBDateTime `json:"updated_at"`
}

func (p PriorityCreate) ToInsert() PriorityDB {
	return PriorityDB{
		Name:             p.Name,
		Level:            p.Level,
		SessionsPerMonth: p.SessionsPerMonth,
		CreatedAt:        Time.Now(),
		UpdatedAt:        Time.Now(),
		DeletedAt:        Time.Zero(),
	}
}
func (p PriorityCreate) ToUpdate() PriorityDB {
	return PriorityDB{
		Name:             p.Name,
		Level:            p.Level,
		SessionsPerMonth: p.SessionsPerMonth,
		UpdatedAt:        Time.Now(),
	}
}

func (p PriorityDB) ToResponse() PriorityResponse {
	return PriorityResponse{
		ID:               p.ID.Hex(),
		Name:             p.Name,
		Level:            p.Level,
		SessionsPerMonth: p.SessionsPerMonth,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}

func (p PriorityDB) IsEmpty() bool {
	return p.ID == (DBID{})
}

func (PriorityDB) TableName() string {
	return "priorities"
}

var _ DBModelInterface[PriorityCreate, PriorityResponse] = (*PriorityDB)(nil)
