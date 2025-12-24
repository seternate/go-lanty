package logging

import (
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	//Log level [disable, trace, debug, info, warning, error, panic, fatal]
	LogLevel string

	//Determines if log files are written
	FileLoggingEnabled bool

	//Directory to store the log files
	Directory string

	//File name of the log file
	Filename string

	//Maximum file size (megabytes) of log file before it will be rotated if logging is enabled
	MaxSize int

	//Maximum number of old log files to retain
	MaxBackups int

	//Maximum number of days to retain old log files
	MaxAge int
}

// Configures a zerolog.Logger with the given configuration.
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
	writers = append(writers, os.Stderr)

	mw := io.MultiWriter(writers...)

	return zerolog.New(mw).With().Timestamp().Logger()
}
