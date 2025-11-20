package models

import "time"

type Snapshot struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	DocumentID string    `gorm:"type:uuid;index" json:"document_id"`
	Data       []byte    `gorm:"type:bytea" json:"-"`
	CreatedAt  time.Time `json:"created_at"`
}
