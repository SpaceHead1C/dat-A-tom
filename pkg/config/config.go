package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/ardanlabs/conf"
	"reflect"
	"strings"
)

const (
	configPathTag = "config_file"
	tomlParam     = "toml"
	zeroTag       = "zero"
	noParam       = "no"
)

type settings struct {
	configFilePathField string
}

type option func(s *settings)

func Configure(args []string, cfg any, opts ...option) error {
	if cfg == nil {
		return fmt.Errorf("configuration object parameter can't be nil")
	}
	if err := conf.Parse(args, "DAT_A_TOM", cfg); err != nil {
		return fmt.Errorf("arguments parse error: %w", err)
	}
	s := settings{}
	for _, o := range opts {
		o(&s)
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
		if strings.TrimSpace(rType.Field(i).Tag.Get(zeroTag)) == noParam && field.IsZero() {
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
