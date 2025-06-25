package models

type UniversityDBMongo struct {
	ID        DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	Location  string     `json:"location,omitempty" bson:"location,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
	DeletedAt DBDateTime `json:"deleted_at,omitzero" bson:"deleted_at,omitempty"`
}

type UniversityDBMongoReceiver struct {
	ID        any        `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	Location  string     `json:"location,omitempty" bson:"location,omitempty"`
	CreatedAt DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
	DeletedAt DBDateTime `json:"deleted_at,omitzero" bson:"deleted_at,omitzero"`
}

func (UniversityDBMongo) Receiver() UniversityDBMongoReceiver {
	return UniversityDBMongoReceiver{}
}
func (UniversityDBMongo) ReceiverList() []UniversityDBMongoReceiver {
	u := make([]UniversityDBMongoReceiver, 1)
	u[0] = UniversityDBMongo{}.Receiver()
	return u
}

// UniversityCreate represents the request body for creating a new university
type UniversityCreate struct {
	Name     string `json:"name" gorm:"not null"`
	Location string `json:"location" gorm:"not null"`
}

// UniversityResponse represents the response body for a university
type UniversityResponse struct {
	ID        string     `json:"_id,omitempty" bson:"_id,omitempty"`
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
	}
}
func (u UniversityCreate) ToUpdate() UniversityDBMongo {
	return UniversityDBMongo{
		Name:      u.Name,
		Location:  u.Location,
		UpdatedAt: TimeNow(),
	}
}

func (u UniversityDBMongoReceiver) ToDB() UniversityDBMongo {
	id, _ := DBIDFrom(u.ID)
	return UniversityDBMongo{
		ID:        id,
		Name:      u.Name,
		Location:  u.Location,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
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
func (UniversityDBMongoReceiver) TableName() string {
	return UniversityDBMongo{}.TableName()
}

var _ DBModelInterface = (*UniversityDBMongo)(nil)
var _ DBModelInterface = (*UniversityDBMongoReceiver)(nil)
