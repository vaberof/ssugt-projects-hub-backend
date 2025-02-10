package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	_appName = "ssugt-projects-hub-backend"

	_portKey    = "port"
	_timeoutKey = "timeout"
	_secretKey  = "secret_key"
)

func Init() {
	if err := godotenv.Load(".config/app.env"); err != nil {
		panic(fmt.Errorf("failed to read environment variables file: %w", err))
	}

	viper.AddConfigPath(".config/")
	viper.SetConfigType("json")
	viper.SetConfigName(configName())
	err := viper.ReadInConfig()
	viper.WatchConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}
}

func Port() int {
	return viper.GetInt(_portKey)
}

func PostgresConnection() string {
	return fmt.Sprintf("%s&search_path=%s", viper.GetString("postgres"), configName())
}

func Timeout() time.Duration {
	return time.Duration(viper.GetInt(_timeoutKey)) * time.Second
}

func SecretKey() string {
	return viper.GetString(_secretKey)
}

func IsProd() bool {
	serviceEnv := os.Getenv(ServiceEnvVarName)
	return serviceEnv == EnvProd
}

func ApplicationName() string {
	return _appName
}

func configName() string {
	serviceEnv := os.Getenv(ServiceEnvVarName)

	if serviceEnv == "prod" {
		return EnvProd
	}

	return EnvDev
}
