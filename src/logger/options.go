package logger

const (
	DEFAULT_LOG_LEVELS       = LEVEL_ALL
	DEFAULT_SHOW_TIME        = true
	DEFAULT_SHOW_FILE        = false
	DEFAULT_SHOW_LINE        = false
	DEFAULT_SHOW_APP_VERSION = false
	DEFAULT_COLORIZE         = false
)

var (
	DEFAULT_FORMATTER = SimpleFormatter.New()
	DEFAULT_WRITERS   = []Writer{ConsoleWriter.NewLine().New()}
)

type options struct {
	levels logLevel

	showTime       bool
	showFile       bool
	showLine       bool
	showAppVersion bool
	colorize       bool

	formatter Formatter
	writers   []Writer
}

type optionBuilder struct {
	opts options
}

var Opt optionBuilder = optionBuilder{
	opts: options{
		levels:         DEFAULT_LOG_LEVELS,
		showTime:       DEFAULT_SHOW_TIME,
		showFile:       DEFAULT_SHOW_FILE,
		showLine:       DEFAULT_SHOW_LINE,
		showAppVersion: DEFAULT_SHOW_APP_VERSION,
		colorize:       DEFAULT_COLORIZE,
		formatter:      DEFAULT_FORMATTER,
		writers:        DEFAULT_WRITERS,
	},
}

func (o *optionBuilder) SetLevels(levels logLevel) *optionBuilder {
	o.opts.levels = levels
	return o
}
func (o *optionBuilder) AddLevel(level logLevel) *optionBuilder {
	o.opts.levels |= level
	return o
}
func (o *optionBuilder) RemoveLevel(level logLevel) *optionBuilder {
	o.opts.levels &= ^level
	return o
}

func (o *optionBuilder) ShowDate() *optionBuilder {
	o.opts.showTime = true
	return o
}
func (o *optionBuilder) ShowFile() *optionBuilder {
	o.opts.showFile = true
	return o
}
func (o *optionBuilder) SetShowLine(show bool) *optionBuilder {
	o.opts.showLine = show
	return o
}
func (o *optionBuilder) SetShowAppVersion(show bool) *optionBuilder {
	o.opts.showAppVersion = show
	return o
}

func (o *optionBuilder) SetColorize(colorize bool) *optionBuilder {
	o.opts.colorize = colorize
	return o
}

func (o *optionBuilder) SetFormatter(formatter Formatter) *optionBuilder {
	o.opts.formatter = formatter
	return o
}
func (o *optionBuilder) AddWriter(writer Writer) *optionBuilder {
	o.opts.writers = append(o.opts.writers, writer)
	return o
}

func (o *optionBuilder) Build() options {
	if o.opts.formatter == nil {
		o.opts.formatter = SimpleFormatter.New()
	}

	if len(o.opts.writers) == 0 {
		o.opts.writers = []Writer{
			ConsoleWriter.NewLine().New(),
			FileWriter.NewLine().New(),
		}
	}

	returnOptions := o.opts
	o.opts = options{
		levels:         DEFAULT_LOG_LEVELS,
		showTime:       DEFAULT_SHOW_TIME,
		showFile:       DEFAULT_SHOW_FILE,
		showLine:       DEFAULT_SHOW_LINE,
		showAppVersion: DEFAULT_SHOW_APP_VERSION,
		colorize:       DEFAULT_COLORIZE,
		formatter:      DEFAULT_FORMATTER,
		writers:        DEFAULT_WRITERS,
	}

	return returnOptions
}
