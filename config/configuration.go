package config

import "github.com/spf13/viper"

type Configuration struct {
	Port int
	MongoConnStr string
}

func Load() Configuration {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	result := &Configuration{}

	if err := viper.Unmarshal(&result); err != nil {
		panic(err)
	}

	return *result
}
