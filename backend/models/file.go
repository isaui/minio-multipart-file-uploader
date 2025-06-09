package models

import (
	"time"

	"gorm.io/gorm"
)

// File represents a file uploaded to the system
type File struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Filename   string         `json:"filename" gorm:"not null"`
	Filepath   string         `json:"filepath" gorm:"not null"`
	Size       int64          `json:"size" gorm:"not null"`
	UploadedAt time.Time      `json:"uploaded_at" gorm:"autoCreateTime"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
