package logger

import (
	"fmt"
	"os"
)

type Writer interface {
	Write(text string) error
}
type WriterBuilder interface {
	New() Writer
}

type consoleWriter struct {
	NewLineTerminated bool
}

func (w *consoleWriter) Write(text string) error {
	if w.NewLineTerminated {
		text += "\n"
	}
	_, err := fmt.Print(text)
	return err
}

type ConsoleWriterBuilder struct {
	writer consoleWriter
}

func (b ConsoleWriterBuilder) NewLine() ConsoleWriterBuilder {
	b.writer.NewLineTerminated = true
	return b
}
func (b ConsoleWriterBuilder) New() Writer {
	return &consoleWriter{
		NewLineTerminated: b.writer.NewLineTerminated,
	}
}

type fileWriter struct {
	NewLineTerminated bool
	FilePath          string
}

func (w *fileWriter) Write(text string) error {
	if w.NewLineTerminated {
		text += "\n"
	}

	file, err := os.OpenFile(w.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	return err
}

type FileWriterBuilder struct {
	writer fileWriter
}

func (b FileWriterBuilder) FilePath(path string) FileWriterBuilder {
	b.writer.FilePath = path
	return b
}
func (b FileWriterBuilder) NewLine() FileWriterBuilder {
	b.writer.NewLineTerminated = true
	return b
}
func (b FileWriterBuilder) New() Writer {
	if b.writer.FilePath == "" {
		b.writer.FilePath = "logs.log"
	}

	return &fileWriter{
		NewLineTerminated: b.writer.NewLineTerminated,
		FilePath:          b.writer.FilePath,
	}
}

var _ Writer = (*consoleWriter)(nil)
var _ Writer = (*fileWriter)(nil)

var ConsoleWriter ConsoleWriterBuilder
var FileWriter FileWriterBuilder
