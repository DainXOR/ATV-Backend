package models

import (
	"dainxor/atv/utils"
	"testing"
)

func TestPriorityDB_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		input PriorityDB
		want  bool
	}{
		{"empty struct", PriorityDB{}, true},
		{"non-empty struct", PriorityDB{ID: utils.Test.GenerateObjectID()}, false},
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

func TestPriorityDB_ToResponse(t *testing.T) {
	p := PriorityDB{
		ID:               utils.Test.GenerateObjectID(),
		Name:             "Test",
		Level:            1,
		SessionsPerMonth: 2,
		CreatedAt:        Time.Now(),
		UpdatedAt:        Time.Now(),
	}
	r := p.ToResponse()
	if r.Name != p.Name || r.Level != p.Level || r.SessionsPerMonth != p.SessionsPerMonth {
		t.Errorf("ToResponse() did not copy fields correctly")
	}
}

func TestPriorityCreate_ToInsertAndToUpdate(t *testing.T) {
	c := PriorityCreate{Name: "Test", Level: 1, SessionsPerMonth: 2}
	pInsert := c.ToInsert()
	if pInsert.Name != c.Name || pInsert.Level != c.Level || pInsert.SessionsPerMonth != c.SessionsPerMonth {
		t.Errorf("ToInsert() did not copy fields correctly")
	}
	pUpdate := c.ToUpdate()
	if pUpdate.Name != c.Name || pUpdate.Level != c.Level || pUpdate.SessionsPerMonth != c.SessionsPerMonth {
		t.Errorf("ToUpdate() did not copy fields correctly")
	}
}
