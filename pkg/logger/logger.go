package logger

import (
	"os"
	"path/filepath"
	"sync"

	"errors"

	"github.com/rs/zerolog"
)

const (
	logDir    = "logs"
	infoFile  = "info.json"
	debugFile = "debug.json"
	warnFile  = "warn.json"
	errorFile = "error.json"
)

// Logger - обертка над zerolog
// Каждый уровень пишет в свой файл
// Все логи в формате JSON
// Thread-safe

type Logger struct {
	InfoLogger  zerolog.Logger
	DebugLogger zerolog.Logger
	WarnLogger  zerolog.Logger
	ErrorLogger zerolog.Logger
	mu          sync.Mutex
}

func getProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return "", errors.New("project root (go.mod) not found")
		}
		wd = parent
	}
}

func NewLogger() (*Logger, error) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		return nil, err
	}
	logDir := filepath.Join(projectRoot, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	infoF, err := os.OpenFile(filepath.Join(logDir, infoFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	debugF, err := os.OpenFile(filepath.Join(logDir, debugFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	warnF, err := os.OpenFile(filepath.Join(logDir, warnFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	errorF, err := os.OpenFile(filepath.Join(logDir, errorFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		InfoLogger:  zerolog.New(infoF).With().Timestamp().Logger().Level(zerolog.InfoLevel),
		DebugLogger: zerolog.New(debugF).With().Timestamp().Logger().Level(zerolog.DebugLevel),
		WarnLogger:  zerolog.New(warnF).With().Timestamp().Logger().Level(zerolog.WarnLevel),
		ErrorLogger: zerolog.New(errorF).With().Timestamp().Logger().Level(zerolog.ErrorLevel),
	}, nil
}

func (l *Logger) Info() *zerolog.Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.InfoLogger.Info()
}

func (l *Logger) Debug() *zerolog.Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.DebugLogger.Debug()
}

func (l *Logger) Warn() *zerolog.Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.WarnLogger.Warn()
}

func (l *Logger) Error() *zerolog.Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.ErrorLogger.Error()
}

func (l *Logger) With() zerolog.Context {
	return l.InfoLogger.With()
}

func (l *Logger) WithContext(ctx interface{}) zerolog.Context {
	return l.InfoLogger.With().Interface("ctx", ctx)
}

func NewTestLogger() *Logger {
	return &Logger{
		InfoLogger:  zerolog.Nop(),
		DebugLogger: zerolog.Nop(),
		WarnLogger:  zerolog.Nop(),
		ErrorLogger: zerolog.Nop(),
	}
}

