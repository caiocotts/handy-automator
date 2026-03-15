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
	if errors.Is(err, persistence.ErrNotFound) {
		return model.User{}, "", ErrInvalidCredentials
	}
	if err != nil {
		return model.User{}, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return model.User{}, "", ErrInvalidCredentials
	}

	t, err := createJWT(u.Id, time.Now().AddDate(0, 1, 0))
	if err != nil {
		return model.User{}, "", err
	}

	refreshToken, err := signJWT(t)
	if err != nil {
		return model.User{}, "", err
	}

	err = s.userRepository.UpdateRefreshToken(ctx, u.Id, refreshToken)
	if err != nil {
		return model.User{}, "", err
	}

	t, err = createJWT(u.Id, time.Now().Add(time.Hour*2))
	if err != nil {
		return model.User{}, "", err
	}

	accessToken, err := signJWT(t)
	if err != nil {
		return model.User{}, "", err
	}

	u.RefreshToken = &refreshToken

	return u, accessToken, nil
}

func (s Service) Refresh(ctx context.Context) (string, error) {
	uid := ctx.Value("userId").(string)
	if uid == "" {
		return "", errors.New("error: user ID not present in context")
	}

	t, err := createJWT(uid, time.Now().Add(time.Hour*2))
	if err != nil {
		return "", err
	}

	accessToken, err := signJWT(t)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s Service) ValidateToken(token string) (string, error) {
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

func signJWT(t jwt.Token) (string, error) {
	signedToken, err := jwt.Sign(t, jwt.WithKey(jwa.HS256(), config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("error: failed to sign JWT token: %w", err)
	}

	return string(signedToken), nil
}

func createJWT(uid string, expiration time.Time) (jwt.Token, error) {
	t, err := jwt.NewBuilder().
		Subject(uid).
		Issuer(consts.JWTIssuer).
		IssuedAt(time.Now()).
		Expiration(expiration).
		Build()

	if err != nil {
		return nil, fmt.Errorf("error: failed to build JWT token: %w", err)
	}

	return t, nil
}
