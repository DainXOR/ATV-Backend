package logger

import "fmt"

type ansiCode string

func (c ansiCode) String() string {
	return string(c)
}

const ( // Color constants
	TXT_BLACK   ansiCode = "30m"
	TXT_RED     ansiCode = "31m"
	TXT_GREEN   ansiCode = "32m"
	TXT_YELLOW  ansiCode = "33m"
	TXT_BLUE    ansiCode = "34m"
	TXT_MAGENTA ansiCode = "35m"
	TXT_CYAN    ansiCode = "36m"
	TXT_WHITE   ansiCode = "37m"

	BG_BLACK   ansiCode = "40"
	BG_RED     ansiCode = "41"
	BG_GREEN   ansiCode = "42"
	BG_YELLOW  ansiCode = "43"
	BG_BLUE    ansiCode = "44"
	BG_MAGENTA ansiCode = "45"
	BG_CYAN    ansiCode = "46"
	BG_WHITE   ansiCode = "47"

	CLR_START ansiCode = "\033["
	CLR_RESET ansiCode = "\033[0m"
)

var (
	CLR_DEBUG = AnsiStyle{Background: BG_GREEN, Text: TXT_BLACK}
	CLR_INFO  = AnsiStyle{Background: BG_CYAN, Text: TXT_BLACK}
	CLR_WARN  = AnsiStyle{Background: BG_YELLOW, Text: TXT_BLACK}
	CLR_ERROR = AnsiStyle{Background: BG_RED, Text: TXT_BLACK}
	CLR_FATAL = AnsiStyle{Background: BG_RED, Text: TXT_WHITE}

	CLR_DEPRECATE   = AnsiStyle{Background: BG_MAGENTA, Text: TXT_WHITE}
	CLR_DEPR_REASON = AnsiStyle{Background: BG_YELLOW, Text: TXT_WHITE}

	CLR_LAVA       = AnsiStyle{Background: BG_WHITE, Text: TXT_BLACK}
	CLR_COLD_LAVA  = AnsiStyle{Background: BG_YELLOW, Text: TXT_BLACK}
	CLR_DRIED_LAVA = AnsiStyle{Background: BG_RED, Text: TXT_BLACK}

	CLR_FILE = AnsiStyle{Background: BG_BLUE, Text: TXT_WHITE}
)

type ColorScheme[S Style] interface {
	GetStyle(name string) S
}
type Style interface {
	Apply(text string) string
}

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
	Background ansiCode
	Text       ansiCode
}

func (s AnsiStyle) Apply(text string) string {
	return fmt.Sprintf("%s%s%s", CLR_START+s.Background+";"+s.Text, text, CLR_RESET)
}

var _ ColorScheme[AnsiStyle] = (*AnsiColorScheme)(nil)
var _ Style = (*AnsiStyle)(nil)
