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
		//AppVersion: types.V(configs.App.ApiVersion()),
	}

	if len(extra) > 0 {
		if e := extra[0]; internal.CallOriginOffset().Check(e.First) {
			i, err := strconv.Atoi(e.Second)

			if err == nil {
				rec.File, rec.Line = utils.CallOrigin(int(i) + 2)
			}

			extra = extra[1:] // Remove the first element

		} else {
			rec.File, rec.Line = utils.CallOrigin(2)
		}

		rec.Context = make(map[string]string, len(extra))
		for _, pair := range extra {
			if internal.AppVersion().Check(pair.First) {
				// If the key is app_version, convert the value to a Version type
				if version, err := types.VersionFrom(pair.Second); err == nil {
					rec.AppVersion = version
				} else {
					rec.AppVersion = types.V("0.0.0") // Default if parsing fails
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
