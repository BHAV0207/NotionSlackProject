package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/BHAV0207/user-service/internal/middleware"
)

func (h *UserHandler) UserInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Extract email from context
	email, ok := r.Context().Value(middleware.UserEmailKey).(string)
	if !ok || email == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	query := `SELECT name , email , phone FROM users where email = $1`
	err := h.DB.QueryRow(ctx, query, email).Scan(&user.Name, &user.Email, &user.Phone)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}
