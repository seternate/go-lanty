package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var validConf = Config{
	DBString:          "postgres://lanty:lanty@database:5432/lanty?sslmode=disable",
	Port:              8080,
	LogLevel:          "info",
	LogBackups:        0,
	LogFileSize:       10,
	LogAge:            0,
	EnableFileLogging: false,
}

func TestParseEnv(t *testing.T) {
	t.Setenv("DB", "postgres://lanty:lanty@database:5432/lanty?sslmode=disable")
	t.Setenv("PORT", "9090")
	t.Setenv("LOG_LEVEL", "trace")
	t.Setenv("LOG_FILESIZE", "33")
	t.Setenv("ENABLE_FILE_LOGGING", "true")

	config, err := Parse()

	assert.NoError(t, err)
	assert.Equal(t, "postgres://lanty:lanty@database:5432/lanty?sslmode=disable", config.DBString)
	assert.Equal(t, 9090, config.Port)
	assert.Equal(t, "trace", config.LogLevel)
	assert.Equal(t, 33, config.LogFileSize)
	assert.True(t, config.EnableFileLogging)
}

func TestParseArgs(t *testing.T) {
	config, err := Parse(
		"--db", "postgres://lanty:lanty@database:5432/lanty?sslmode=disable",
		"--port", "9090",
		"--loglevel", "trace",
		"--logfilesize", "33",
		"--enablefilelogging",
		"--help",
	)

	assert.NoError(t, err)
	assert.Equal(t, "postgres://lanty:lanty@database:5432/lanty?sslmode=disable", config.DBString)
	assert.Equal(t, 9090, config.Port)
	assert.Equal(t, "trace", config.LogLevel)
	assert.Equal(t, 33, config.LogFileSize)
	assert.True(t, config.EnableFileLogging)
	assert.True(t, config.PrintHelp)
}

func TestParseOverwriteEnvs(t *testing.T) {
	t.Setenv("DB", "dbenv")
	t.Setenv("PORT", "9090")
	t.Setenv("LOG_LEVEL", "error")
	t.Setenv("LOG_FILESIZE", "44")
	t.Setenv("ENABLE_FILE_LOGGING", "false")

	config, err := Parse(
		"--db", "postgres://lanty:lanty@database:5432/lanty?sslmode=disable",
		"--loglevel", "trace",
		"--logfilesize", "33",
		"--enablefilelogging",
	)

	assert.NoError(t, err)
	assert.Equal(t, "postgres://lanty:lanty@database:5432/lanty?sslmode=disable", config.DBString)
	assert.Equal(t, 9090, config.Port)
	assert.Equal(t, "trace", config.LogLevel)
	assert.Equal(t, 33, config.LogFileSize)
	assert.True(t, config.EnableFileLogging)
}

func TestParseValidationFail(t *testing.T) {
	config, err := Parse("--db", "")

	assert.Error(t, err)
	assert.Empty(t, config.DBString)
}

func TestParseFromFlags(t *testing.T) {
	config := parseFromArgs(
		"--db", "postgres://lanty:lanty@database:5432/lanty?sslmode=disable",
		"--port", "9090",
		"--loglevel", "trace",
		"--logbackup", "5",
		"--logfilesize", "33",
		"--logage", "20",
		"--enablefilelogging",
		"--help",
		"--version",
	)

	assert.Equal(t, "postgres://lanty:lanty@database:5432/lanty?sslmode=disable", config.DBString)
	assert.Equal(t, 9090, config.Port)
	assert.Equal(t, "trace", config.LogLevel)
	assert.Equal(t, 5, config.LogBackups)
	assert.Equal(t, 33, config.LogFileSize)
	assert.Equal(t, 20, config.LogAge)
	assert.Equal(t, true, config.EnableFileLogging)
	assert.Equal(t, true, config.PrintHelp)
	assert.Equal(t, true, config.PrintVersion)
}

func TestParseFromFlagsShorthand(t *testing.T) {
	config := parseFromArgs(
		"-p", "9090",
		"-h",
	)

	assert.Equal(t, 9090, config.Port)
	assert.Equal(t, true, config.PrintHelp)
}

func TestParseFromArgsDefaults(t *testing.T) {
	config := parseFromArgs()

	assert.Empty(t, config.DBString)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, "info", config.LogLevel)
	assert.Empty(t, config.LogBackups)
	assert.Equal(t, 10, config.LogFileSize)
	assert.Empty(t, config.LogAge)
	assert.Empty(t, config.EnableFileLogging)
	assert.Empty(t, config.PrintHelp)
	assert.Empty(t, config.PrintVersion)
}

func TestParseFromEnv(t *testing.T) {
	t.Setenv("DB", "postgres://lanty:lanty@database:5432/lanty?sslmode=disable")
	t.Setenv("PORT", "9090")
	t.Setenv("LOG_LEVEL", "trace")
	t.Setenv("LOG_BACKUPS", "5")
	t.Setenv("LOG_FILESIZE", "33")
	t.Setenv("LOG_AGE", "20")
	t.Setenv("ENABLE_FILE_LOGGING", "true")

	config := parseFromEnvs()

	assert.Equal(t, "postgres://lanty:lanty@database:5432/lanty?sslmode=disable", config.DBString)
	assert.Equal(t, 9090, config.Port)
	assert.Equal(t, "trace", config.LogLevel)
	assert.Equal(t, 5, config.LogBackups)
	assert.Equal(t, 33, config.LogFileSize)
	assert.Equal(t, 20, config.LogAge)
	assert.Equal(t, true, config.EnableFileLogging)
	assert.Empty(t, config.PrintHelp)
	assert.Empty(t, config.PrintVersion)
}

func TestParseFromEnvNoDefaults(t *testing.T) {
	t.Setenv("DB", "")
	t.Setenv("PORT", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("LOG_BACKUPS", "")
	t.Setenv("LOG_FILESIZE", "")
	t.Setenv("LOG_AGE", "")
	t.Setenv("ENABLE_FILE_LOGGING", "")

	config := parseFromEnvs()

	assert.Empty(t, config)
}

func TestMergeOverwriteEmpty(t *testing.T) {
	destination := validConf
	destination.DBString = "postgres://lanty:lanty@database:5432/lanty?sslmode=disable"

	source := validConf
	source.DBString = "testingconn"
	source.LogBackups = 33

	config, err := merge(destination, source)

	assert.Equal(t, source.LogBackups, config.LogBackups)
	assert.Equal(t, destination.DBString, config.DBString)
	assert.NoError(t, err)
}

func TestMergeOverwriteDefaults(t *testing.T) {
	destination := validConf

	source := validConf
	source.Port = 1337
	source.LogLevel = "trace"
	source.LogFileSize = 23

	config, err := merge(destination, source)

	assert.Equal(t, source.Port, config.Port)
	assert.Equal(t, source.LogLevel, config.LogLevel)
	assert.Equal(t, source.LogFileSize, config.LogFileSize)
	assert.NoError(t, err)
}

func TestValidationTrue(t *testing.T) {
	config := validConf
	err := config.validate()
	assert.NoError(t, err)
}

func TestValidationMissingDBFlag(t *testing.T) {
	config := validConf
	config.DBString = ""
	err := config.validate()
	assert.Error(t, err)
}

func TestValidationWrongLogLevel(t *testing.T) {
	config := validConf
	config.LogLevel = ""
	err := config.validate()
	assert.Error(t, err)

	config = validConf
	config.LogLevel = "whatever"
	err = config.validate()
	assert.Error(t, err)
}

func TestValidationLessThanZero(t *testing.T) {
	config := validConf
	config.Port = -1
	config.LogBackups = -1
	config.LogFileSize = -1
	config.LogAge = -1
	err := config.validate()
	assert.Error(t, err)
}
