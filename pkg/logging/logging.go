package logging

import (
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	LogLevel              string
	ConsoleLoggingEnabled bool
	FileLoggingEnabled    bool
	Directory             string
	Filename              string
	MaxSize               int
	MaxBackups            int
	MaxAge                int
}

func Configure(config Config) zerolog.Logger {
	var writers []io.Writer

	switch config.LogLevel {
	case "disable":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	}

	if config.FileLoggingEnabled {
		filelogger := &lumberjack.Logger{
			Filename:   path.Join(config.Directory, config.Filename),
			MaxBackups: config.MaxBackups,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
		}
		writers = append(writers, filelogger)
	}
	if config.ConsoleLoggingEnabled {
		writers = append(writers, os.Stderr)
	}

	mw := io.MultiWriter(writers...)

	return zerolog.New(mw).With().Timestamp().Logger()
}
