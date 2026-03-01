package persistence

import (
	"context"
	"decisionMaker/model"
)

type UserRepository interface {
	Create(ctx context.Context, username, password string) (model.User, error)
	Delete(ctx context.Context, id string) error
}
