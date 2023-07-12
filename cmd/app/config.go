package main

import (
	cfg "datatom/pkg/config"
)

type config struct {
	ConfigFilePath string `conf:"flag:config_file_path,short:c,env:CONFIG_FILE_PATH"`

	Title       string `conf:"flag:title,env:TITLE" toml:"title"`
	Description string `conf:"flag:description,env:DESCRIPTION" toml:"description"`

	RESTPort       uint `conf:"flag:rest_port,short:r,env:REST_PORT" toml:"rest_port" zero:"no"`
	RESTTimeoutSec uint `conf:"flag:rest_timeout,short:r,env:REST_TIMEOUT" toml:"rest_timeout"`

	PostgresAddress  string `conf:"flag:postgres_address,env:POSTGRES_ADDRESS" toml:"postgres_address" zero:"no"`
	PostgresPort     uint   `conf:"flag:postgres_port,env:POSTGRES_PORT" toml:"postgres_port" zero:"no"`
	PostgresDBName   string `conf:"flag:postgres_db_name,env:POSTGRES_DB_NAME" toml:"postgres_db_name" zero:"no"`
	PostgresUser     string `conf:"flag:postgres_user,env:POSTGRES_USER" toml:"postgres_user" zero:"no"`
	PostgresPassword string `conf:"flag:postgres_password,env:POSTGRES_PASSWORD" toml:"postgres_password" zero:"no"`

	DatawayGRPCAddress string `conf:"flag:dataway_grpc_address,env:DATAWAY_GRPC_ADDRESS" toml:"dataway_grpc_address"`
	DatawayGRPCPort    uint   `conf:"flag:dataway_grpc_port,env:DATAWAY_GRPC_PORT" toml:"dataway_grpc_port"`

	RMQAddress      string `conf:"flag:rmq_address,env:RMQ_ADDRESS" toml:"rmq_address"`
	RMQPort         uint   `conf:"flag:rmq_port,env:RMQ_PORT" toml:"rmq_port"`
	RMQUser         string `conf:"flag:rmq_user,env:RMQ_USER" toml:"rmq_user"`
	RMQPassword     string `conf:"flag:rmq_password,env:RMQ_PASSWORD" toml:"rmq_password"`
	RMQVHost        string `conf:"flag:rmq_vhost,env:RMQ_VHOST" toml:"rmq_vhost"`
	RMQConsumeQueue string `conf:"flag:rmq_consume_queue,env:RMQ_CONSUME_QUEUE" toml:"rmq_consume_queue"`
	RMQDLE          string `conf:"flag:rmq_dle,env:RMQ_DLE" toml:"rmq_dle"`
}

func newConfig() *config {
	return &config{}
}

func parse(args []string, c *config) error {
	return cfg.Configure(args, c, cfg.WithConfigFilePathField("ConfigFilePath"))
}
