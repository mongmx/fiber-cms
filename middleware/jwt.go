package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte("SECRET"),
		ErrorHandler: jwtError,
		ContextKey:   "jwt",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

//func validToken(t *jwt.Token, id string) bool {
//	n, err := strconv.Atoi(id)
//	if err != nil {
//		return false
//	}
//
//	claims := t.Claims.(jwt.MapClaims)
//	uid := int(claims["user_id"].(float64))
//
//	if uid != n {
//		return false
//	}
//
//	return true
//}
//
//func validUser(id string, p string) bool {
//	db := database.DB
//	var user model.User
//	db.First(&user, id)
//	if user.Username == "" {
//		return false
//	}
//	if !CheckPasswordHash(p, user.Password) {
//		return false
//	}
//	return true
//}
