package logger

import (
	"os"

	"github.com/aiagent/boilerplate/internal/infrastructure/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger is the application logger
type Logger struct {
	*zerolog.Logger
}

// New creates a new logger instance
func New(cfg *config.LoggerConfig) *Logger {
	var l zerolog.Logger

	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set output format
	if cfg.Format == "text" {
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().
			Timestamp().
			Caller().
			Logger()
	} else {
		l = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger()
	}

	// Set as global logger
	log.Logger = l

	return &Logger{Logger: &l}
}

// Info logs an info message
func Info(msg string, fields ...map[string]interface{}) {
	event := log.Info()
	for _, f := range fields {
		for k, v := range f {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Error logs an error message
func Error(msg string, err error, fields ...map[string]interface{}) {
	event := log.Error().Err(err)
	for _, f := range fields {
		for k, v := range f {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Debug logs a debug message
func Debug(msg string, fields ...map[string]interface{}) {
	event := log.Debug()
	for _, f := range fields {
		for k, v := range f {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Warn logs a warning message
func Warn(msg string, fields ...map[string]interface{}) {
	event := log.Warn()
	for _, f := range fields {
		for k, v := range f {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, err error, fields ...map[string]interface{}) {
	event := log.Fatal().Err(err)
	for _, f := range fields {
		for k, v := range f {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}
