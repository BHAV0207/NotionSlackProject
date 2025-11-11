package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserResponse struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone int    `json:"phone"`
}

func GetUserByID(userID string) (*UserResponse, error) {
	url := fmt.Sprintf("http://user-service:8000/me/%s", userID)
	// Make HTTP GET request to user service (omitted for brevity)
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned status: %d", resp.StatusCode)
	}

	var userResp UserResponse
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	if err != nil {
		return nil, err
	}

	return &userResp, nil

}
