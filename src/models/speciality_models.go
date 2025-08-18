package models

type SpecialityDB struct {
	ID        DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt DBDateTime `json:"deleted_at,omitzero" bson:"deleted_at,omitempty"`
}

// SpecialityCreate represents the request body for creating a new Speciality
type SpecialityCreate struct {
	Name string `json:"name" gorm:"not null"`
}

// SpecialityResponse represents the response body for a Speciality
type SpecialityResponse struct {
	ID        string     `json:"id,omitempty" bson:"id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
}

func (u SpecialityCreate) ToInsert() SpecialityDB {
	return SpecialityDB{
		Name:      u.Name,
		CreatedAt: TimeNow(),
		UpdatedAt: TimeNow(),
		DeletedAt: TimeZero(),
	}
}
func (u SpecialityCreate) ToUpdate() SpecialityDB {
	return SpecialityDB{
		Name:      u.Name,
		UpdatedAt: TimeNow(),
	}
}
func (u SpecialityDB) ToResponse() SpecialityResponse {
	return SpecialityResponse{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func (u SpecialityDB) IsEmpty() bool {
	return u == (SpecialityDB{})
}

func (SpecialityDB) TableName() string {
	return "specialities"
}

var _ DBModelInterface = (*SpecialityDB)(nil)
