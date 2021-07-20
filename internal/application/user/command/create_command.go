package command

import (
	"context"

	"github.com/wosai/go-web-scaffold/internal/domain/user"
)

type (
	CreateUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	CreateUserResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	CreateUserHandler struct {
		repo user.Repository
	}
)

func NewCreateUserHandler(repo user.Repository) *CreateUserHandler {
	return &CreateUserHandler{repo: repo}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd *CreateUserRequest) (*CreateUserResponse, error) {
	entity, err := user.NewUser(cmd.Username, cmd.Password)
	if err != nil {
		return nil, err
	}
	if err = entity.Validate(); err != nil {
		return nil, err
	}

	if err = h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return &CreateUserResponse{ID: entity.ID, Username: entity.Username}, nil
}
