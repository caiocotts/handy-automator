package workflow

import (
	"context"
	"decisionMaker/model"
	"decisionMaker/persistence"
)

type Service struct {
	workflowRepository persistence.WorkflowRepository
}

func NewService(r persistence.WorkflowRepository) *Service {
	return &Service{workflowRepository: r}
}

func (s Service) Create(ctx context.Context, name string) (model.Workflow, error) {
	return s.workflowRepository.Create(ctx, name)
}

func (s Service) GetById(ctx context.Context, id string) (model.Workflow, error) {
	return s.workflowRepository.Get(ctx, id)
}
