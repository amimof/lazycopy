package logger

import (
	"fmt"
	"io"
	"sync"
	"os"
)

// Define colors
const (
	CLR_0 = "\x1b[30;1m"
	CLR_R = "\x1b[31;1m"
	CLR_G = "\x1b[32;1m"
	CLR_Y = "\x1b[33;1m"
	CLR_B = "\x1b[34;1m"
	CLR_M = "\x1b[35;1m"
	CLR_C = "\x1b[36;1m"
	CLR_W = "\x1b[37;1m"
	CLR_N = "\x1b[0m"
	DEBUG = 3
	INFO = 2
	WARN = 1
	ERROR = 0
)

var (
	std = NewStdLogger(os.Stderr)
	new *Logger
)

type StandardLogger struct {
	mu sync.Mutex
	out io.Writer
	buf []byte
}

type Logger struct {
	name  string
	Level *Level
}

type Level struct {
  num int
  name string
}

func NewStdLogger(out io.Writer) *StandardLogger {
	return &StandardLogger{out: out}
}

func (l *Logger) New(name string, level *Level) *Logger {
	l.name = name
	l.Level = level
	return l
}

func (l *Level) SetLevel(level int) *Level {
  switch level {
    case 0:
      l.num = 0
      l.name = "ERROR"
    case 1:
      l.num = 1
      l.name = "WARN"
    case 2:
      l.num = 2
      l.name = "INFO"
    case 3:
      l.num = 3
      l.name = "DEBUG"
    default:
      l.num = 1
      l.name = "INFO"
  }
  return l
}

// Writes the specified string to std and colors it accordingly
func (l *StandardLogger) Output(color, str string) error {
	l.mu.Lock()
	lgr := getLogger()
	defer l.mu.Unlock()
	l.buf = l.buf[:0]
	l.buf = append(l.buf, color...)
	l.buf = append(l.buf, lgr.Level.name...)
	l.buf = append(l.buf, ' ')
	l.buf = append(l.buf, lgr.name...)
	l.buf = append(l.buf, ' ')
	l.buf = append(l.buf, str...)
	l.buf = append(l.buf, CLR_N...)
	_, err := l.out.Write(l.buf)
	return err
}
// Return the level
func (l *Logger) GetLevel() *Level {
	return l.Level
}

// Prints according to fmt.Sprintf format specifier and returns the resulting string
func (l *Logger) Debugf(format string, message ...interface{}) {
  if l.Level.num >= DEBUG {
		std.Output(CLR_G, fmt.Sprintf(format, message...))
  }
}

// Prints a debug message on a new line
func (l *Logger) Debug(message ...interface{}) {
  if l.Level.num >= DEBUG {
		std.Output(CLR_G, fmt.Sprintln(message...))
  }
}

// Prints according to fmt.Sprintf format specifier and returns the resulting string
func (l *Logger) Infof(format string, message ...interface{}) {
  if l.Level.num >= INFO {
		std.Output(CLR_W, fmt.Sprintf(format, message...))
  }
}

// Prints n info message on a new line
func (l *Logger) Info(message ...interface{}) {
  if l.Level.num >= INFO {
		std.Output(CLR_W, fmt.Sprintln(message...))
  }
}

// Prints according to fmt.Sprintf format specifier and returns the resulting string
func (l *Logger) Warnf(format string, message ...interface{}) {
	if l.Level.num >= WARN {
		std.Output(CLR_Y, fmt.Sprintf(format, message...))
	}
}

// Prints a warning message on a new line
func (l *Logger) Warn(message ...interface{}) {
  if l.Level.num >= WARN {
		std.Output(CLR_Y, fmt.Sprintln(message...))
  }
}

// Prints according to fmt.Sprintf format specifier and returns the resulting string
func (l *Logger) Errorf(format string, message ...interface{}) {
	if l.Level.num >= ERROR {
		std.Output(CLR_R, fmt.Sprintf(format, message...))
	}
}

// Prints an error message on a new line
func (l *Logger) Error(message ...interface{}) {
  if l.Level.num >= ERROR {
    std.Output(CLR_R, fmt.Sprintln(message...))
  }
}

// Return logger
func getLogger() *Logger {
	return new
}

// Set the standard logger output
func SetOutput(out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = out
}

// Setup and return logger
func SetupNew(name string) *Logger {
	new = &Logger{
		name:  name,
		Level: &Level{1, "INFO"},
	}
	return new
}
