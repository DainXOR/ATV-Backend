package logger

import (
	"dainxor/atv/utils"
	"fmt"
	"strings"
)

type Formatter interface {
	Format(record *Record) (string, error)
	FormatStrings(record *StringRecord) (string, error)
	DateFormat() string
}
type FormatterBuilder interface {
	New() Formatter
}

type simpleFormatter struct {
	dateFormat string
}

func (f *simpleFormatter) Format(record *Record) (string, error) {
	formatTime := record.Time.Format(f.DateFormat())
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
func (f *simpleFormatter) FormatStrings(record *StringRecord) (string, error) {
	return fmt.Sprintf("|%s| %s %s: %s %s",
		record.Level,
		record.Time,
		record.File,
		record.Message,
		record.Context,
	), nil
}
func (f *simpleFormatter) DateFormat() string {
	return f.dateFormat
}

type simpleFormatterBuilder struct {
	formatter simpleFormatter
}

var SimpleFormatter simpleFormatterBuilder = simpleFormatterBuilder{
	formatter: simpleFormatter{
		dateFormat: "",
	},
}

func (b simpleFormatterBuilder) TimeFormat(format string) simpleFormatterBuilder {
	b.formatter.dateFormat = format
	return b
}
func (b simpleFormatterBuilder) New() Formatter {
	if b.formatter.dateFormat == "" {
		b.formatter.dateFormat = "02/01/2006 15:04:05 -07:00"
	}

	return &simpleFormatter{
		dateFormat: b.formatter.dateFormat,
	}
}

type consoleColorFormatter struct {
	formatter   Formatter
	colorScheme AnsiColorScheme
}

// DateFormat implements Formatter.
func (f consoleColorFormatter) DateFormat() string {
	panic("unimplemented")
}

// FormatStrings implements Formatter.
func (f consoleColorFormatter) FormatStrings(record *StringRecord) (string, error) {
	panic("unimplemented")
}

func (f *consoleColorFormatter) Format(record *Record) (string, error) {
	levelString := record.Level.String()
	timeString := record.Time.Format(f.DateFormat())
	lineString := fmt.Sprint(record.Line)

	styleLevel := f.colorScheme.GetStyle(record.Level.String())
	styleTime := f.colorScheme.GetStyle("time")
	styleMessage := f.colorScheme.GetStyle("message")
	styleFile := f.colorScheme.GetStyle("file")
	styleLine := f.colorScheme.GetStyle("line")
	styleVersion := f.colorScheme.GetStyle("version")
	StyleCtxKey := f.colorScheme.GetStyle("context-key")
	StyleCtxValue := f.colorScheme.GetStyle("context-value")

	StyledCtx := utils.DFlatten(record.Context, func(k, v string) string {
		return fmt.Sprintf("%s: %s", StyleCtxKey.Apply(k), StyleCtxValue.Apply(v))
	})

	recordString := StringRecord{
		Level:      styleLevel.Apply(levelString),
		Time:       styleTime.Apply(timeString),
		Message:    styleMessage.Apply(record.Message),
		File:       styleFile.Apply(record.File),
		Line:       styleLine.Apply(lineString),
		AppVersion: styleVersion.Apply(record.AppVersion),
		Context:    utils.Reduce(StyledCtx, func(acc, ctx string) string { return acc + ", " + ctx }, ""),
	}

	formatted, err := f.formatter.FormatStrings(&recordString)
	if err != nil {
		return "", err
	}

	return formatted, nil
}

type colorFormatterBuilder struct {
	formatter consoleColorFormatter
}

func (b colorFormatterBuilder) Formatter(formatter Formatter) colorFormatterBuilder {
	b.formatter.formatter = formatter
	return b
}
func (b colorFormatterBuilder) AddColor(level logLevel, colorCode string) colorFormatterBuilder {
	//if b.formatter.colors == nil {
	//	b.formatter.colors = make(map[logLevel]string)
	//}
	//b.formatter.colors[level] = colorCode
	return b
}
func (b colorFormatterBuilder) DefaultColor(colorCode string) colorFormatterBuilder {
	//if b.formatter.colors == nil {
	//	b.formatter.colors = make(map[logLevel]string)
	//}
	//
	//b.formatter.colors[LEVEL_NONE] = colorCode
	return b
}
func (b colorFormatterBuilder) New() Formatter {
	//if b.formatter.formatter == nil {
	//	b.formatter.formatter = SimpleFormatter.New()
	//}
	//
	//if len(b.formatter.colors) == 0 {
	//	b.formatter.colors = map[logLevel]string{
	//		LEVEL_DEBUG:   "\033[34m", // Blue
	//		LEVEL_INFO:    "\033[32m", // Green
	//		LEVEL_WARNING: "\033[33m", // Yellow
	//		LEVEL_ERROR:   "\033[31m", // Red
	//		LEVEL_FATAL:   "\033[35m", // Magenta
	//	}
	//}

	return &consoleColorFormatter{
		formatter:   b.formatter.formatter,
		colorScheme: b.formatter.colorScheme,
	}
}

var _ Formatter = (*simpleFormatter)(nil)
var _ Formatter = (*consoleColorFormatter)(nil)
var _ FormatterBuilder = (*simpleFormatterBuilder)(nil)
var _ FormatterBuilder = (*colorFormatterBuilder)(nil)
