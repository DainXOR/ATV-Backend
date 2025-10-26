package models

import (
	"dainxor/atv/utils"
	"testing"
)

func TestStudentDB_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		input StudentDB
		want  bool
	}{
		{"empty struct", StudentDB{}, true},
		{"non-empty struct", StudentDB{ID: utils.Test.GenerateObjectID()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.IsZero()
			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStudentDB_ToResponse(t *testing.T) {
	s := StudentDB{
		ID:               utils.Test.GenerateObjectID(),
		NumberID:         "123",
		FirstName:        "John",
		LastName:         "Doe",
		PersonalEmail:    "john@example.com",
		InstitutionEmail: "john@uni.edu",
		ResidenceAddress: "123 St",
		Semester:         2,
		PhoneNumber:      "555-1234",
		CreatedAt:        Time.Now(),
		UpdatedAt:        Time.Now(),
	}
	r := s.ToResponse()
	if r.FirstName != s.FirstName || r.LastName != s.LastName || r.NumberID != s.NumberID {
		t.Errorf("ToResponse() did not copy fields correctly")
	}
}

func TestStudentCreate_ToInsert(t *testing.T) {
	c := StudentCreate{NumberID: "123", FirstName: "John", LastName: "Doe", PersonalEmail: "john@example.com", InstitutionEmail: "john@uni.edu", ResidenceAddress: "123 St", Semester: 2, IDUniversity: utils.Test.GenerateObjectID().Hex(), PhoneNumber: "555-1234"}
	s := c.ToInsert()
	if s.FirstName != c.FirstName || s.LastName != c.LastName || s.NumberID != c.NumberID {
		t.Errorf("ToInsert() did not copy fields correctly")
	}
}
