package logger

import "time"

type Record struct {
	Level      logLevel
	Time       time.Time
	Message    string
	File       string
	Line       int
	AppVersion string
	Args       []any
	Context    map[string]string
}

type StringRecord struct {
	Level      string
	Time       string
	Message    string
	File       string
	Line       string
	AppVersion string
	Context    string
}
