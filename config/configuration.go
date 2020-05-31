package config

// Configuration is config type.
type Configuration struct {
	Port                       int
	BasePath                   string
	BcryptCost                 int
	GracefulShutdownTimeoutSec int

	JwtSecretKey           string
	JwtAccessExpiresInSec  int64
	JwtRefreshExpiresInSec int64

	CorsAllowOrigins []string
	CorsAllowMethods []string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresName     string
	PostgresPassword string

	RedisAddr     string
	RedisPassword string
}
