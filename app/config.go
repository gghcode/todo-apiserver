package app

type Configuration = map[string]interface{}

type ConfigurationBuilder interface {
	SetBasePath(path string)
	SetValue(key string, value interface{})

	AddJsonFile(path string)
	AddEnvironmentVariables()

	Build() (Configuration, error)
}
