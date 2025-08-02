package logger

type wf struct {
	writer    Writer
	formatter *Formatter
}

type configs struct {
	LogLevel logLevel
	LogFlags logFlag

	MaxLogAttempts     uint8
	WarningLogAttempts uint8

	AppVersion string

	wf []wf
}

// NewConfigs initializes a new configs instance with default values
func NewConfigs() configs {
	fFormatter := SimpleFormatter.New()
	cFormatter := ConsoleColorFormatter.New()

	return configs{
		LogLevel: Level.All(),
		LogFlags: Flag.DateTime() | Flag.File() | Flag.Line() | Flag.AppVersion(),

		WarningLogAttempts: 10,
		MaxLogAttempts:     15,

		AppVersion: "0.1.0",

		wf: []wf{
			{writer: ConsoleWriter.NewLine().New(), formatter: &cFormatter},
			{writer: FileWriter.NewLine().New(), formatter: &fFormatter},
		},
	}
}

func (c *configs) AddWriter(writer Writer, formatter *Formatter) *configs {
	c.wf = append(c.wf, wf{writer: writer, formatter: formatter})
	return c
}
func (c *configs) Writers() []wf {
	return c.wf
}
func (c *configs) Writer(index int) (*Writer, *Formatter) {
	if index < 0 || index >= len(c.wf) {
		return nil, nil
	}
	return &c.wf[index].writer, c.wf[index].formatter
}
