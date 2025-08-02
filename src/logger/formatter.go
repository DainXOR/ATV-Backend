package logger

import (
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"fmt"
	"strings"
	"time"
)

/*
Formatter interface defines the methods required for formatting log records.
*/
type Formatter interface {
	Format(original *Record, current *stringRecord) (string, error)
	DateFormat() string
	Next() types.Optional[Formatter]
}
type FormatterBase struct {
	next *Formatter
}

func (f *FormatterBase) Next() types.Optional[Formatter] {
	return types.OptionalOf(*f.next, f.next != nil)
}

type FormatterBuilder interface {
	Prev(Formatter) FormatterBuilder
	New() Formatter
}

/* SimpleFormatter implements a basic text formatter for log records.
 * It formats the log record into a string with a specific structure.
 * The default date format is "02/01/2006 15:04:05 -07:00".
 * You can customize the date format using the TimeFormat method.
 */
type simpleFormatter struct {
	FormatterBase
	dateFormat string
}

func (f *simpleFormatter) Format(original *Record, current *stringRecord) (string, error) {
	formatTime := original.Time.Format(f.DateFormat())
	formatFile := fmt.Sprint(original.File) + ":" + fmt.Sprint(original.Line)

	formatCtx := "{"
	for k, v := range original.Context {
		formatCtx += fmt.Sprintf("%s: %v, ", k, v)
	}
	formatCtx = strings.TrimSuffix(formatCtx, ", ")
	formatCtx += "}"

	res := strings.TrimSpace(fmt.Sprintf("|%s| %s %s: %s %s",
		original.LogLevel.Name(),
		formatTime,
		formatFile,
		original.Message,
		formatCtx,
	))

	if current == nil {
		if f.Next().IsPresent() {
			nextFormatter := f.Next().Get()
			return nextFormatter.Format(original, nil)
		} else {
			current = &stringRecord{
				LogLevel:   original.LogLevel.Name(),
				Time:       formatTime,
				Message:    original.Message,
				File:       original.File,
				Line:       fmt.Sprint(original.Line),
				AppVersion: original.AppVersion,
				Context:    original.Context,
			}
		}

	}

	return res, nil
}
func (f *simpleFormatter) FormatStrings(record *stringRecord) (string, error) {
	formatCtx := "{"
	for k, v := range record.Context {
		formatCtx += fmt.Sprintf("%s: %v, ", k, v)
	}
	formatCtx = strings.TrimSuffix(formatCtx, ", ")
	formatCtx += "}"

	formatTime, err := time.Parse(f.DateFormat(), record.Time)

	return fmt.Sprintf("|%s| %s %s:%s: %s %s",
		record.LogLevel,
		formatTime,
		record.File,
		record.Line,
		record.Message,
		formatCtx,
	), err
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
func (b simpleFormatterBuilder) Prev(formatter Formatter) FormatterBuilder {
	b.formatter.FormatterBase.next = &formatter
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
	FormatterBase
	formatter   Formatter
	colorScheme AnsiColorScheme
}

func (f *consoleColorFormatter) Format(original *Record, record *stringRecord) (string, error) {
	levelString := original.LogLevel.Name()
	timeString := original.Time.Format(f.DateFormat())
	lineString := fmt.Sprint(original.Line)

	styleLevel := f.colorScheme.GetStyle(original.LogLevel.CodeName())
	styleTime := f.colorScheme.GetStyle("time")
	styleMessage := f.colorScheme.GetStyle("message")
	styleFile := f.colorScheme.GetStyle("file")
	styleLine := f.colorScheme.GetStyle("line")
	styleVersion := f.colorScheme.GetStyle("version")
	StyleCtxKey := f.colorScheme.GetStyle("context-key")
	StyleCtxValue := f.colorScheme.GetStyle("context-value")

	StyledCtx := utils.DFlatten(original.Context, func(k, v string) string {
		return fmt.Sprintf("%s: %s", StyleCtxKey.Apply(k), StyleCtxValue.Apply(v))
	})

	utils.Reduce(StyledCtx, func(acc, ctx string) string { return acc + ", " + ctx }, "")

	recordString := stringRecord{
		LogLevel:   styleLevel.Apply(levelString),
		Time:       styleTime.Apply(timeString),
		Message:    styleMessage.Apply(original.Message),
		File:       styleFile.Apply(original.File),
		Line:       styleLine.Apply(lineString),
		AppVersion: styleVersion.Apply(original.AppVersion),
		Context:    original.Context,
	}

	return recordString.Context["formatted"], nil
}
func (f *consoleColorFormatter) FormatStrings(record *stringRecord) (string, error) {
	styleLevel := f.colorScheme.GetStyle(record.LogLevel)
	styleTime := f.colorScheme.GetStyle("time")
	styleMessage := f.colorScheme.GetStyle("message")
	styleFile := f.colorScheme.GetStyle("file")
	styleLine := f.colorScheme.GetStyle("line")
	styleVersion := f.colorScheme.GetStyle("version")
	StyleCtxKey := f.colorScheme.GetStyle("context-key")
	StyleCtxValue := f.colorScheme.GetStyle("context-value")

	record.LogLevel = styleLevel.Apply(record.LogLevel)
	record.Time = styleTime.Apply(record.Time)
	record.Message = styleMessage.Apply(record.Message)
	record.File = styleFile.Apply(record.File)
	record.Line = styleLine.Apply(record.Line)
	record.AppVersion = styleVersion.Apply(record.AppVersion)
	record.Context = utils.DMap(record.Context, func(k, v string) (string, string) {
		return StyleCtxKey.Apply(k), StyleCtxValue.Apply(v)
	})

	return "formatted", nil

}
func (f *consoleColorFormatter) DateFormat() string {
	return f.formatter.DateFormat()
}

type consoleColorFormatterBuilder struct {
	formatter consoleColorFormatter
}

var ConsoleColorFormatter consoleColorFormatterBuilder = consoleColorFormatterBuilder{
	formatter: consoleColorFormatter{
		formatter: SimpleFormatter.New(),
		colorScheme: AnsiColorScheme{styles: map[string]AnsiStyle{
			Level.Debug().CodeName():   CLR_DEBUG,
			Level.Info().CodeName():    CLR_INFO,
			Level.Warning().CodeName(): CLR_WARN,
			Level.Error().CodeName():   CLR_ERROR,
			Level.Fatal().CodeName():   CLR_FATAL,

			Level.Deprecate().CodeName():             CLR_DEPRECATE,
			Level.DeprecateWarning().CodeName():      CLR_DEPRECATE_WARNING,
			Level.DeprecateError().CodeName():        CLR_DEPRECATE_ERROR,
			Level.DeprecateFatal().CodeName():        CLR_DEPRECATE_FATAL,
			Level.Deprecate().CodeName() + "_reason": CLR_DEPR_REASON,

			Level.Lava().CodeName():     CLR_LAVA,
			Level.LavaCold().CodeName(): CLR_COLD_LAVA,
			Level.LavaDry().CodeName():  CLR_DRIED_LAVA,

			"file": CLR_FILE,

			"default": CLR_DEFAULT,
		}},
	},
}

func (b consoleColorFormatterBuilder) Formatter(formatter Formatter) consoleColorFormatterBuilder {
	b.formatter.formatter = formatter
	return b
}
func (b consoleColorFormatterBuilder) AddColor(level logLevel, colorCode AnsiStyle) consoleColorFormatterBuilder {
	//if b.formatter.colors == nil {
	//	b.formatter.colors = make(map[logLevel]string)
	//}
	//b.formatter.colors[level] = colorCode
	return b
}
func (b consoleColorFormatterBuilder) DefaultColor(colorCode string) consoleColorFormatterBuilder {
	//if b.formatter.colors == nil {
	//	b.formatter.colors = make(map[logLevel]string)
	//}
	//
	//b.formatter.colors[LEVEL_NONE] = colorCode
	return b
}
func (b consoleColorFormatterBuilder) Prev(formatter Formatter) FormatterBuilder {
	b.formatter.FormatterBase.next = &formatter
	return b
}
func (b consoleColorFormatterBuilder) New() Formatter {
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
var _ FormatterBuilder = (*consoleColorFormatterBuilder)(nil)
