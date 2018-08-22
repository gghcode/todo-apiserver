package app

type Configuration struct {
	ListenPort int `mapstructure:"LISTEN_PORT"`
}

type ConfigurationBuilder interface {
	SetBasePath(path string)
	SetValue(key string, value interface{})

	AddJsonFile(path string)
	AddEnvironmentVariables()

	Build() (*Configuration, error)
}
