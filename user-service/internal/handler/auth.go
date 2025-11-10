package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BHAV0207/user-service/internal/service"
	"github.com/BHAV0207/user-service/pkg/models"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	DB *pgxpool.Pool
}

var jwtSecret = []byte("supersecretkey")

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var User models.User

	if err := json.NewDecoder(r.Body).Decode(&User); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//  verify if user alread exist via email

	userExists := service.FindUserEmail(ctx, h.DB, User.Email)
	if userExists {
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

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Implementation for user login
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//  parse request body
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	exists := service.FindUserEmail(ctx, h.DB, req.Email)
	if !exists {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	var HashedPassword string
	query := `SELECT password FROM users WHERE email = $1`
	err := h.DB.QueryRow(ctx, query, req.Email).Scan(&HashedPassword)
	if err != nil {
		http.Error(w, "Error fetching user data", http.StatusInternalServerError)
		return
	}

	isTrue := service.CheckPasswordHash(req.Password, HashedPassword)
	if !isTrue {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	//jwt login karna hai mittar aab

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": req.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   tokenString,
	})

}
