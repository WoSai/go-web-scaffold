package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wosai/go-web-scaffold/internal/pkg/eventbus"
	"github.com/wosai/go-web-scaffold/internal/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID             string           `validate:"required"`
		Username       string           `validate:"required"`
		HashedPassword []byte           `validate:"required"`
		Events         *eventbus.Events `validate:"required"`
	}
)

func NewUser(username, password string) (*User, error) {
	id := uuid.New().String()

	eb := eventbus.NewEvents()
	eb.Add(&EventUserCreated{
		ID:        id,
		Name:      username,
		CreatedAt: time.Now(),
	})
	entity := &User{ID: id, Username: username, Events: eb}
	if err := entity.setPassword(password); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *User) setPassword(pwd string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = hashed
	return nil
}

func (u *User) Validate() error {
	return validator.Struct(u)
}

func (u *User) CheckPassword(pwd string) bool {
	if err := bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(pwd)); err != nil {
		return false
	}
	return true
}

func (u *User) ChangePassword(n, o string) error {
	if !u.CheckPassword(o) {
		return errors.New("旧密码校验失败")
	}
	return u.setPassword(n)
}
