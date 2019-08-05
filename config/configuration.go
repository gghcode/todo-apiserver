package config

// Configuration is config type.
type Configuration struct {
	Addr     string `mapstructure:"addr"`
	BasePath string `mapstructure:"basePath"`
}
