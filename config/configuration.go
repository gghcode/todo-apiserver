package config

// DefaultConfig is default config
var DefaultConfig = Configuration{
	Addr:     ":8080",
	BasePath: "api",

	Cors:     DefaultCorsConfig,
	Postgres: DefaultPostgresConfig,
	Redis:    DefaultRedisConfig,
	Jwt:      DefaultJwtConfig,
}

// Configuration is config type.
type Configuration struct {
	Addr     string `mapstructure:"addr"`
	BasePath string `mapstructure:"basePath"`

	Cors     CorsConfig     `mapstructure:"cors"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Jwt      JwtConfig      `mapstructure:"jwt"`
}

// DefaultCorsConfig is default config
var DefaultCorsConfig = CorsConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
}

// CorsConfig godoc
type CorsConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowMethods []string `mapstructure:"allow_methods"`
}

// DefaultPostgresConfig is default config
var DefaultPostgresConfig = PostgresConfig{
	Driver:   "postgres",
	Host:     "localhost",
	Port:     "5432",
	User:     "postgres",
	Name:     "postgres",
	Password: "postgres",
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

// DefaultRedisConfig is config
var DefaultRedisConfig = RedisConfig{
	Addr: "localhost:6378",
}

// RedisConfig is redis config
type RedisConfig struct {
	Addr string `mapstructure:"addr"`
}

// DefaultJwtConfig is default config
var DefaultJwtConfig = JwtConfig{
	SecretKey:           "debugKey",
	AccessExpiresInSec:  3600,
	RefreshExpiresInSec: 60 * 60 * 24,
}

// JwtConfig is jwt config
type JwtConfig struct {
	SecretKey           string `mapstructure:"secret_key"`
	AccessExpiresInSec  int64  `mapstructure:"access_expires_sec"`
	RefreshExpiresInSec int64  `mapstructure:"refresh_expires_sec"`
}
