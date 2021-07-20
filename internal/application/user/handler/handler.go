package handler

import (
	"context"
	"fmt"

	domain "github.com/wosai/go-web-scaffold/internal/domain/user"
	"github.com/wosai/go-web-scaffold/internal/pkg/eventbus"
)

type (
	UserCreatedSubscriber struct{}
)

func (uc UserCreatedSubscriber) Handle(ctx context.Context, event eventbus.Event) {
	ev := event.(domain.EventUserCreated)
	fmt.Println("create a new user", ev.Name)
}
