package logger

type Options struct {
	Formatter Formatter

	showTime       bool
	showFile       bool
	showLine       bool
	showAppVersion bool

	writers map[string]Writer
}
