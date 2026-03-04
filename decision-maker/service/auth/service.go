package auth

import (
	"context"
	"decisionMaker/config"
	"decisionMaker/consts"
	"decisionMaker/model"
	"decisionMaker/persistence"
	"errors"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("authentication failed")

type Service struct {
	userRepository persistence.UserRepository
}

func NewService(r persistence.UserRepository) *Service {
	return &Service{
		userRepository: r,
	}
}

func (s Service) Login(ctx context.Context, username, password string) (model.User, string, error) {
	u, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return model.User{}, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return model.User{}, "", ErrInvalidCredentials
	}
	t, err := jwt.NewBuilder().
		Subject(u.Id).
		Issuer(consts.JWTIssuer).
		IssuedAt(time.Now()).
		Expiration(time.Now().AddDate(0, 1, 0)).
		Build()
	if err != nil {
		return model.User{}, "", err
	}
	refreshToken, err := jwt.Sign(t, jwt.WithKey(jwa.HS256(), config.JWTSecret))
	if err != nil {
		fmt.Printf("failed to sign token: %s\n", err)
		return model.User{}, "", err
	}
	err = s.userRepository.UpdateRefreshToken(ctx, u.Id, string(refreshToken))
	if err != nil {
		return model.User{}, "", err
	}

	t, err = jwt.NewBuilder().
		Subject(u.Id).
		Issuer(consts.JWTIssuer).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(time.Hour * 2)).
		Build()
	accessToken, err := jwt.Sign(t, jwt.WithKey(jwa.HS256(), config.JWTSecret))
	if err != nil {
		fmt.Printf("failed to sign token: %s\n", err)
		return model.User{}, "", err
	}
	rt := string(refreshToken)
	u.RefreshToken = &rt

	return u, string(accessToken), nil
}

func (s Service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	return "", nil
}

func (s Service) ValidateAccessToken(token string) (string, error) {
	verifiedToken, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256(), config.JWTSecret))
	if err != nil {
		return "", err
	}

	uid, exists := verifiedToken.Subject()
	if !exists {
		return "", jwt.ValidateError()
	}

	return uid, nil
}
