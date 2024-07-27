package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DBpool *pgxpool.Pool

func InitDB() {
	var err error
	databaseUrl := "postgres://tac:test@localhost:5433/gps"
	DBpool, err = pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
}
