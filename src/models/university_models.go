package models

type UniversityDBMongo struct {
	ID        DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	Location  string     `json:"location,omitempty" bson:"location,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt DBDateTime `json:"deleted_at,omitzero" bson:"deleted_at,omitempty"`
}

// UniversityCreate represents the request body for creating a new university
type UniversityCreate struct {
	Name     string `json:"name" gorm:"not null"`
	Location string `json:"location" gorm:"not null"`
}

// UniversityResponse represents the response body for a university
type UniversityResponse struct {
	ID        string     `json:"id,omitempty" bson:"id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	Location  string     `json:"location,omitempty" bson:"location,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
}

func (u UniversityCreate) ToInsert() UniversityDBMongo {
	return UniversityDBMongo{
		Name:      u.Name,
		Location:  u.Location,
		CreatedAt: TimeNow(),
		UpdatedAt: TimeNow(),
		DeletedAt: TimeZero(),
	}
}
func (u UniversityCreate) ToUpdate() UniversityDBMongo {
	return UniversityDBMongo{
		Name:      u.Name,
		Location:  u.Location,
		UpdatedAt: TimeNow(),
	}
}

func (u UniversityDBMongo) ToResponse() UniversityResponse {
	return UniversityResponse{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		Location:  u.Location,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (UniversityDBMongo) TableName() string {
	return "universities"
}

var _ DBModelInterface = (*UniversityDBMongo)(nil)
