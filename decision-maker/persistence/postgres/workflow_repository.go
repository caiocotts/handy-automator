package postgres

import (
	"context"
	"database/sql"
	"decisionMaker/model"
	"decisionMaker/persistence"
	"errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type WorkflowRepository struct {
	database *sql.DB
}

func NewWorkflowRepository(db *sql.DB) *WorkflowRepository {
	return &WorkflowRepository{
		database: db,
	}
}

func (r WorkflowRepository) Create(ctx context.Context, name, userId string) (model.Workflow, error) {
	id, err := gonanoid.New(12)
	if err != nil {
		return model.Workflow{}, err
	}
	_, err = r.database.ExecContext(ctx, `insert into "workflow" values ($1, $2, $3)`, id, name, userId)
	if err != nil {
		return model.Workflow{}, err
	}
	return model.Workflow{
		Id:     id,
		Name:   name,
		UserId: userId,
	}, nil

}

func (r WorkflowRepository) Get(ctx context.Context, id string) (model.Workflow, error) {
	w := model.Workflow{}
	err := r.database.QueryRowContext(ctx, `select id, name from "workflow" where id = $1`, id).Scan(&w.Id, &w.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Workflow{}, persistence.ErrNotFound
	}
	if err != nil {
		return model.Workflow{}, err
	}
	return w, nil
}

func (r WorkflowRepository) GetAll(ctx context.Context) ([]model.Workflow, error) {
	rows, err := r.database.QueryContext(ctx, `select * from "workflow"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var workflows []model.Workflow
	for rows.Next() {
		var w model.Workflow
		if err := rows.Scan(&w.Id, &w.Name, &w.UserId); err != nil {
			return nil, err
		}
		workflows = append(workflows, w)
	}
	if err = rows.Err(); err != nil {
		return workflows, err
	}
	return workflows, nil
}

func (r WorkflowRepository) Update() (model.Workflow, error) {
	//TODO implement me
	panic("implement me")
}

func (r WorkflowRepository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
