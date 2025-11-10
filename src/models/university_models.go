package models

type UniversityDB struct {
	ID        DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	Location  string     `json:"location,omitempty" bson:"location,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt DBDateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt DBDateTime `json:"deleted_at" bson:"deleted_at"`
}

// UniversityCreate represents the request body for creating a new university
type UniversityCreate struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

// UniversityResponse represents the response body for a university
type UniversityResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Location  string     `json:"location"`
	CreatedAt DBDateTime `json:"created_at"`
	UpdatedAt DBDateTime `json:"updated_at"`
}

func (u UniversityCreate) ToInsert() UniversityDB {
	return UniversityDB{
		Name:      u.Name,
		Location:  u.Location,
		CreatedAt: Time.Now(),
		UpdatedAt: Time.Now(),
		DeletedAt: Time.Zero(),
	}
}
func (u UniversityCreate) ToUpdate() UniversityDB {
	return UniversityDB{
		Name:      u.Name,
		Location:  u.Location,
		UpdatedAt: Time.Now(),
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
