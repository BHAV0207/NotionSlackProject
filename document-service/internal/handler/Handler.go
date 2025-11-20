package handler

import (
	"github.com/BHAV0207/documet-service/internal/websockets"
	"gorm.io/gorm"
)

type Handler struct {
	DB  *gorm.DB
	Hub *websockets.Hub
}

func NewHandler(d *gorm.DB, h *websockets.Hub) *Handler {
	return &Handler{DB: d, Hub: h}
}
