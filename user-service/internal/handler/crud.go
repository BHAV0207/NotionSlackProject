package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/BHAV0207/user-service/internal/middleware"
	"github.com/BHAV0207/user-service/pkg/models"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
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

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Create a context tied to the HTTP request (cancels if request ends)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Extract ID from URL params
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		http.Error(w, "Missing user ID in request", http.StatusBadRequest)
		return
	}

	// Query user from DB
	var user models.User
	query := `SELECT name, email, phone FROM users WHERE id = $1`
	err := h.DB.QueryRow(ctx, query, userID).Scan(&user.Name, &user.Email, &user.Phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // ✅ no record found
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
