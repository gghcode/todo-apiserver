package config

// Configuration is config type.
type Configuration struct {
	Addr       string `mapstructure:"addr"`
	BasePath   string `mapstructure:"basePath"`
	BcryptCost int    `mapstructure:"bcryptCost"`

	Cors     CorsConfig     `mapstructure:"cors"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Jwt      JwtConfig      `mapstructure:"jwt"`
}

// DefaultConfig is default config
func DefaultConfig() Configuration {
	return Configuration{
		Addr:       ":8080",
		BasePath:   "api",
		BcryptCost: 12,

		Cors:     DefaultCorsConfig(),
		Postgres: DefaultPostgresConfig(),
		Redis:    DefaultRedisConfig(),
		Jwt:      DefaultJwtConfig(),
	}
}

// CorsConfig godoc
type CorsConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowMethods []string `mapstructure:"allow_methods"`
}

// DefaultCorsConfig is default config
func DefaultCorsConfig() CorsConfig {
	return CorsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}
}

// PostgresConfig is postgres config
type PostgresConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
}

// DefaultPostgresConfig is default config
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Name:     "postgres",
		Password: "postgres",
	}
}

// RedisConfig is redis config
type RedisConfig struct {
	Addr string `mapstructure:"addr"`
}

// DefaultRedisConfig is config
func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Addr: "localhost:6378",
	}
}

// JwtConfig is jwt config
type JwtConfig struct {
	SecretKey           string `mapstructure:"secret_key"`
	AccessExpiresInSec  int64  `mapstructure:"access_expires_sec"`
	RefreshExpiresInSec int64  `mapstructure:"refresh_expires_sec"`
}

// DefaultJwtConfig is default config
func DefaultJwtConfig() JwtConfig {
	return JwtConfig{
		SecretKey:           "debugKey",
		AccessExpiresInSec:  3600,
		RefreshExpiresInSec: 60 * 60 * 24,
	}
}
