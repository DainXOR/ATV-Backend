package models

import (
	"dainxor/atv/utils"
	"testing"
)

func TestSpecialityDB_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		input SpecialityDB
		want  bool
	}{
		{"empty struct", SpecialityDB{}, true},
		{"non-empty struct", SpecialityDB{ID: utils.Test.GenerateObjectID()}, false},
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

// Add more tests for conversion and field logic as needed
