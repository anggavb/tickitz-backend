package config

import (
	"context"
<<<<<<< HEAD
	"fmt"
=======
>>>>>>> b9ee6f3b7daa7e17199dec072791cf7dbe5d369b
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPsql() (*pgxpool.Pool, error) {
	pgc, _ := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	return pgc, pgc.Ping(context.Background())
}
