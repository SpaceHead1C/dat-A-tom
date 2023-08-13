package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/ardanlabs/conf"
)

const (
	configPathTag = "config_file"
	tomlParam     = "toml"
	zeroTag       = "zero"
	noParam       = "no"
)

type settings struct {
	configFilePathField string
	envNamespace        string
}

type option func(s *settings)

func Configure(args []string, cfg any, opts ...option) error {
	if cfg == nil {
		return fmt.Errorf("configuration object parameter can't be nil")
	}
	s := settings{}
	for _, o := range opts {
		o(&s)
	}
	if err := conf.Parse(args, s.envNamespace, cfg); err != nil {
		return fmt.Errorf("arguments parse error: %w", err)
	}
	rObject := reflect.ValueOf(cfg)
	if rObject.Kind() == reflect.Ptr {
		rObject = rObject.Elem()
	}
	rType := rObject.Type()
	for i := 0; i < rObject.NumField(); i++ {
		field := rObject.Field(i)
		fieldType := rType.Field(i)

		if s.configFilePathField != "" && fieldType.Name == s.configFilePathField && fieldType.Type.Kind() == reflect.String {
			if _, err := toml.DecodeFile(field.String(), cfg); err != nil {
				return fmt.Errorf("config file parse error: %w", err)
			}
		}
		if strings.TrimSpace(fieldType.Tag.Get(zeroTag)) == noParam && field.IsZero() {
			return fmt.Errorf("option %s contains zero value", fieldType.Name)
		}
	}
	return nil
}

func WithConfigFilePathField(fieldName string) option {
	return func(s *settings) {
		s.configFilePathField = fieldName
	}
}

func WithEnvNamespace(namespace string) option {
	return func(s *settings) {
		s.envNamespace = namespace
	}
}
