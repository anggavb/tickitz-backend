package config

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPsql() (*pgxpool.Pool, error) {
	pgc, _ := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	return pgc, pgc.Ping(context.Background())
}
