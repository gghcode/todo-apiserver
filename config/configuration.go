package config

import "github.com/spf13/viper"

type Configuration struct {
	Port int
	MongoConnStr string
}

func Load() (*Configuration, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var result *Configuration

	if err := viper.Unmarshal(&result); err != nil {
		return nil, err
	}

	return result, nil
}
