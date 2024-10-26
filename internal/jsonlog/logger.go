package jsonlog

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case 1:
		return "ERROR"
	case 2:
		return "FATAL"
	default:
		return "INFO"
	}
}

type Logger struct {
	out       io.Writer
	errorFile *os.File
	mu        sync.Mutex
}

func New(out io.Writer, errorFilePath string) (*Logger, error) {
	errorFile, err := os.OpenFile(errorFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{
		errorFile: errorFile,
		out:       out,
	}, nil
}

func (l *Logger) PrintInfo(message string, properties map[string]string) {
	err := l.print(LevelInfo, message, properties)
	if err != nil {
		log.Println("jsonlog failed to initialize correctly")
	}

}

func (l *Logger) PrintError(err error, properties map[string]string) {
	logError := l.print(LevelError, err.Error(), properties)
	if logError != nil {
		log.Println("jsonlog failed to initialize correctly")
	}
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	logError := l.print(LevelFatal, err.Error(), properties)
	if logError != nil {
		log.Println("jsonlog failed to initialize correctly")
	}
}

func (l *Logger) print(level Level, message string, properties map[string]string) error {
	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties"`
		Trace      string            `json:"trace"`
	}{
		Level:      level.String(),
		Time:       time.Now().Format(time.RFC1123),
		Message:    message,
		Properties: properties,
	}
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal the log message" + err.Error())
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err = l.out.Write(append(line, '\n'))
	if err != nil {
		return err
	}
	if level >= LevelError && l.errorFile != nil {
		_, err = l.errorFile.Write(append(line, '\n'))
		if err != nil {
			return err
		}
	}
	return nil
}
