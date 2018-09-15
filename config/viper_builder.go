package config

import (
	"github.com/spf13/viper"
	"log"
)

type viperBuilder struct {
	viper *viper.Viper
}

func NewViperBuilder() Builder {
	return &viperBuilder{
		viper: viper.New(),
	}
}

func (builder *viperBuilder) SetBasePath(path string) Builder {
	builder.viper.AddConfigPath(path)

	return builder
}

func (builder *viperBuilder) SetValue(key string, value interface{}) Builder {
	builder.viper.Set(key, value)

	return builder
}

func (builder *viperBuilder) AddJsonFile(path string) Builder {
	builder.viper.SetConfigType("json")
	builder.viper.SetConfigName(path)

	if err := builder.viper.MergeInConfig(); err != nil {
		log.Fatal(err)
	}

	return builder
}

func (builder *viperBuilder) AddEnvironmentVariables() Builder {
	builder.viper.SetEnvPrefix("apas")
	builder.viper.BindEnv("port")
	//builder.viper.AutomaticEnv()

	return builder
}

func (builder *viperBuilder) Build() (Configuration, error) {
	result := Configuration{}

	if err := builder.viper.Unmarshal(&result); err != nil {
		return Configuration{}, err
	}

	return result, nil
}

