package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"log"
)

// Session middleware for start to use session
func Session(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := sess.Save()
			if err != nil {
				log.Println(err)
			}
		}()
		c.Locals("session", sess)
		return c.Next()
	}
}
