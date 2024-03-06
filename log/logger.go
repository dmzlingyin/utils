package log

import (
	"io"
	"log"
	"os"
)

const (
	LevelDebug int = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger interface {
	Debug(v ...any)
	Debugf(format string, v ...any)
	Info(v ...any)
	Infof(format string, v ...any)
	Warn(v ...any)
	Warnf(format string, v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
	SetLevel(lvl int)
}

// default logger
var logger = New(LevelDebug)

func New(level int, writers ...io.Writer) Logger {
	var w io.Writer = os.Stdout
	if len(writers) > 0 {
		if v := writers[0]; v != nil {
			w = v
		}
	}

	return &defaultLogger{
		debugLogger: log.New(w, "[Debug]: ", log.Ldate|log.Ltime),
		infoLogger:  log.New(w, "[Info]: ", log.Ldate|log.Ltime),
		warnLogger:  log.New(w, "[Warn]: ", log.Ldate|log.Ltime),
		errorLogger: log.New(w, "[Error]: ", log.Ldate|log.Ltime),
		level:       level,
	}
}

type defaultLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	level       int
}

func (l *defaultLogger) Debug(v ...any) {
	if l.level > LevelDebug {
		return
	}
	l.debugLogger.Println(v...)
}

func (l *defaultLogger) Debugf(format string, v ...any) {
	if l.level > LevelDebug {
		return
	}
	l.debugLogger.Printf(format, v...)
}

func (l *defaultLogger) Info(v ...any) {
	if l.level > LevelInfo {
		return
	}
	l.infoLogger.Println(v...)
}

func (l *defaultLogger) Infof(format string, v ...any) {
	if l.level > LevelInfo {
		return
	}
	l.infoLogger.Printf(format, v...)
}

func (l *defaultLogger) Warn(v ...any) {
	if l.level > LevelWarn {
		return
	}
	l.warnLogger.Println(v...)
}

func (l *defaultLogger) Warnf(format string, v ...any) {
	if l.level > LevelWarn {
		return
	}
	l.warnLogger.Printf(format, v...)
}

func (l *defaultLogger) Error(v ...any) {
	if l.level > LevelError {
		return
	}
	l.errorLogger.Println(v...)
}

func (l *defaultLogger) Errorf(format string, v ...any) {
	if l.level > LevelError {
		return
	}
	l.errorLogger.Printf(format, v...)
}

func (l *defaultLogger) SetLevel(lvl int) {
	l.level = lvl
}
