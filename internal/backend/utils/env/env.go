package env

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	HttpPort                   = "HTTP_PORT"
	HttpReadTimeoutMs          = "HTTP_READ_TIMEOUT_MS"
	HttpWriteTimeoutMs         = "HTTP_WRITE_TIMEOUT_MS"
	HttpSecutityJwtSecret      = "HTTP_SECURITY_JWT_SECRET"
	HttpSecurityCorsAllowedURL = "HTTP_SECURTITY_CORS_ALLOWED_URL"

	LoggingServiceName = "LOGGING_SERVICE_NAME"
	LoggingMinLevel    = "LOGGING_MIN_LEVEL"
	LoggingAddSource   = "LOGGING_ADD_SOURCE"
	LoggingHandler     = "LOGGING_HANDLER"

	PostgreSQLConnectionString = "POSTGRESQL_CONNECTION_STRING"

	GracefulShutdownTimeoutSec = "GRACEFUL_SHUTDOWN_TIMEOUT_SEC"
)

const (
	ENVIRONMENT_VAR = "ENVIRONMENT"
	PROD_ENV_NAME   = "PROD"

	PROD_ENV_FILE = "prod.env"
	DEV_ENV_FILE  = ".env"
)

func LoadEnvFromFile() (map[string]string, error) {
	var envFileName string
	switch os.Getenv(ENVIRONMENT_VAR) {
	case PROD_ENV_NAME:
		envFileName = PROD_ENV_FILE
	default:
		envFileName = DEV_ENV_FILE
	}
	return godotenv.Read(envFileName)
}
