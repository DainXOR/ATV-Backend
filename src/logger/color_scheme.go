package logger

import "fmt"

type AnsiCode string

func (c AnsiCode) String() string {
	return string(c)
}

const ( // Color constants
	TXT_BLACK   AnsiCode = "30m"
	TXT_RED     AnsiCode = "31m"
	TXT_GREEN   AnsiCode = "32m"
	TXT_YELLOW  AnsiCode = "33m"
	TXT_BLUE    AnsiCode = "34m"
	TXT_MAGENTA AnsiCode = "35m"
	TXT_CYAN    AnsiCode = "36m"
	TXT_WHITE   AnsiCode = "37m"

	BG_BLACK   AnsiCode = "40"
	BG_RED     AnsiCode = "41"
	BG_GREEN   AnsiCode = "42"
	BG_YELLOW  AnsiCode = "43"
	BG_BLUE    AnsiCode = "44"
	BG_MAGENTA AnsiCode = "45"
	BG_CYAN    AnsiCode = "46"
	BG_WHITE   AnsiCode = "47"

	CLR_START AnsiCode = "\033["
	CLR_RESET AnsiCode = "\033[0m"

	CLR_NONE AnsiCode = "" // No format
)

var (
	CLR_DEBUG = AnsiStyle{Background: BG_GREEN, Text: TXT_BLACK}
	CLR_INFO  = AnsiStyle{Background: BG_CYAN, Text: TXT_BLACK}
	CLR_WARN  = AnsiStyle{Background: BG_YELLOW, Text: TXT_BLACK}
	CLR_ERROR = AnsiStyle{Background: BG_RED, Text: TXT_BLACK}
	CLR_FATAL = AnsiStyle{Background: BG_RED, Text: TXT_WHITE}

	CLR_DEPRECATE         = AnsiStyle{Background: BG_WHITE, Text: TXT_MAGENTA}
	CLR_DEPRECATE_WARNING = AnsiStyle{Background: BG_YELLOW, Text: TXT_MAGENTA}
	CLR_DEPRECATE_ERROR   = AnsiStyle{Background: BG_RED, Text: TXT_CYAN}
	CLR_DEPRECATE_FATAL   = AnsiStyle{Background: BG_RED, Text: TXT_WHITE}
	CLR_DEPR_REASON       = AnsiStyle{Background: BG_YELLOW, Text: TXT_WHITE}

	CLR_LAVA       = AnsiStyle{Background: BG_WHITE, Text: TXT_BLACK}
	CLR_COLD_LAVA  = AnsiStyle{Background: BG_YELLOW, Text: TXT_BLACK}
	CLR_DRIED_LAVA = AnsiStyle{Background: BG_RED, Text: TXT_BLACK}

	CLR_FILE = AnsiStyle{Background: BG_BLUE, Text: TXT_WHITE}

	CLR_DEFAULT = AnsiStyle{Background: CLR_NONE, Text: CLR_NONE}
)

type ColorScheme[S Style] interface {
	GetStyle(name string) S
}
type Style interface {
	Apply(text string) string
}

/* AnsiColorScheme implements ColorScheme for ANSI styles
 * It provides a mapping of log levels to ANSI styles for console output.
 * It allows for easy customization of log output colors using identifiers.
 * The default identifiers used are listed here for reference:
 * - debug
 * - info
 * - warning
 * - error
 * - fatal
 * - deprecate
 * - deprecate_warning
 * - deprecate_error
 * - deprecate_fatal
 * - deprecate_reason
 * - lava
 * - lava_hot
 * - lava_cold
 * - lava_dry
 * - time
 * - file
 * - line
 * - version
 * - message
 * - context-key
 * - context-value
 * - default (used when no specific style is found)

 * You may add identifiers as needed for your custom formatters.
 */
type AnsiColorScheme struct {
	styles map[string]AnsiStyle
}

func (cs AnsiColorScheme) GetStyle(name string) AnsiStyle {
	if style, exists := cs.styles[name]; exists {
		return style
	}
	return AnsiStyle{Background: BG_BLACK, Text: TXT_WHITE} // Default style
}

type AnsiStyle struct {
	Background AnsiCode
	Text       AnsiCode
}

func (s AnsiStyle) Apply(text string) string {
	return fmt.Sprintf("%s%s%s", CLR_START+s.Background+";"+s.Text, text, CLR_RESET)
}

var _ ColorScheme[AnsiStyle] = (*AnsiColorScheme)(nil)
var _ Style = (*AnsiStyle)(nil)
