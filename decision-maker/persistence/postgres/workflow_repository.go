package postgres

import (
	"context"
	"database/sql"
	"decisionMaker/model"

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

func (r WorkflowRepository) Create(ctx context.Context, name string) (model.Workflow, error) {
	id, err := gonanoid.New(12)
	if err != nil {
		return model.Workflow{}, err
	}
	_, err = r.database.ExecContext(ctx, `insert into "workflow" values ($1, $2)`, id, name)
	if err != nil {
		return model.Workflow{}, err
	}
	return model.Workflow{
		Id:   id,
		Name: name,
	}, nil

}

func (r WorkflowRepository) Get(ctx context.Context, id string) (model.Workflow, error) {
	//TODO implement me
	panic("implement me")
}

func (r WorkflowRepository) GetAll(ctx context.Context) ([]model.Workflow, error) {
	//TODO implement me
	panic("implement me")
}

func (r WorkflowRepository) Update() (model.Workflow, error) {
	//TODO implement me
	panic("implement me")
}

func (r WorkflowRepository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
