package models

import (
	"github.com/google/uuid"
)

// User struct to describe User object.
type User struct {
	ID           uuid.UUID `db:"id" json:"id" validate:"required,unique"`
	Email        string    `json:"email" validate:"required,email,lte=255"`
	PasswordHash string    `db:"password_hash" json:"password_hash,omitempty" validate:"required,lte=255"`
	AccessKey    string    `db:"password_hash" json:"password_hash,omitempty" validate:"required,lte=255"`
	SecretKey    string    `db:"password_hash" json:"password_hash,omitempty" validate:"required,lte=255"`
	//UserRole     string    `db:"user_role" json:"user_role" validate:"required,lte=25"`
}
