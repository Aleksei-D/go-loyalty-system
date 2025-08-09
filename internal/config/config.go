package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"os"
)

const (
	defaultServerAddr           = "localhost:8080"
	databaseUriDefault          = "infrastructure://user:password@localhost:5432/loyalty-system"
	accrualSystemAddressDefault = "localhost:22222"
	secretKeyDefault            = "secretKey"
	pollIntervalDefault         = 2
	RateLimitDefault            = 3
	waitDefault                 = 15
)

func NewServerConfig() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	serverFlagSet := flag.NewFlagSet("Server", flag.ExitOnError)
	serverAddr := serverFlagSet.String("a", defaultServerAddr, "input endpoint")
	databaseUri := serverFlagSet.String("d", databaseUriDefault, "infrastructure uri")
	accrualSystemAddress := serverFlagSet.String("r", accrualSystemAddressDefault, "accrual system address")
	secretKey := serverFlagSet.String("s", secretKeyDefault, "secret key")
	pollInterval := serverFlagSet.Uint("p", pollIntervalDefault, "poll interval")
	rateLimit := serverFlagSet.Uint("l", RateLimitDefault, "accrual system address")
	wait := serverFlagSet.Uint("w", waitDefault, "secret key")
	err = serverFlagSet.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	if config.ServerAddr == nil {
		config.ServerAddr = serverAddr
	}
	if config.DatabaseUri == nil {
		config.DatabaseUri = databaseUri
	}
	if config.AccrualSystemAddress == nil {
		config.AccrualSystemAddress = accrualSystemAddress
	}
	if config.SecretKey == nil {
		config.SecretKey = secretKey
	}
	if config.PollInterval == nil {
		config.PollInterval = pollInterval
	}
	if config.RateLimit == nil {
		config.RateLimit = rateLimit
	}
	if config.Wait == nil {
		config.Wait = wait
	}
	return &config, nil
}

type Config struct {
	ServerAddr           *string `env:"RUN_ADDRESS"`
	DatabaseUri          *string `env:"DATABASE_URI"`
	AccrualSystemAddress *string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKey            *string `env:"SECRET_KEY"`
	PollInterval         *uint   `env:"POLL_INTERVAL"`
	RateLimit            *uint   `env:"RATE_LIMIT"`
	Wait                 *uint   `env:"WAIT"`
}
