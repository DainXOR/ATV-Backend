package logger

import "time"

type Record struct {
	LogLevel   logLevel
	Time       time.Time
	Message    string
	File       string
	Line       int
	AppVersion string
	Context    map[string]string
}

type stringRecord struct {
	LogLevel   string
	Time       string
	Message    string
	File       string
	Line       string
	AppVersion string
	Context    map[string]string
}
