package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(name string) *User {
	return &User{
		ID:		uuid.NewString(),
		Name:	name,
		CreatedAt: time.Now(),
	}
}

func (u *User) Validate() error {
	if u.ID == "" {
		return errors.New("user must have an ID")
	}
	if u.Name == "" {
		return errors.New("user must have a name")
	}
	return nil
}