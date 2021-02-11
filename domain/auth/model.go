package auth

import (
	"github.com/google/uuid"
	"time"
)

type model struct {
	ID        int64
	UUID      uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type User struct {
	model
	Email string `json:"email"`
	Auth  *Auth  `json:"-"`
}

type Auth struct {
	UserID int64
	Type   string
	Secret string
}

type Topic struct {
	model
	UserID int64  `json:"-"`
	Title  string `json:"title"`
}

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
