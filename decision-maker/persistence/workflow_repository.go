package persistence

import (
	"context"
	"decisionMaker/model"
)

type WorkflowRepository interface {
	Create(ctx context.Context, name string) (model.Workflow, error)
	Get(ctx context.Context, id string) (model.Workflow, error)
	GetAll(ctx context.Context) ([]model.Workflow, error)
	Update() (model.Workflow, error)
	Delete(ctx context.Context, id string) error
}
