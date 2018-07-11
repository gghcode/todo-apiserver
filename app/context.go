package app

type Context struct {
	Config *Config
}

func LoadContext() (*Context, error) {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	result := &Context{ config }

	return result, nil
}
