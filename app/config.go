package app

import "github.com/spf13/viper"

type Config struct {
	Port int
	MongoConnStr string
}

func loadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var result *Config

	if err := viper.Unmarshal(&result); err != nil {
		return nil, err
	}

	return result, nil
}