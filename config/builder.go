package config

import "github.com/spf13/viper"

// Builder build configuration wtih some options.
type Builder struct {
	viper *viper.Viper
}

// NewBuilder return new builder instance.
func NewBuilder() *Builder {
	return &Builder{
		viper: viper.New(),
	}
}

// AddConfigFile read config file from filepath.
func (builder *Builder) AddConfigFile(filepath string) *Builder {
	builder.viper.SetConfigFile(filepath)
	builder.viper.MergeInConfig()

	return builder
}

// UseEnv use environment variables.
func (builder *Builder) UseEnv(prefix string) *Builder {
	builder.viper.SetEnvPrefix(prefix)
	builder.viper.AutomaticEnv()

	return builder
}

// Build return built new configuration instance.
func (builder *Builder) Build() (*Configuration, error) {
	var result Configuration

	if err := builder.viper.Unmarshal(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
