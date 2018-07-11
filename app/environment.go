package app

type Environment struct {
	Config *Config
}

func LoadEnvironment() (*Environment, error) {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	result := &Environment{ config }

	return result, nil
}
