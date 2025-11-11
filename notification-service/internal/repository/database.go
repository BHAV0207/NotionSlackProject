package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDb(Uri string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, Uri)
	if err != nil {
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("❌ Unable to connect to Postgres: %v", err)
	}

	fmt.Println("✅ Connected to PostgreSQL successfully!")
	DB = pool
	return pool

}
