package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/vaberof/ssugt-projects/pkg/database/mongodb"
)

var appConfigPaths = flag.String("config.files", "not-found.yaml", "List of application config files separated by comma")
var environmentVariablesPath = flag.String("env.vars.file", "not-found.env", "Path to environment variables file")

func main() {
	flag.Parse()

	if err := loadEnvironmentVariables(); err != nil {
		panic(err)
	}

	appConfig := getAppConfig(*appConfigPaths)
	
	_, err := mongodb.New(&appConfig.Database)
	if err != nil {
		panic(err)
	}
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}
