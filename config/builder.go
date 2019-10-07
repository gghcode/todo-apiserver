package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// Builder build configuration wtih some options.
type Builder struct {
	pipelines [](func(*viper.Viper) error)
}

// NewBuilder return new builder instance.
func NewBuilder() *Builder {
	return &Builder{}
}

// AddConfigFile read config file from filepath.
func (builder *Builder) AddConfigFile(filepath string, optional bool) *Builder {
	builder.pipelines = append(builder.pipelines, func(v *viper.Viper) error {
		v.SetConfigFile(filepath)

		err := v.MergeInConfig()
		if err != nil {
			if _, ok := err.(*os.PathError); ok && !optional {
				return err
			}
		}

		return nil
	})

	return builder
}

// BindEnvs bind environment variables.
func (builder *Builder) BindEnvs(prefix string) *Builder {
	builder.pipelines = append(builder.pipelines, func(v *viper.Viper) error {
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

// Build return new configuration instance.
func (builder *Builder) Build() (Configuration, error) {
	result := DefaultConfig
	v := viper.New()

	for _, pipeline := range builder.pipelines {
		if err := pipeline(v); err != nil {
			return result, err
		}
	}

	if err := v.Unmarshal(&result); err != nil {
		return result, err
	}

	return result, nil
}
