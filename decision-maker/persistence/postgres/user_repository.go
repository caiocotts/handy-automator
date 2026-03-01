package postgres

import (
	"context"
	"database/sql"
	"decisionMaker/model"
	"decisionMaker/persistence"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	database *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		database: db,
	}
}

func (r UserRepository) Create(ctx context.Context, username, password string) (model.User, error) {
	id, err := gonanoid.New(12)
	if err != nil {
		return model.User{}, err
	}

	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	_, err = r.database.ExecContext(ctx, `insert into "user" values ($1, $2, $3)`, id, username, string(h))
	if err != nil {
		return model.User{}, persistence.ParseDBError(persistence.PostgresError, err)
	}
	return model.User{
		Id:       id,
		Username: username,
	}, err
}

func (r UserRepository) Delete(ctx context.Context, id string) error {
	res, err := r.database.ExecContext(ctx, `delete from "user" where id = $1`, id)
	if err != nil {
		return persistence.ParseDBError(persistence.PostgresError, err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return persistence.ErrNotFound
	}
	return nil
}
