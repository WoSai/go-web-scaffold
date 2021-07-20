package user

import (
	"github.com/wosai/go-web-scaffold/internal/application/user/command"
	"github.com/wosai/go-web-scaffold/internal/application/user/handler"
	"github.com/wosai/go-web-scaffold/internal/domain/user"
	"github.com/wosai/go-web-scaffold/internal/pkg/eventbus"
)

type (
	Application struct {
		Commands commands
		Queries  query
	}

	commands struct {
		CreateUser *command.CreateUserHandler
	}

	query struct{}
)

func BuildApplication(repo user.Repository) *Application {
	app := &Application{
		Commands: commands{
			CreateUser: command.NewCreateUserHandler(repo),
		},
		Queries: query{},
	}

	// subscribe event
	eventbus.Subscribe(user.ETUserCreated, handler.UserCreatedSubscriber{})
	return app
}
