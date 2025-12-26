package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/caarlos0/env/v11"
	flag "github.com/spf13/pflag"
)

type Config struct {
	//Database connection string (e.g. postgres://lanty:lanty@database:5432/lanty?sslmode=disable)
	DBString string `env:"DB"`

	//Scheme of the API server [http, https]
	Scheme string `env:"SCHEME"`

	//Host of the API server
	Host string `env:"HOST"`

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
var defaultConfig = Config{
	DBString:          "",
	Scheme:            "http",
	Host:              "localhost",
	Port:              8080,
	LogLevel:          "info",
	LogBackups:        0,
	LogFileSize:       10,
	LogAge:            0,
	EnableFileLogging: false,
	PrintHelp:         false,
	PrintVersion:      false,
}

func ParseConfig(args ...string) (*Config, *flag.FlagSet, error) {
	argConfig, flagset, err := parseFromArgs(args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse args: %w", err)
	}

	envConfig, err := parseFromEnvs()
	if err != nil {
		return nil, flagset, fmt.Errorf("failed to parse envs: %w", err)
	}

	config, err := merge(*argConfig, *envConfig, flagset)
	if err != nil {
		return nil, flagset, fmt.Errorf("failed to merge arg and env config: %w", err)
	}

	err = config.validate(flagset)
	if err != nil {
		return nil, flagset, fmt.Errorf("failed to validate config: %w", err)
	}

	return config, flagset, nil
}

func parseFromArgs(arguments ...string) (*Config, *flag.FlagSet, error) {
	config := &Config{}
	flagset := flag.NewFlagSet("lantyd", flag.ContinueOnError)

	flagset.Usage = func() {
		fmt.Printf("Usage: %s [options]\nOptions:\n", filepath.Base(arguments[0]))
		flagset.PrintDefaults()
	}

	flagset.StringVar(&config.DBString, "db", defaultConfig.DBString, "[REQUIRED] Connection string to a postgres database (eg. 'postgres://lanty:lanty@localhost:5432/lanty?sslmode=disable')")
	flagset.StringVar(&config.Scheme, "scheme", defaultConfig.Scheme, "Scheme of the API server [http, https]")
	flagset.StringVar(&config.Host, "host", defaultConfig.Host, "Host of the server")
	flagset.IntVarP(&config.Port, "port", "p", defaultConfig.Port, "Port of the server")
	flagset.StringVar(&config.LogLevel, "loglevel", defaultConfig.LogLevel, "Log level of the application [disable, trace, debug, info, warning, error, panic, fatal]")
	flagset.IntVar(&config.LogBackups, "logbackup", defaultConfig.LogBackups, "Number of log rotation files to keep")
	flagset.IntVar(&config.LogFileSize, "logfilesize", defaultConfig.LogFileSize, "Size of log rotation files [MegaBytes]")
	flagset.IntVar(&config.LogAge, "logage", defaultConfig.LogAge, "Days to retain logs")
	flagset.BoolVar(&config.EnableFileLogging, "enablefilelogging", defaultConfig.EnableFileLogging, "Enables file logging to \"lanty.log\"")
	flagset.BoolVarP(&config.PrintHelp, "help", "h", defaultConfig.PrintHelp, "Prints help")
	flagset.BoolVar(&config.PrintVersion, "version", defaultConfig.PrintVersion, "Prints version information")

	err := flagset.Parse(arguments[1:])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse into config: %w", err)
	}
	return config, flagset, nil
}

func parseFromEnvs() (*Config, error) {
	config := defaultConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse into config: %w", err)
	}
	return &config, nil
}

func merge(arg Config, env Config, flagset *flag.FlagSet) (*Config, error) {
	config := env

	db := flagset.Lookup("db")
	if db != nil && db.Changed {
		config.DBString = arg.DBString
	}

	scheme := flagset.Lookup("scheme")
	if scheme != nil && scheme.Changed {
		config.Scheme = arg.Scheme
	}

	host := flagset.Lookup("host")
	if host != nil && host.Changed {
		config.Host = arg.Host
	}

	port := flagset.Lookup("port")
	if port != nil && port.Changed {
		config.Port = arg.Port
	}

	loglevel := flagset.Lookup("loglevel")
	if loglevel != nil && loglevel.Changed {
		config.LogLevel = arg.LogLevel
	}

	logbackup := flagset.Lookup("logbackup")
	if logbackup != nil && logbackup.Changed {
		config.LogBackups = arg.LogBackups
	}

	logfilesize := flagset.Lookup("logfilesize")
	if logfilesize != nil && logfilesize.Changed {
		config.LogFileSize = arg.LogFileSize
	}

	logage := flagset.Lookup("logage")
	if logage != nil && logage.Changed {
		config.LogAge = arg.LogAge
	}

	enablefilelogging := flagset.Lookup("enablefilelogging")
	if enablefilelogging != nil && enablefilelogging.Changed {
		config.EnableFileLogging = arg.EnableFileLogging
	}

	printhelp := flagset.Lookup("help")
	if printhelp != nil {
		config.PrintHelp = arg.PrintHelp
	}

	printversion := flagset.Lookup("version")
	if printversion != nil {
		config.PrintVersion = arg.PrintVersion
	}

	return &config, nil
}

func (config Config) validate(flagset *flag.FlagSet) error {
	var err error

	if config.PrintHelp || config.PrintVersion {
		return nil
	}

	db := flagset.Lookup("db")
	if db != nil && !db.Changed {
		err = errors.Join(err, errors.New("required flag was not set: db"))
	} else if db != nil && db.Changed && len(config.DBString) == 0 {
		err = errors.Join(err, errors.New("required flag is empty: db"))
	}
	if config.Scheme != "http" && config.Scheme != "https" {
		err = errors.Join(err, fmt.Errorf("unknown scheme: %s", config.Scheme))
	}

	if config.Host == "" {
		err = errors.Join(err, errors.New("host can not be empty"))
	}

	if config.Port < 0 {
		err = errors.Join(err, errors.New("port can not be less than 0"))
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

	if config.LogBackups < 0 {
		err = errors.Join(err, errors.New("flag logbackup can not be less than 0"))
	}

	if config.LogFileSize < 0 {
		err = errors.Join(err, errors.New("flag logfilesize can not be less than 0"))
	}

	if config.LogAge < 0 {
		err = errors.Join(err, errors.New("flag logage can not be less than 0"))
	}

	return err
}
