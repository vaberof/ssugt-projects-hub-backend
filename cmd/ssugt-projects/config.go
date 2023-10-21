package main

import (
	"errors"
	"github.com/vaberof/ssugt-projects/pkg/config"
	"github.com/vaberof/ssugt-projects/pkg/database/mongodb"
	"github.com/vaberof/ssugt-projects/pkg/http/httpserver"
)

type AppConfig struct {
	Server   httpserver.ServerConfig
	Database mongodb.MongoDatabaseConfig
}

func getAppConfig(sources ...string) AppConfig {
	config, err := tryGetAppConfig(sources...)
	if err != nil {
		panic(err)
	}

	if config == nil {
		panic(errors.New("config cannot be nil"))
	}

	return *config
}

func tryGetAppConfig(sources ...string) (*AppConfig, error) {
	if len(sources) == 0 {
		return nil, errors.New("at least 1 source must be set for app config")
	}

	provider := config.MergeConfigs(sources)

	var serverConfig httpserver.ServerConfig
	err := config.ParseConfig(provider, "app.http.server", &serverConfig)
	if err != nil {
		return nil, err
	}

	var mongoDatabaseConfig mongodb.MongoDatabaseConfig
	err = config.ParseConfig(provider, "app.mongodb", &mongoDatabaseConfig)
	if err != nil {
		return nil, err
	}

	config := AppConfig{
		Server:   serverConfig,
		Database: mongoDatabaseConfig,
	}

	return &config, nil
}
