package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const connectionTimeout = 5 * time.Second

type MongoDatabaseConfig struct {
	AppName  string `yaml:"app-name"`
	URI      string `yaml:"uri"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	PoolSize uint64 `yaml:"pool-size"`
}

type ManagedDatabase struct {
	Db      *mongo.Database
	errorCh chan error
}

func New(config *MongoDatabaseConfig) (*ManagedDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	clientOptions := options.Client()

	if config.User != "" && config.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: config.User,
			Password: config.Password,
		})
	}

	clientOptions.ApplyURI(config.URI)
	clientOptions.SetAppName(config.AppName)
	clientOptions.SetTimeout(time.Second * 5)
	clientOptions.SetConnectTimeout(time.Second * 5)
	clientOptions.SetMaxConnIdleTime(time.Second * 5)
	clientOptions.SetMaxPoolSize(config.PoolSize)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	errorCh := make(chan error)

	//go verifyClientConnectionAbility(client, errorCh)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("FAILED TO PING")
		return nil, fmt.Errorf("failed to ping mongoDb: %w", err)
	}

	mongoDb := client.Database(config.Database)
	if mongoDb == nil {
		return nil, errors.New("failed to get database")
	}

	managedDatabase := ManagedDatabase{
		Db:      mongoDb,
		errorCh: errorCh,
	}

	log.Println("database is running")

	return &managedDatabase, nil
}

func verifyClientConnectionAbility(client *mongo.Client, errorCh chan<- error) {
	pingTimeout := time.Second * 5
	pingInterval := time.Second * 5

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	timer := time.NewTimer(pingInterval)
	defer timer.Stop()
	for {
		if err := client.Ping(ctx, nil); err != nil {
			log.Println("FAILED TO PING")
			errorCh <- fmt.Errorf("failed to ping mongodb while verifying client connection: %w", err)
			return
		}

		select {
		case <-timer.C:
			timer.Reset(pingInterval)
		}
	}
}

func (mongo *ManagedDatabase) Error() <-chan error {
	return mongo.errorCh
}

func (mongo *ManagedDatabase) Disconnect(ctx context.Context) error {
	if mongo.Db.Client() == nil {
		return nil
	}
	return mongo.Db.Client().Disconnect(ctx)
}
