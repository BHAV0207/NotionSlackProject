package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	Kafka "github.com/BHAV0207/user-service/internal/kafka")

type UserHandler struct {
	DB       *pgxpool.Pool
	Producer *Kafka.Producer
}
