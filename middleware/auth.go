package middleware

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/google/uuid"
	"time"
)

// Auth middleware for get use session
func Auth(rs *redis.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess := c.Locals("session").(*session.Session)
		sessionToken, ok := sess.Get("session_token").(string)
		if !ok {
			sessionToken := uuid.New().String()
			sess.Set("session_token", sessionToken)
			return c.Next()
		}
		u := getSessionUser(rs, sessionToken)
		if u == nil {
			return c.Next()
		}
		c.Locals("user", u)
		return c.Next()
	}
}

func getSessionUser(rs *redis.Storage, sessionToken string) *User {
	b, err := rs.Get("user--" + sessionToken)
	if err != nil {
		return nil
	}
	var user User
	err = json.Unmarshal(b, &user)
	if err != nil {
		return nil
	}
	return &user
}

// User to use in auth middleware
type User struct {
	ID        int64      `json:"id"`
	UUID      uuid.UUID  `json:"uuid"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
