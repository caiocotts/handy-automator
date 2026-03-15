package user

import (
	"context"
	"decisionMaker/model"
	"decisionMaker/persistence"
	"errors"
)

var ErrPasswordTooLong = errors.New("password must not be longer than 72 characters")

type Service struct {
	userRepository persistence.UserRepository
}

func NewService(r persistence.UserRepository) *Service {
	return &Service{
		userRepository: r,
	}
}

func (s Service) Register(ctx context.Context, username, password string) (model.User, error) {
	if len(password) > 72 {
		return model.User{}, ErrPasswordTooLong
	}
	u, err := s.userRepository.Create(ctx, username, password)
	if err != nil {
		return model.User{}, err
	}

	return u, nil
}

func (s Service) Delete(ctx context.Context, id string) error {
	return s.userRepository.Delete(ctx, id)
}
