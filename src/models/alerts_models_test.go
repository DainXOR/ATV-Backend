package models

import (
	"dainxor/atv/utils"
	"testing"
)

func TestAlertDB_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		input AlertDB
		want  bool
	}{
		{"empty struct", AlertDB{}, true},
		{"non-empty struct", AlertDB{ID: utils.Test.GenerateObjectID()}, false},
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
