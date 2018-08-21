package app

import (
	"github.com/spf13/viper"
	"log"
)

type viperConfigurationBuilder struct {
	viper *viper.Viper
}

func NewViperConfigurationBuilder() viperConfigurationBuilder {
	return viperConfigurationBuilder{viper: viper.New()}
}

func (builder *viperConfigurationBuilder) SetBasePath(path string) {
	builder.viper.AddConfigPath(path)
}

func (builder *viperConfigurationBuilder) SetValue(key string, value interface{}) {
	builder.viper.Set(key, value)
}

func (builder *viperConfigurationBuilder) AddJsonFile(path string) {
	builder.viper.SetConfigType("json")
	builder.viper.SetConfigName(path)

	if err := builder.viper.MergeInConfig(); err != nil {
		log.Fatal(err)
	}
}

func (builder *viperConfigurationBuilder) AddEnvironmentVariables() {
	builder.viper.SetEnvPrefix("apas")
	builder.viper.BindEnv("port")
	//builder.viper.AutomaticEnv()
}

func (builder *viperConfigurationBuilder) Build() (Configuration, error) {
	result := Configuration{}

	if err := builder.viper.Unmarshal(&result); err != nil {
		return nil, err
	}

	return result, nil
}
