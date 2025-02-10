package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
)

func NewPgx(ctx context.Context, connectionString string) *sqlx.DB {
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Fatal(err)
	}

	db := sqlx.NewDb(stdlib.OpenDBFromPool(pool), "pgx")
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func NewOld(connectionString string) *sqlx.DB {
	db := sqlx.MustConnect("postgres", connectionString)

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	db.SetMaxOpenConns(20)

	return db
}
