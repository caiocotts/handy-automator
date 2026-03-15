package workflow

import (
	"context"
	"decisionMaker/model"
	"decisionMaker/persistence"
	"errors"
)

var ErrEmptyName = errors.New("workflow name is empty")

type Service struct {
	workflowRepository persistence.WorkflowRepository
}

func NewService(r persistence.WorkflowRepository) *Service {
	return &Service{workflowRepository: r}
}

func (s Service) Create(ctx context.Context, name string) (model.Workflow, error) {
	if name == "" {
		return model.Workflow{}, ErrEmptyName
	}
	uid := ctx.Value("userId").(string)
	if uid == "" {
		return model.Workflow{}, errors.New("error: user ID not present in context")
	}

	w, err := s.workflowRepository.Create(ctx, name, uid)
	if err != nil {
		return model.Workflow{}, err
	}

	return w, nil
}

func (s Service) GetById(ctx context.Context, id string) (model.Workflow, error) {
	w, err := s.workflowRepository.Get(ctx, id)
	if err != nil {
		return model.Workflow{}, err
	}

	return w, nil
}

func (s Service) GetAll(ctx context.Context) ([]model.Workflow, error) {
	workflows, err := s.workflowRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return workflows, nil
}
