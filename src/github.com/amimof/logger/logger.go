package logger

import "fmt"

// Define colors
const CLR_0 = "\x1b[30;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"
const CLR_Y = "\x1b[33;1m"
const CLR_B = "\x1b[34;1m"
const CLR_M = "\x1b[35;1m"
const CLR_C = "\x1b[36;1m"
const CLR_W = "\x1b[37;1m"
const CLR_N = "\x1b[0m"
const DEBUG = 3
const INFO = 2
const WARN = 1
const ERROR = 0

type Logger struct {
	name  string
	Level *Level
}

type Level struct {
  num int
  name string
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

func (l *Logger) GetLevel() *Level {
	return l.Level
}

func (l *Logger) Debug(message ...interface{}) {
  if l.Level.num >= DEBUG {
    fmt.Printf("%s %s %s %s %s\n", CLR_G, l.Level.name, l.name, message, CLR_N)
  }
}

func (l *Logger) Info(message ...interface{}) {
  if l.Level.num >= INFO {
    fmt.Printf("%s %s %s %s %s\n", CLR_W, l.Level.name, l.name, message, CLR_N)
  }
}

func (l *Logger) Warn(message ...interface{}) {
  if l.Level.num >= WARN {
    fmt.Printf("%s %s %s %s %s\n", CLR_Y, l.Level.name, l.name, message, CLR_N)
  }
}

func (l *Logger) Error(message ...interface{}) {
  if l.Level.num >= ERROR {
    fmt.Printf("%s %s %s %s %s\n", CLR_R, l.Level.name, l.name, message, CLR_N)
  }
}

func SetupNew(name string) *Logger {
	new := &Logger{
		name:  name,
		Level: &Level{1, "INFO"},
	}
	return new
}
