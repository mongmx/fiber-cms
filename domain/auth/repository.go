package auth

import (
	"encoding/json"
	"github.com/gofiber/storage/redis"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository - auth store APIs.
type Repository interface {
	createUser(user *User) (*User, error)
	createAuth(auth *Auth) (*Auth, error)
	findUserByEmail(email string) (*User, error)
	storeSessionUser(token string, user *User) error
	getSessionUser(token string) (User, error)
}

type repo struct {
	db *sqlx.DB
	rs *redis.Storage
}

// NewRepository is a factory function of auth store.
func NewRepository(db *sqlx.DB, rs *redis.Storage) Repository {
	return &repo{
		db: db,
		rs: rs,
	}
}

func (r repo) createUser(user *User) (*User, error) {
	query := `
		INSERT INTO users (email) VALUES ($1) RETURNING id
	`
	var lastInsertedID int64
	err := r.db.QueryRowx(query, user.Email).Scan(&lastInsertedID)
	if err != nil {
		return nil, err
	}
	user.ID = lastInsertedID
	return user, nil
}

func (r repo) createAuth(auth *Auth) (*Auth, error) {
	query := `
		INSERT INTO auths (user_id, type, secret) VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, auth.UserID, auth.Type, auth.Secret)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (r repo) findUserByEmail(email string) (*User, error) {
	var queryUser struct {
		ID     int64  `db:"id"`
		UUID   string `db:"uuid"`
		Email  string `db:"email"`
		Secret string `db:"secret"`
	}
	query := `
		SELECT u.id, u.uuid, u.email, a.secret FROM users u JOIN auths a on u.id = a.user_id WHERE email = $1
	`
	err := r.db.QueryRowx(query, email).StructScan(&queryUser)
	uid, _ := uuid.Parse(queryUser.UUID)
	u := &User{
		model: model{
			ID:   queryUser.ID,
			UUID: uid,
		},
		Email: queryUser.Email,
		Auth: &Auth{
			UserID: queryUser.ID,
			Type:   "email",
			Secret: queryUser.Secret,
		},
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r repo) storeSessionUser(token string, user *User) error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = r.rs.Set("user--"+token, b, 0)
	if err != nil {
		return err
	}
	return nil
}

func (r repo) getSessionUser(token string) (User, error) {
	b, err := r.rs.Get("user--" + token)
	if err != nil {
		return User{}, err
	}
	var user User
	if err := json.Unmarshal(b, &user); err != nil {
		return User{}, err
	}
	return user, nil
}
