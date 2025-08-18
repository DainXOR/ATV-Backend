package models

type UniversityDB struct {
	ID       DBID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Location string `json:"location,omitempty" bson:"location,omitempty"`
	DBModelBase
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

func (u UniversityCreate) ToInsert() UniversityDB {
	return UniversityDB{
		Name:     u.Name,
		Location: u.Location,
		DBModelBase: DBModelBase{
			CreatedAt: Time.Now(),
			UpdatedAt: Time.Now(),
			DeletedAt: Time.Zero(),
		},
	}
}
func (u UniversityCreate) ToUpdate() UniversityDB {
	return UniversityDB{
		Name:     u.Name,
		Location: u.Location,
		DBModelBase: DBModelBase{
			UpdatedAt: Time.Now(),
		},
	}
}

func (u UniversityDB) ToResponse() UniversityResponse {
	return UniversityResponse{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		Location:  u.Location,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func (u UniversityDB) IsEmpty() bool {
	return u == (UniversityDB{})
}

func (UniversityDB) TableName() string {
	return "universities"
}

var _ DBModelInterface = (*UniversityDB)(nil)
