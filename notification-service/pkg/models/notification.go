package models

type Notification struct {
	ID        string `jsin:"id"`
	UserID    string `json:"user_id"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
