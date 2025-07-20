package logger

import (
	"fmt"
	"os"
)

type Writer interface {
	Write(text string) error
}

type ConsoleWriter struct {
	NewLineTerminated bool
}

func (w *ConsoleWriter) Write(text string) error {
	if w.NewLineTerminated {
		text += "\n"
	}
	_, err := fmt.Print(text)
	return err
}

type FileWriter struct {
	NewLineTerminated bool
	FilePath          string
}

func (w *FileWriter) Write(text string) error {
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
