package main

import (
	cfg "datatom/pkg/config"
)

type config struct {
	ConfigFilePath string `conf:"flag:config_file_path,short:c,env:CONFIG_FILE_PATH"`
}

func newConfig() *config {
	return &config{}
}

func parse(args []string, c *config) error {
	return cfg.Configure(args, c, cfg.WithConfigFilePathField("ConfigFilePath"))
}
