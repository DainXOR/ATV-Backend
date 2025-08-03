package logger

import (
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"time"
)

type Record struct {
	LogLevel   logLevel
	Time       time.Time
	Message    string
	File       string
	Line       int
	AppVersion string
	Context    map[string]string
}

func NewRecord(msg string, extra ...types.SPair[string]) Record {
	rec := Record{
		Time:       time.Now(),
		Message:    msg,
		AppVersion: AppVersion(),
	}

	rec.File, rec.Line = utils.CallOrigin(2)

	if len(extra) > 0 {
		rec.Context = make(map[string]string, len(extra))
		for _, pair := range extra {
			rec.Context[pair.First] = pair.Second
		}
	}

	return rec
}

type formatRecord struct {
	LogLevel     string
	Time         string
	File         string
	Line         string
	Message      string
	AppVersion   string
	Context      map[string](types.SPair[string])
	ContextBegin string
	ContextEnd   string
}
