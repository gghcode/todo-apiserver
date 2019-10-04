package config

// Configuration is config type.
type Configuration struct {
	Addr     string `mapstructure:"addr"`
	BasePath string `mapstructure:"basePath"`

	Cors     CorsConfig     `mapstructure:"cors"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Jwt      JwtConfig      `mapstructure:"jwt"`
}

// CorsConfig godoc
type CorsConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowMethods []string `mapstructure:"allow_methods"`
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

// RedisConfig is redis config
type RedisConfig struct {
	Addr string `mapstructure:"addr"`
}

// JwtConfig is jwt config
type JwtConfig struct {
	SecretKey           string `mapstructure:"secret_key"`
	AccessExpiresInSec  int64  `mapstructure:"access_expires_sec"`
	RefreshExpiresInSec int64  `mapstructure:"refresh_expires_sec"`
}
