package app

type Configuration = map[string]interface{}

type ConfigurationBuilder interface {
	SetBasePath(path string)
	AddJsonFile(path string)
	AddEnvironmentVariables()
	Build() (Configuration, error)
}