package models

type CompanionDB struct {
	ID               DBID       `json:"_id,omitempty" bson:"_id,omitempty"`
	NumberID         string     `json:"number_id,omitempty" bson:"number_id,omitempty"`
	FirstName        string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName         string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email            string     `json:"email,omitempty" bson:"email,omitempty"`
	InstitutionEmail string     `json:"institution_email,omitempty" bson:"institution_email,omitempty"`
	PhoneNumber      string     `json:"phone_number" bson:"phone_number"`
	IDSpeciality     DBID       `json:"id_speciality" bson:"id_speciality"`
	CreatedAt        DBDateTime `json:"created_at,omitzero" bson:"created_at,omitempty"`
	UpdatedAt        DBDateTime `json:"updated_at,omitzero" bson:"updated_at,omitempty"`
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

func (c CompanionCreate) ToInsert() CompanionDB {
	obj := CompanionDB{
		NumberID:         c.NumberID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		InstitutionEmail: c.InstitutionEmail,
		PhoneNumber:      c.PhoneNumber,
		CreatedAt:        Time.Now(),
		UpdatedAt:        Time.Now(),
		DeletedAt:        Time.Zero(),
	}

	if !EnsureID(c.IDSpeciality, &obj.IDSpeciality, "IDSpeciality") {
		return CompanionDB{}
	}

	return obj
}
func (c CompanionCreate) ToUpdate() CompanionDB {
	obj := CompanionDB{
		NumberID:         c.NumberID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		InstitutionEmail: c.InstitutionEmail,
		PhoneNumber:      c.PhoneNumber,
		UpdatedAt:        Time.Now(),
	}

	if !OmitEmptyID(c.IDSpeciality, &obj.IDSpeciality, "IDSpeciality") {
		return CompanionDB{}
	}

	return obj
}
func (c CompanionDB) ToResponse() CompanionResponse {
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

func (c CompanionDB) IsEmpty() bool {
	return c == (CompanionDB{})
}

func (CompanionDB) TableName() string {
	return "companions"
}

var _ DBModelInterface = (*CompanionDB)(nil)
