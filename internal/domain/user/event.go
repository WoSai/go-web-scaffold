package user

import (
	"time"

	"github.com/wosai/go-web-scaffold/internal/pkg/eventbus"
)

type (
	EventUserCreated struct {
		ID        string
		Name      string
		CreatedAt time.Time
	}
)

const (
	ETUserCreated eventbus.EventType = "user.created"
)

func (e EventUserCreated) Type() eventbus.EventType {
	return ETUserCreated
}
