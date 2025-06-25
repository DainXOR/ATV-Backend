package models

type SpecialityDBMongo struct {
	ID        DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt DBDateTime `json:"deleted_at,omitzero" bson:"deleted_at,omitempty"`
}

type SpecialityDBMongoReceiver struct {
	ID        any        `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
	DeletedAt DBDateTime `json:"deleted_at,omitzero" bson:"deleted_at,omitzero"`
}

func (SpecialityDBMongo) Receiver() SpecialityDBMongoReceiver {
	return SpecialityDBMongoReceiver{}
}
func (SpecialityDBMongo) ReceiverList() []SpecialityDBMongoReceiver {
	u := make([]SpecialityDBMongoReceiver, 1)
	u[0] = SpecialityDBMongo{}.Receiver()
	return u
}

// SpecialityCreate represents the request body for creating a new Speciality
type SpecialityCreate struct {
	Name string `json:"name" gorm:"not null"`
}

// SpecialityResponse represents the response body for a Speciality
type SpecialityResponse struct {
	ID        string     `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
}

func (u SpecialityCreate) ToInsert() SpecialityDBMongo {
	return SpecialityDBMongo{
		Name:      u.Name,
		CreatedAt: TimeNow(),
		UpdatedAt: TimeNow(),
	}
}
func (u SpecialityCreate) ToUpdate() SpecialityDBMongo {
	return SpecialityDBMongo{
		Name:      u.Name,
		UpdatedAt: TimeNow(),
	}
}

func (u SpecialityDBMongoReceiver) ToDB() SpecialityDBMongo {
	id, _ := PrimitiveIDFrom(u.ID)
	return SpecialityDBMongo{
		ID:        id,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func (u SpecialityDBMongo) ToResponse() SpecialityResponse {
	return SpecialityResponse{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (SpecialityDBMongo) TableName() string {
	return "specialities"
}
func (SpecialityDBMongoReceiver) TableName() string {
	return SpecialityDBMongo{}.TableName()
}

var _ DBModelInterface = (*SpecialityDBMongo)(nil)
var _ DBModelInterface = (*SpecialityDBMongoReceiver)(nil)
