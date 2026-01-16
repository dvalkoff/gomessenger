package config

type Config struct {
	httpConfig HttpConfig
	dbConfig DbConfig
}

type HttpConfig struct {
	Port int
	ReadTimeoutMs int
	WriteTimeoutMs int
	ShutdownTimeoutSec int
}

type DbConfig struct {
	ConnectionStr string
}
