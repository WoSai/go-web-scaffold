package user

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/wosai/go-web-scaffold/internal/application/user/command"
	"github.com/wosai/go-web-scaffold/internal/application/user/handler"
	"github.com/wosai/go-web-scaffold/internal/domain/user"
	"github.com/wosai/go-web-scaffold/internal/pkg/eventbus"
)

type (
	Application struct {
		Commands   commands
		Queries    query
		CreateUser endpoint.Endpoint
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
		Queries:    query{},
		CreateUser: command.MakeCreateUserEndpoint(repo),
	}

	// subscribe event
	eventbus.Subscribe(user.ETUserCreated, handler.UserCreatedSubscriber{})
	return app
}
