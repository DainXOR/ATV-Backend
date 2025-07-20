package logger

import "time"

type Record struct {
	Level   logLevel
	Time    time.Time
	Message string
	Args    []any
	// File       string
	// Line       int
	// AppVersion string
	Context map[string]any
}
