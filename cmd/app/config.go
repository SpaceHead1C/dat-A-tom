package main

import (
	cfg "datatom/pkg/config"
)

type config struct {
	ConfigFilePath string `conf:"flag:config_file_path,short:c,env:CONFIG_FILE_PATH"`

	RESTPort uint `conf:"flag:rest_port,short:r,env:REST_PORT" toml:"rest_port" zero:"no"`
	RESTTimeoutSec uint `conf:"flag:rest_timeout,short:r,env:REST_TIMEOUT" toml:"rest_timeout"`

	PostgresAddress  string `conf:"flag:postgres_address,env:POSTGRES_ADDRESS" toml:"postgres_address" zero:"no"`
	PostgresPort     uint   `conf:"flag:postgres_port,env:POSTGRES_PORT" toml:"postgres_port" zero:"no"`
	PostgresDBName   string `conf:"flag:postgres_db_name,env:POSTGRES_DB_NAME" toml:"postgres_db_name" zero:"no"`
	PostgresUser     string `conf:"flag:postgres_user,env:POSTGRES_USER" toml:"postgres_user" zero:"no"`
	PostgresPassword string `conf:"flag:postgres_password,env:POSTGRES_PASSWORD" toml:"postgres_password" zero:"no"`
}

func newConfig() *config {
	return &config{}
}

func parse(args []string, c *config) error {
	return cfg.Configure(args, c, cfg.WithConfigFilePathField("ConfigFilePath"))
}
