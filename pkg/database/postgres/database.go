package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgDatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type ManagedDatabase struct {
	Db *gorm.DB
}

func New(config *PgDatabaseConfig) (*ManagedDatabase, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", config.Host, config.Port, config.User, config.Password, config.Database)
	gormDb, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	managedDatabase := ManagedDatabase{
		Db: gormDb,
	}

	return &managedDatabase, nil
}
