package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BHAV0207/user-service/internal/service"
	"github.com/BHAV0207/user-service/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	DB *pgxpool.Pool
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var User models.User

	if err := json.NewDecoder(r.Body).Decode(&User); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//  verify if user alread exist via email

	user := service.FindUserEmail(ctx, h.DB, User.ID)
	if user {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	//  hash password
	hashedPassword, err := service.HashPassword(User.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	User.Password = hashedPassword

	//  store user in db
	query := `INSERT INTO users (name, email, phone, password, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = h.DB.Exec(ctx, query, User.Name, User.Email, User.Phone, User.Password, User.CreatedAt, User.UpdatedAt)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})

}
