package persistence

import (
	"context"
	"decisionMaker/model"
)

type WorkflowRepository interface {
	Create(ctx context.Context, name, userId string) (model.Workflow, error)
	Get(ctx context.Context, id string) (model.Workflow, error)
	GetByUserIdAndGestureId(ctx context.Context, userId string, gestureId int) (model.Workflow, error)
	GetAll(ctx context.Context) ([]model.Workflow, error)
	Update() (model.Workflow, error)
	UpdateState(ctx context.Context, id, state string) error
	Delete(ctx context.Context, id string) error
	AssociateDevices(ctx context.Context, workflowId string, deviceIds []string) ([]string, error)
}
