package models

type CompanionDBMongo struct {
	ID               DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	NumberID         string     `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email            string     `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string     `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	PhoneNumber      string     `json:"phone_number" bson:"phone_number"`
	IDSpeciality     DBID       `json:"id_speciality" bson:"id_speciality"`
	CreatedAt        DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt        DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
	DeletedAt        DBDateTime `json:"deleted_at" bson:"deleted_at"`
}
type CompanionCreate struct {
	NumberID         string `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email            string `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	PhoneNumber      string `json:"phone_number" bson:"phone_number"`
	IDSpeciality     string `json:"id_speciality" bson:"id_speciality"`
}
type CompanionResponse struct {
	ID               string     `json:"id,omitempty" bson:"id,omitempty"`
	NumberID         string     `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email            string     `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string     `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	PhoneNumber      string     `json:"phone_number" bson:"phone_number"`
	IDSpeciality     string     `json:"id_speciality" bson:"id_speciality"`
	CreatedAt        DBDateTime `json:"created_at,omitzero" bson:"created_at,omitzero"`
	UpdatedAt        DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitzero"`
}

func (c CompanionCreate) ToInsert() CompanionDBMongo {
	obj := CompanionDBMongo{
		NumberID:         c.NumberID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		InstitutionEmail: c.InstitutionEmail,
		PhoneNumber:      c.PhoneNumber,
		CreatedAt:        TimeNow(),
		UpdatedAt:        TimeNow(),
		DeletedAt:        TimeZero(),
	}

	if !EnsureID(c.IDSpeciality, &obj.IDSpeciality, "IDSpeciality") {
		return CompanionDBMongo{}
	}

	return obj
}
func (c CompanionCreate) ToUpdate() CompanionDBMongo {
	obj := CompanionDBMongo{
		NumberID:         c.NumberID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		InstitutionEmail: c.InstitutionEmail,
		PhoneNumber:      c.PhoneNumber,
		CreatedAt:        TimeNow(),
	}

	if !OmitEmptyID(c.IDSpeciality, &obj.IDSpeciality, "IDSpeciality") {
		return CompanionDBMongo{}
	}

	return obj
}
func (c CompanionDBMongo) ToResponse() CompanionResponse {
	return CompanionResponse{
		ID:               c.ID.Hex(),
		NumberID:         c.NumberID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		InstitutionEmail: c.InstitutionEmail,
		PhoneNumber:      c.PhoneNumber,
		IDSpeciality:     c.IDSpeciality.Hex(),
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}
func (c CompanionDBMongo) IsEmpty() bool {
	return c == (CompanionDBMongo{})
}
func (c *CompanionDBMongo) SetID(id any) error {
	var err error
	c.ID, err = ID.ToDB(id)
	return err
}

func (CompanionDBMongo) TableName() string {
	return "companions"
}

var _ DBModelInterface = (*CompanionDBMongo)(nil)
