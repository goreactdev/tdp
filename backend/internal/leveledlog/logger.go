package leveledlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgHiRed, color.Bold).SprintFunc()
	yellow = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	cyan   = color.New(color.FgHiCyan, color.Bold).SprintFunc()
)

type Level int8

const (
	LevelAll Level = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelWarning:
		return "WARNING"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	useJSON  bool
	colorize bool
	mu       sync.Mutex
}

func NewLogger(out io.Writer, minLevel Level, colorize bool) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
		colorize: colorize,
	}
}

func NewJSONLogger(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
		useJSON:  true,
	}
}

func (l *Logger) Info(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.print(LevelInfo, message, nil)
}

func (l *Logger) Warning(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.print(LevelWarning, message, nil)
}

func (l *Logger) Error(err error, trace []byte) {
	l.print(LevelError, err.Error(), trace)
}

func (l *Logger) Fatal(err error, trace []byte) {
	l.print(LevelFatal, err.Error(), trace)
	os.Exit(1)
}

func (l *Logger) Write(message []byte) (n int, err error) {
	messageStr := string(message)

	if strings.Contains(messageStr, "runtime/panic") {
		return l.print(LevelError, messageStr, nil)
	}

	return l.print(LevelWarning, messageStr, nil)
}

func (l *Logger) print(level Level, message string, trace []byte) (n int, err error) {
	if level < l.minLevel {
		return 0, nil
	}

	var line string

	if l.useJSON {
		line = jsonLine(level, message, trace)
	} else {
		line = textLine(level, message, trace, l.colorize)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return fmt.Fprintln(l.out, line)
}

func textLine(level Level, message string, trace []byte, colorize bool) string {
	line := fmt.Sprintf("level=%q time=%q message=%q", level, time.Now().UTC().Format(time.RFC3339), strings.TrimSpace(message))

	if colorize {
		switch level {
		case LevelError, LevelFatal:
			line = red(line)
		case LevelWarning:
			line = yellow(line)
		case LevelInfo:
			line = cyan(line)
		}
	}

	if trace != nil {
		line += fmt.Sprintf("\n%s", string(trace))
	}

	return line
}

func jsonLine(level Level, message string, trace []byte) string {
	aux := struct {
		Level   string `json:"level"`
		Time    string `json:"time"`
		Message string `json:"message"`
		Trace   string `json:"trace,omitempty"`
	}{
		Level:   level.String(),
		Time:    time.Now().UTC().Format(time.RFC3339),
		Message: message,
	}

	if trace != nil {
		aux.Trace = string(trace)
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		return fmt.Sprintf("%s: unable to marshal log message: %s", LevelError.String(), err.Error())
	}

	return string(line)
}
