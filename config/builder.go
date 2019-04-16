package config

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// Builder build configuration wtih some options.
type Builder struct {
	process [](func(*viper.Viper) error)
}

// NewBuilder return new builder instance.
func NewBuilder() *Builder {
	return &Builder{}
}

// AddConfigFile read config file from filepath.
func (builder *Builder) AddConfigFile(filepath string) *Builder {
	builder.process = append(builder.process, func(v *viper.Viper) error {
		v.SetConfigFile(filepath)
		return v.MergeInConfig()
	})

	return builder
}

// BindEnvs bind environment variables.
func (builder *Builder) BindEnvs(prefix string) *Builder {
	builder.process = append(builder.process, func(v *viper.Viper) error {
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.SetEnvPrefix(prefix)

		bindEnvsToViper(v, Configuration{})

		return nil
	})

	return builder
}

func bindEnvsToViper(viper *viper.Viper, iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			bindEnvsToViper(viper, v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}

// Build return built new configuration instance.
func (builder *Builder) Build() (*Configuration, error) {
	v := viper.New()

	for _, ps := range builder.process {
		if err := ps(v); err != nil {
			return nil, err
		}
	}

	var result Configuration

	if err := v.Unmarshal(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
