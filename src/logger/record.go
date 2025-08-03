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

type formatRecord struct {
	LogLevel     string
	Time         string
	File         string
	Line         string
	Message      string
	AppVersion   string
	Context      map[string]string
	ContextBegin string
	ContextEnd   string
}
