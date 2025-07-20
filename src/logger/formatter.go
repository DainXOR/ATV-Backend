package logger

import (
	"fmt"
	"strings"
)

type Formatter interface {
	Format(record *Record) (string, error)
	New() Formatter
}

type SimpleFormatter struct {
	timeFormat string
}

func (SimpleFormatter) New() Formatter {
	return &SimpleFormatter{
		timeFormat: "02/01/2006 15:04:05 -07:00",
	}
}
func (f *SimpleFormatter) Format(record *Record) (string, error) {
	formatTime := record.Time.Format(f.timeFormat)
	formatFile := fmt.Sprint(record.Context["file"]) + ":" + fmt.Sprint(record.Context["line"])

	formatArgs := "["
	for _, arg := range record.Args {
		formatArgs += fmt.Sprintf("%v, ", arg)
	}
	formatArgs = strings.TrimSuffix(formatArgs, ", ")
	formatArgs += "]"

	formatCtx := "{"
	for k, v := range record.Context {
		formatCtx += fmt.Sprintf("%s: %v, ", k, v)
	}
	formatCtx = strings.TrimSuffix(formatCtx, ", ")
	formatCtx += "}"

	return strings.TrimSpace(fmt.Sprintf("|%s| %s %s: %s %s %s",
		record.Level,
		formatTime,
		formatFile,
		record.Message,
		formatArgs,
		formatCtx,
	)), nil
}

var _ Formatter = (*SimpleFormatter)(nil)
