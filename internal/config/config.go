package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"os"
	"strconv"
)

const (
	defaultServerAddr           = "localhost:8080"
	databaseURIDefault          = "postgres://user:password@localhost:5432/loyalty-system?sslmode=disable"
	accrualSystemAddressDefault = "http://localhost:4444"
	secretKeyDefault            = "secretKey"
	pollIntervalDefault         = 2
	RateLimitDefault            = 3
	waitDefault                 = 15
	updateTimeoutDefault        = 10
)

func InitConfig() (*Config, error) {
	var newConfig Config
	err := env.Parse(&newConfig)
	if err != nil {
		return nil, err
	}
	return &newConfig, err
}

func NewServerConfig() (*Config, error) {
	newConfig, err := InitConfig()
	if err != nil {
		return nil, err
	}

	serverFlagSet := flag.NewFlagSet("Server", flag.ExitOnError)
	serverAddr := serverFlagSet.String("a", defaultServerAddr, "input endpoint")
	DatabaseURI := serverFlagSet.String("d", databaseURIDefault, "infrastructure uri")
	accrualSystemAddress := serverFlagSet.String("r", accrualSystemAddressDefault, "accrual system address")
	secretKey := serverFlagSet.String("s", secretKeyDefault, "secret key")
	pollInterval := serverFlagSet.Uint("p", pollIntervalDefault, "poll interval")
	rateLimit := serverFlagSet.Uint("l", RateLimitDefault, "accrual system address")
	wait := serverFlagSet.Uint("w", waitDefault, "secret key")
	updateTimeout := serverFlagSet.Uint("u", updateTimeoutDefault, "secret key")
	err = serverFlagSet.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	if newConfig.ServerAddr == nil {
		newConfig.ServerAddr = serverAddr
	}
	if newConfig.DatabaseURI == nil {
		newConfig.DatabaseURI = DatabaseURI
	}
	if newConfig.AccrualSystemAddress == nil {
		newConfig.AccrualSystemAddress = accrualSystemAddress
	}
	if newConfig.SecretKey == nil {
		newConfig.SecretKey = secretKey
	}
	if newConfig.PollInterval == nil {
		newConfig.PollInterval = pollInterval
	}
	if newConfig.RateLimit == nil {
		newConfig.RateLimit = rateLimit
	}
	if newConfig.Wait == nil {
		newConfig.Wait = wait
	}
	if newConfig.UpdateTimeout == nil {
		newConfig.UpdateTimeout = updateTimeout
	}
	return newConfig, nil
}

type Config struct {
	ServerAddr           *string `env:"RUN_ADDRESS"`
	DatabaseURI          *string `env:"DATABASE_URI"`
	AccrualSystemAddress *string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKey            *string `env:"SECRET_KEY"`
	PollInterval         *uint   `env:"POLL_INTERVAL"`
	RateLimit            *uint   `env:"RATE_LIMIT"`
	Wait                 *uint   `env:"WAIT"`
	UpdateTimeout        *uint   `env:"UPDATE_TIMEOUT"`
}

func InitDefaultEnv() error {
	envDefaults := map[string]string{
		"RUN_ADDRESS":            defaultServerAddr,
		"DATABASE_URI":           databaseURIDefault,
		"ACCRUAL_SYSTEM_ADDRESS": accrualSystemAddressDefault,
		"SECRET_KEY":             secretKeyDefault,
		"POLL_INTERVAL":          strconv.Itoa(pollIntervalDefault),
		"RATE_LIMIT":             strconv.Itoa(RateLimitDefault),
		"WAIT":                   strconv.Itoa(waitDefault),
		"UPDATE_TIMEOUT":         strconv.Itoa(updateTimeoutDefault),
	}
	for k, v := range envDefaults {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
