package users

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"go.starlark.net/lib/time"
)

type User struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Login      string    `json:"login"`
	Password   string    `json:"password`
	CreatedAt  time.Time `json:"create_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    bool      `json:"-"`
	LastLogin  time.Time `json:"last_login"`
}

var (
	ErrPasswordRequired = errors.New("password can't be blank")
	ErrPasswordLen      = errors.New("password must have at least 6 characters")
	ErrNameRequired     = errors.New("Name is required")
	ErrLoginRequired    = errors.New("Login is required")
)

func New(name, login, password string) (*User, error) {
	user := &User{
		Name:       name,
		Login:      login,
		ModifiedAt: time.Now(),
		Deleted:    false,
	}

	err := user.SetPassword(password)
	if err != nil {
		return nil, err
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) SetPassword(password string) error {
	if password == "" {
		return ErrPasswordRequired
	}

	if len(password) < 6 {
		return ErrPasswordLen
	}

	u.Password = fmt.Sprintf("%x", (md5.Sum([]byte(password))))

	return nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrNameRequired
	}

	if u.Login == "" {
		return ErrLoginRequired
	}

	if u.Password == fmt.Sprintf("%x", md5.Sum([]byte(""))) {
		return ErrPasswordRequired
	}

	return nil
}
