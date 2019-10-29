package config

// Builder build configuration with some options.
type Builder interface {
	AddConfigFile(filepath string, optional bool) Builder
	BindEnvs(prefix string) Builder
	Build() (Configuration, error)
}
