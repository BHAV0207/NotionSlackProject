package models

import "time"

type Document struct {
	ID        string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title     string     `gorm:"type:text" json:"title"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Snapshots []*Snapshot `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}
