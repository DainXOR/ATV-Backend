package logger

import (
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"strconv"
	"time"
)

type Record struct {
	LogLevel   logLevel
	Time       time.Time
	Message    string
	File       string
	Line       uint64
	AppVersion types.Version
	Context    map[string]string
}

func NewRecord(msg string, extra ...types.SPair[string]) Record {
	rec := Record{
		Time:    time.Now(),
		Message: msg,
	}

	rec.File, rec.Line = utils.CallOrigin(2)
	rec.AppVersion = types.V0()
	rec.Context = make(map[string]string, len(extra))

	if len(extra) > 0 {
		for _, pair := range extra {
			if internal.AppVersion().Check(pair.First) {
				// If the key is app_version, convert the value to a Version type
				if version, err := types.VersionFrom(pair.Second); err == nil {
					rec.AppVersion = version
				}
				continue
			}
			if internal.CallOriginOffset().Check(pair.First) {
				if i, err := strconv.Atoi(pair.Second); err == nil {
					rec.File, rec.Line = utils.CallOrigin(2 + int(i))
				} else {
					rec.File, rec.Line = "UnknownFile", 0
				}
				continue
			}

			rec.Context[pair.First] = pair.Second
		}
	}

	return rec
}

type FormatRecord struct {
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
