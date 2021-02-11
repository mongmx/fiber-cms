package auth

import "golang.org/x/crypto/bcrypt"

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

func (c useCase) register(user *User) error {
	user, err := c.repo.createUser(user)
	if err != nil {
		return err
	}
	user.Auth.UserID = user.ID
	_, err = c.repo.createAuth(user.Auth)
	if err != nil {
		return err
	}
	return nil
}

func (c useCase) login(cred Credentials, sessionToken string) error {
	user, err := c.repo.findUserByEmail(cred.Email)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Auth.Secret), []byte(cred.Password))
	if err != nil {
		return err
	}
	err = c.repo.storeSessionUser(sessionToken, user)
	if err != nil {
		return err
	}
	return nil
}

func (c useCase) profile(sessionToken string) (User, error) {
	user, err := c.repo.getSessionUser(sessionToken)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
