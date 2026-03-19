package postgres

import (
	"context"
	"database/sql"
	"decisionMaker/model"
	"decisionMaker/persistence"
	"encoding/json"

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
		return model.User{}, persistence.ParseDBError(err)
	}

	return model.User{
		Id:       id,
		Username: username,
	}, nil
}

func (r UserRepository) Get(ctx context.Context, id string) (model.User, error) {
	u := model.User{}
	var rawEmbedding []byte
	err := r.database.
		QueryRowContext(ctx, `select id, username, hash, refresh_token, face_embedding from "user" where id = $1`, id).
		Scan(&u.Id, &u.Username, &u.PasswordHash, &u.RefreshToken, &rawEmbedding)
	if err != nil {
		return model.User{}, persistence.ParseDBError(err)
	}
	if rawEmbedding != nil {
		var v []float64
		if err := json.Unmarshal(rawEmbedding, &v); err != nil {
			return model.User{}, persistence.ParseDBError(err)
		}
		u.FaceEmbedding = &v
	}

	return u, nil
}

func (r UserRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	u := model.User{}
	err := r.database.
		QueryRowContext(ctx, `select id, username, hash, refresh_token, face_embedding from "user" where username = $1`, username).
		Scan(&u.Id, &u.Username, &u.PasswordHash, &u.RefreshToken, &u.FaceEmbedding)

	if err != nil {
		return model.User{}, persistence.ParseDBError(err)
	}

	return u, nil
}

func (r UserRepository) UpdateRefreshToken(ctx context.Context, id, refreshToken string) error {
	query := `
update "user"
set refresh_token = $1
where id = $2;
`
	_, err := r.database.ExecContext(ctx, query, refreshToken, id)
	if err != nil {
		return persistence.ParseDBError(err)
	}
	return nil
}

func (r UserRepository) Delete(ctx context.Context, id string) error { // TODO this should cascade delete all records whose FK is this user id
	res, err := r.database.ExecContext(ctx, `delete from "user" where id = $1`, id)
	if err != nil {
		return persistence.ParseDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return persistence.ErrNotFound
	}
	return nil
}
