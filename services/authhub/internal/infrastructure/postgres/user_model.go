package postgres

import (
	"time"

	"github.com/google/uuid"
)

type userModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}
