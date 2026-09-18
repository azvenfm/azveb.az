package models
import (
	"time"
	"github.com/google/uuid"
)
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FullName     string    `json:"full_name" db:"full_name"`
	AvatarURL    string    `json:"avatar_url" db:"avatar_url"`
	Role         string    `json:"role" db:"role"`
	Plan         string    `json:"plan" db:"plan"`
	Currency     string    `json:"currency" db:"currency"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
