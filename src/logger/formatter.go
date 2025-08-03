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
	Format(original *Record, current *formatRecord) (string, error)
	Next() types.Optional[Formatter]
	DateFormat() string
}
type FormatterBase struct {
	next       *Formatter
	dateFormat string
}

func (f *FormatterBase) DateFormat() string {
	return f.dateFormat
}
func (f *FormatterBase) Next() types.Optional[Formatter] {
	return types.OptionalOf(*f.next, f.next != nil)
}
func (f *FormatterBase) FinalString(original *Record, formatRecord *formatRecord) string {
	formattedLevel := fmt.Sprintf(formatRecord.LogLevel, original.LogLevel.Name())
	formattedTime := fmt.Sprintf(formatRecord.Time, original.Time.Format(f.DateFormat()))
	formattedFile := fmt.Sprintf(formatRecord.File, original.File)
	formattedLine := fmt.Sprintf(formatRecord.Line, original.Line)
	formattedMessage := fmt.Sprintf(formatRecord.Message, original.Message)
	formattedVersion := fmt.Sprintf(formatRecord.AppVersion, original.AppVersion)

	formattedContext := formatRecord.ContextBegin
	for k, pair := range utils.DZipWithDefault(original.Context, formatRecord.Context, "%s: %s, ") {
		v, formatStr := pair.First, pair.Second
		formattedContext += fmt.Sprintf(formatStr, k, v)
	}
	formattedContext = strings.Trim(formattedContext, ", ")
	formattedContext += formatRecord.ContextEnd

	finalString := fmt.Sprint(
		formattedLevel,
		formattedTime,
		formattedFile,
		formattedLine,
		formattedMessage,
		formattedVersion,
		formattedContext,
	)

	return finalString
}

func (f *FormatterBase) DefaultFormat() string {
	return "%s"
}
func (f *FormatterBase) DefaultContextFormat() map[string]string {
	return nil
}
func (f *FormatterBase) DefaultContextPrefix() string {
	return ""
}
func (f *FormatterBase) DefaultContextPostfix() string {
	return ""
}

func (f *FormatterBase) levelFormatString(_ logLevel) string {
	return f.DefaultFormat()
}
func (f *FormatterBase) dateFormatString(_ time.Time) string {
	return f.DefaultFormat()
}
func (f *FormatterBase) fileFormatString(_ string) string {
	return f.DefaultFormat()
}
func (f *FormatterBase) lineFormatString(_ int) string {
	return "%d"
}
func (f *FormatterBase) messageFormatString(_ string) string {
	return f.DefaultFormat()
}
func (f *FormatterBase) versionFormatString(_ string) string {
	return f.DefaultFormat()
}
func (f *FormatterBase) contextFormatString(_ map[string]string) map[string]string {
	return f.DefaultContextFormat()
}
func (f *FormatterBase) contextPrefixString(_ map[string]string) string {
	return f.DefaultContextPrefix()
}
func (f *FormatterBase) contextPostfixString(_ map[string]string) string {
	return f.DefaultContextPostfix()
}

func (f *FormatterBase) Format(original *Record, currentFormat *formatRecord) (string, error) {
	if currentFormat == nil {
		currentFormat = &formatRecord{
			LogLevel:     f.levelFormatString(original.LogLevel),
			Time:         f.dateFormatString(original.Time),
			File:         f.fileFormatString(original.File),
			Line:         f.lineFormatString(original.Line),
			Message:      f.messageFormatString(original.Message),
			AppVersion:   f.versionFormatString(original.AppVersion),
			Context:      f.contextFormatString(original.Context),
			ContextBegin: f.contextPrefixString(original.Context),
			ContextEnd:   f.contextPostfixString(original.Context),
		}
	}

	err := error(nil)
	if f.Next().IsPresent() {
		_, err = f.Next().Get().Format(original, currentFormat)
	}

	return f.FinalString(original, currentFormat), err
}

type FormatterBuilder interface {
	Next(Formatter) FormatterBuilder
	New() Formatter
}

/* SimpleFormatter implements a basic text formatter for log records.
 * It formats the log record into a string with a specific structure.
 * The default date format is "02/01/2006 15:04:05 -07:00".
 * You can customize the date format using the TimeFormat method.
 */
type simpleFormatter struct {
	FormatterBase
}

func (f *simpleFormatter) levelFormatString(_ logLevel) string {
	return "| %s | "
}
func (f *simpleFormatter) dateFormatString(_ time.Time) string {
	return "%s "
}
func (f *simpleFormatter) fileFormatString(_ string) string {
	return "%s:"
}
func (f *simpleFormatter) lineFormatString(_ int) string {
	return "%d: "
}
func (f *simpleFormatter) messageFormatString(_ string) string {
	return "%s "
}
func (f *simpleFormatter) versionFormatString(_ string) string {
	return "[%s] "
}
func (f *simpleFormatter) contextPrefixString(_ map[string]string) string {
	return "{"
}
func (f *simpleFormatter) contextPostfixString(_ map[string]string) string {
	return "}"
}

func (f *simpleFormatter) Format(original *Record, currentFormat *formatRecord) (string, error) {
	if currentFormat == nil {
		currentFormat = &formatRecord{
			LogLevel:     f.levelFormatString(original.LogLevel),
			Time:         f.dateFormatString(original.Time),
			File:         f.fileFormatString(original.File),
			Line:         f.lineFormatString(original.Line),
			Message:      f.messageFormatString(original.Message),
			AppVersion:   f.versionFormatString(original.AppVersion),
			Context:      f.contextFormatString(original.Context),
			ContextBegin: f.contextPrefixString(original.Context),
			ContextEnd:   f.contextPostfixString(original.Context),
		}
	}

	err := error(nil)
	if f.Next().IsPresent() {
		_, err = f.Next().Get().Format(original, currentFormat)
	}

	return f.FinalString(original, currentFormat), err
}

type simpleFormatterBuilder struct {
	formatter simpleFormatter
}

var SimpleFormatter simpleFormatterBuilder = simpleFormatterBuilder{
	formatter: simpleFormatter{
		FormatterBase: FormatterBase{
			dateFormat: "",
			next:       nil,
		},
	},
}

func (b simpleFormatterBuilder) TimeFormat(format string) simpleFormatterBuilder {
	b.formatter.dateFormat = format
	return b
}
func (b simpleFormatterBuilder) Next(formatter Formatter) FormatterBuilder {
	b.formatter.next = &formatter
	return b
}
func (b simpleFormatterBuilder) New() Formatter {
	if b.formatter.dateFormat == "" {
		b.formatter.dateFormat = "02/01/2006 15:04:05 -07:00"
	}

	t := &simpleFormatter{
		FormatterBase: b.formatter.FormatterBase,
	}

	return t
}

type consoleColorFormatter struct {
	FormatterBase
	baseFormatter Formatter
	colorScheme   AnsiColorScheme
}

func (f *consoleColorFormatter) Format(original *Record, record *formatRecord) (string, error) {
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

	recordString := formatRecord{
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
func (f *consoleColorFormatter) FormatStrings(record *formatRecord) (string, error) {
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

type consoleColorFormatterBuilder struct {
	formatter consoleColorFormatter
}

var ConsoleColorFormatter consoleColorFormatterBuilder = consoleColorFormatterBuilder{
	formatter: consoleColorFormatter{
		baseFormatter: SimpleFormatter.New(),
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
	b.formatter.baseFormatter = formatter
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
func (b consoleColorFormatterBuilder) Next(formatter Formatter) FormatterBuilder {
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
		baseFormatter: b.formatter.baseFormatter,
		colorScheme:   b.formatter.colorScheme,
	}
}

var _ Formatter = (*simpleFormatter)(nil)
var _ Formatter = (*consoleColorFormatter)(nil)
var _ FormatterBuilder = (*simpleFormatterBuilder)(nil)
var _ FormatterBuilder = (*consoleColorFormatterBuilder)(nil)
