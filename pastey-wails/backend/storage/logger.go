package storage

import (
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	path     string
	fileName string
}

func NewLogger(path string) *Logger {
	return &Logger{
		path:     path,
		fileName: "pastey.log",
	}
}

func (l *Logger) Log(message string) {
	file, err := os.OpenFile(filepath.Join(l.path, l.fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString(time.Now().String() + ": " + message + "\n"); err != nil {
		panic(err)
	}
}
