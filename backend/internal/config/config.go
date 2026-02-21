package config

import (
	"strconv"

	"github.com/dvalkoff/gomessenger/backend/internal/utils/env"
)

type Config struct {
	HttpConfig             HttpConfig
	DbConfig               DbConfig
	LoggingConfig          LoggingConfig
	GracefulShutdownConfig GracefulShutdownConfig
}

type HttpConfig struct {
	Port           int
	ReadTimeoutMs  int
	WriteTimeoutMs int
	JwtSecret      string
	CorsAllowedURL string
}

func makeHttpConfig(configMap map[string]string) (HttpConfig, error) {
	httpPort, err := strconv.Atoi(configMap[env.HttpPort])
	readTimeout, err := strconv.Atoi(configMap[env.HttpReadTimeoutMs])
	writeTimeout, err := strconv.Atoi(configMap[env.HttpWriteTimeoutMs])
	if err != nil {
		return HttpConfig{}, err
	}
	return HttpConfig{
		Port:           httpPort,
		ReadTimeoutMs:  readTimeout,
		WriteTimeoutMs: writeTimeout,
		JwtSecret:      configMap[env.HttpSecutityJwtSecret],
		CorsAllowedURL: configMap[env.HttpSecurityCorsAllowedURL],
	}, nil
}

type DbConfig struct {
	ConnectionStr string
}

func makeDbConfig(configMap map[string]string) (DbConfig, error) {
	return DbConfig{
		ConnectionStr: configMap[env.PostgreSQLConnectionString],
	}, nil
}

type LoggingConfig struct {
	ServiceName string
	Level       string
	AddSource   bool
	Handler     string
}

func makeLoggingConfig(configMap map[string]string) (LoggingConfig, error) {
	addSource, err := strconv.ParseBool(configMap[env.LoggingAddSource])
	if err != nil {
		return LoggingConfig{}, err
	}
	return LoggingConfig{
		ServiceName: configMap[env.LoggingServiceName],
		Level:       configMap[env.LoggingMinLevel],
		AddSource:   addSource,
		Handler:     configMap[env.LoggingHandler],
	}, nil
}

type GracefulShutdownConfig struct {
	ShutdownTimeoutSec int
}

func makeGracefulShutdownConfig(configMap map[string]string) (GracefulShutdownConfig, error) {
	shutdownTimeout, err := strconv.Atoi(configMap[env.GracefulShutdownTimeoutSec])
	if err != nil {
		return GracefulShutdownConfig{}, err
	}
	return GracefulShutdownConfig{
		ShutdownTimeoutSec: shutdownTimeout,
	}, nil
}

func MakeAppConfig() (Config, error) {
	configMap, err := env.LoadEnvFromFile()
	if err != nil {
		return Config{}, err
	}
	httpConfig, err := makeHttpConfig(configMap)
	dbConfig, err := makeDbConfig(configMap)
	loggingConfig, err := makeLoggingConfig(configMap)
	gracefulShutdownConfig, err := makeGracefulShutdownConfig(configMap)
	if err != nil {
		return Config{}, err
	}
	return Config{
		HttpConfig:             httpConfig,
		DbConfig:               dbConfig,
		LoggingConfig:          loggingConfig,
		GracefulShutdownConfig: gracefulShutdownConfig,
	}, nil
}
