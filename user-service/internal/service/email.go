package service

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)


func FindUserEmail(ctx context.Context, db *pgxpool.Pool, email string) bool {
	query := `SELECT COUNT(1) FROM users WHERE email = $1`
	var count int
	err := db.QueryRow(ctx, query, email).Scan(&count)
	if err != nil {
		return false // safer fallback if query fails
	}
	return count > 0
}
