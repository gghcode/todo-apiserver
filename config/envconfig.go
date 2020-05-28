package config

import "github.com/vrischmann/envconfig"

// FromEnvs provide configuration by envs
func FromEnvs() (Configuration, error) {
	var cfg Configuration
	if err := envconfig.Init(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
