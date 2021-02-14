package auth

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// UseCase - auth usecase APIs.
type UseCase interface {
	register(user *User) error
	login(cred Credentials, sessionToken string) error
	profile(sessionToken string) (User, error)
}

type useCase struct {
	repo Repository
}

// NewUseCase is a factory function of auth usecase.
func NewUseCase(repo Repository) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (u useCase) register(user *User) error {
	user, err := u.repo.createUser(user)
	if err != nil {
		return err
	}
	user.Auth.UserID = user.ID
	_, err = u.repo.createAuth(user.Auth)
	if err != nil {
		return err
	}
	return nil
}

func (u useCase) login(cred Credentials, sessionToken string) error {
	user, err := u.repo.findUserByEmail(cred.Email)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Auth.Secret), []byte(cred.Password))
	if err != nil {
		return err
	}
	err = u.repo.storeSessionUser(sessionToken, user)
	if err != nil {
		return err
	}
	return nil
}

func (u useCase) profile(sessionToken string) (User, error) {
	user, err := u.repo.getSessionUser(sessionToken)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u useCase) validToken(c *fiber.Ctx) bool {
	token := c.Locals("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	sessionToken := claims["session_id"].(string)
	user, _ := u.repo.getSessionUser(sessionToken)
	if user.ID <= 0 {
		return false
	}
	return true
}
