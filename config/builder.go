package config

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// Builder build configuration wtih some options.
type Builder struct {
	viper *viper.Viper
}

// NewBuilder return new builder instance.
func NewBuilder() *Builder {
	return &Builder{
		viper: viper.New(),
	}
}

// AddConfigFile read config file from filepath.
func (builder *Builder) AddConfigFile(filepath string) *Builder {
	builder.viper.SetConfigFile(filepath)
	builder.viper.MergeInConfig()

	return builder
}

// BindEnvs bind environment variables.
func (builder *Builder) BindEnvs(prefix string) *Builder {
	builder.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	builder.viper.SetEnvPrefix(prefix)

	bindEnvsToViper(builder.viper, Configuration{})

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
	var result Configuration

	if err := builder.viper.Unmarshal(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
