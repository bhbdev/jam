package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

var logger slog.Logger

func init() {
	w := os.Stdout
	logger = *slog.New(tint.NewHandler(w, &tint.Options{
		Level:   slog.LevelInfo,
		NoColor: !isatty.IsTerminal(w.Fd()), // Disable color if not a terminal
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if err, ok := a.Value.Any().(error); ok {
				aErr := tint.Err(err)
				aErr.Key = a.Key
				return aErr
			}
			return a
		},
	}))
}

func Logger() *slog.Logger {
	return &logger
}

func SetDebugLevel(debug bool) {
	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

// Info log
func Info(msg string, v ...interface{}) {
	logger.Info(msg, v...)
}

// Debug log
func Debug(msg string, v ...interface{}) {
	logger.Debug(msg, v...)
}

// Warn log
func Warn(msg string, v ...interface{}) {
	logger.Warn(msg, v...)
}

// Error log
func Error(msg string, v ...interface{}) {
	logger.Error(msg, v...)
}

// Fatal log, this will exit the program
func Fatal(msg string, v ...interface{}) {
	logger.Error(msg, v...)
	os.Exit(1)
}
