package logging

import "testing"

func TestLogConfiguration(t *testing.T) {
	Configure(Config{LogLevel: "disable"})
	Configure(Config{LogLevel: "trace"})
	Configure(Config{LogLevel: "debug"})
	Configure(Config{LogLevel: "info"})
	Configure(Config{LogLevel: "warning"})
	Configure(Config{LogLevel: "error"})
	Configure(Config{LogLevel: "panic"})
	Configure(Config{LogLevel: "fatal"})
	Configure(Config{})
	Configure(Config{LogLevel: "random"})

	Configure(Config{FileLoggingEnabled: true})
}
