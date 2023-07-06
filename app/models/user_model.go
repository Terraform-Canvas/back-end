package models

// User struct to describe User object.
type User struct {
	ID           int    `json:"id" validate:"required,unique"`
	Email        string `json:"email" validate:"required,email,lte=255"`
	PasswordHash string `json:"password_hash" validate:"required,lte=255"`
	AccessKey    string `json:"accessKey" validate:"required,lte=255"`
	SecretKey    string `json:"secretKey" validate:"required,lte=255"`
}
