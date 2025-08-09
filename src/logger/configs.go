package logger

import "dainxor/atv/types"

type writerAndFormatter struct {
	writer    Writer
	formatter Formatter
}

type configurations struct {
	logLevels logLevel
	logFlags  logFlag

	maxLogAttempts     uint8
	warningLogAttempts uint8

	panicOnMaxAttempts      bool
	canPanicOnAbnormalWrite bool

	appVersion types.Version

	writers []writerAndFormatter
}

// NewConfigs initializes a new configs instance with default values
func NewConfigs() configurations {
	return configurations{
		logLevels: Level.All(),
		logFlags:  Flag.DateTime() | Flag.File() | Flag.Line() | Flag.AppVersion(),

		warningLogAttempts: 10,
		maxLogAttempts:     15,

		panicOnMaxAttempts:      true,
		canPanicOnAbnormalWrite: true,

		appVersion: types.V("0.1.0"),

		writers: []writerAndFormatter{
			{writer: ConsoleWriter.NewLine().New(), formatter: SimpleFormatter.New()},
			{writer: FileWriter.NewLine().New(), formatter: SimpleFormatter.New()},
		},
	}
}

func (c *configurations) AddWriter(writer Writer, formatter Formatter) *configurations {
	c.writers = append(c.writers, writerAndFormatter{writer: writer, formatter: formatter})
	return c
}
func (c *configurations) RemoveWriter(index int) *configurations {
	if index < 0 || index >= len(c.writers) {
		return c
	}
	c.writers = append(c.writers[:index], c.writers[index+1:]...)
	return c
}
func (c *configurations) RemoveWriters(index ...int) *configurations {
	indexes := len(index)
	if indexes == 0 {
		return c
	}

	last := len(c.writers) - 1

	for _, idx := range index {
		c.writers[idx].writer.Close()

		swap := c.writers[last]
		c.writers[last] = c.writers[idx]
		c.writers[idx] = swap
		last--
	}

	c.writers = c.writers[:last]
	return c
}
func (c *configurations) Writers() []writerAndFormatter {
	return c.writers
}
func (c *configurations) Writer(index int) (*Writer, *Formatter) {
	if index < 0 || index >= len(c.writers) {
		return nil, nil
	}
	return &c.writers[index].writer, &c.writers[index].formatter
}
