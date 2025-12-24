package config

import (
	"errors"
	"fmt"

	"dario.cat/mergo"
	"github.com/caarlos0/env/v11"
	flag "github.com/spf13/pflag"
)

type Config struct {
	//Database connection string (e.g. postgres://lanty:lanty@database:5432/lanty?sslmode=disable)
	DBString string `env:"DB"`

	//Port the API server will listen
	Port int `env:"PORT"`

	//Log level [disable, trace, debug, info, warning, error, panic, fatal]
	LogLevel string `env:"LOG_LEVEL"`

	//Maximum number of old log files to retain
	LogBackups int `env:"LOG_BACKUPS"`

	//Maximum file size (megabytes) of log file before it will be rotated if logging is enabled
	LogFileSize int `env:"LOG_FILESIZE"`

	//Maximum number of days to retain old log files
	LogAge int `env:"LOG_AGE"`

	//Enables logging to file
	EnableFileLogging bool `env:"ENABLE_FILE_LOGGING"`

	//Determines if the help string is printed on the CLI
	PrintHelp bool

	//Determines if the version string is printed on the CLI
	PrintVersion bool
}

// Defaults for configuration values where necessary
const (
	portDefault        = 8080
	loglevelDefault    = "info"
	logfilesizeDefault = 10
)

// Parses the configuration from the given arguments and environment variables and validates it.
// Arguments take precedence over environment variables and will overwrite them.
//
// If the validation of the configuration fails the parsed configuration is returned for debugging purposes.
func Parse(args ...string) (Config, error) {
	argConfig := parseFromArgs(args...)
	envConfig := parseFromEnvs()

	config, err := merge(argConfig, envConfig)
	if err != nil {
		return config, fmt.Errorf("failed to merge arg and env config: %w", err)
	}

	err = config.validate()
	if err != nil {
		return config, fmt.Errorf("failed to validate config: %w", err)
	}

	return config, nil
}

// Parses a configuration from the given arguments without the command name (normally os.Args[1:]) and sets default
// values if arguments were not set.
// This should be called in conjunction with parseFromEnvs() and merge().
func parseFromArgs(arguments ...string) Config {
	config := Config{}
	flagset := flag.FlagSet{}

	flagset.StringVar(&config.DBString, "db", "", "[REQUIRED] Connection string to a postgres database (eg. 'postgres://lanty:lanty@localhost:5432/lanty?sslmode=disable')")
	flagset.IntVarP(&config.Port, "port", "p", portDefault, "Port of the server")
	flagset.StringVar(&config.LogLevel, "loglevel", loglevelDefault, "Log level of the application [disable, trace, debug, info, warning, error, panic, fatal]")
	flagset.IntVar(&config.LogBackups, "logbackup", 0, "Number of log rotation files to keep")
	flagset.IntVar(&config.LogFileSize, "logfilesize", logfilesizeDefault, "Size of log rotation files [MegaBytes]")
	flagset.IntVar(&config.LogAge, "logage", 0, "Days to retain logs")
	flagset.BoolVar(&config.EnableFileLogging, "enablefilelogging", false, "Enables file logging to \"lanty.log\"")
	flagset.BoolVarP(&config.PrintHelp, "help", "h", false, "Prints help")
	flagset.BoolVar(&config.PrintVersion, "version", false, "Prints version information")

	flagset.Parse(arguments)
	return config
}

// Parses a configuration from environment variables. This does not set any default values as parseFromArgs() does.
// This should be called in conjunction with parseFromArgs() and merge().
func parseFromEnvs() Config {
	config := Config{}
	env.Parse(&config)
	return config
}

// Merges env into argument configuration filling all empty values and overwriting default values if defined by env
// configuration.
// This should be called in conjunction with parseFromArgs() and parseFromEnvs().
func merge(arg, env Config) (Config, error) {
	err := mergo.Merge(&arg, env)
	if err != nil {
		return arg, err
	}

	if arg.Port == portDefault && env.Port != 0 {
		arg.Port = env.Port
	}
	if arg.LogLevel == loglevelDefault && env.LogLevel != "" {
		arg.LogLevel = env.LogLevel
	}
	if arg.LogFileSize == logfilesizeDefault && env.LogFileSize != 0 {
		arg.LogFileSize = env.LogFileSize
	}

	return arg, nil
}

// Validates the configuration for the presence of all required fields and the correctness of all values.
func (config Config) validate() error {
	var err error

	if len(config.DBString) == 0 {
		err = errors.Join(err, errors.New("required flag is empty: db"))
	}

	if config.Port < 0 {
		err = errors.Join(err, errors.New("port can not be less than 0"))
	}

	if config.LogBackups < 0 {
		err = errors.Join(err, errors.New("flag logbackup can not be less than 0"))
	}

	if config.LogFileSize < 0 {
		err = errors.Join(err, errors.New("flag logfilesize can not be less than 0"))
	}

	if config.LogAge < 0 {
		err = errors.Join(err, errors.New("flag logage can not be less than 0"))
	}

	if config.LogLevel != "disable" &&
		config.LogLevel != "trace" &&
		config.LogLevel != "debug" &&
		config.LogLevel != "info" &&
		config.LogLevel != "warning" &&
		config.LogLevel != "error" &&
		config.LogLevel != "panic" &&
		config.LogLevel != "fatal" {
		err = errors.Join(err, fmt.Errorf("unknown loglevel: %s", config.LogLevel))
	}

	return err
}
