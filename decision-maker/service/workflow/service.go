package workflow

import (
	"context"
	"decisionMaker/model"
	"decisionMaker/persistence"
	"errors"
	"fmt"
	"net/http"
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

func (s Service) AssociateDevices(ctx context.Context, workflowId string, deviceIds []string) ([]string, error) {
	ids, err := s.workflowRepository.AssociateDevices(ctx, workflowId, deviceIds)
	if err != nil {
		return nil, err
	}

	return ids, err
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

func (s Service) Trigger(ctx context.Context, gestureId int) error {
	uid := ctx.Value("userId").(string)
	if uid == "" {
		return errors.New("error: user ID not present in context")
	}

	w, err := s.workflowRepository.GetByUserIdAndGestureId(ctx, uid, gestureId)
	if err != nil {
		return err
	}

	for _, device := range w.Devices {
		url := fmt.Sprintf("http://%s/device/toggle", device.Ip.String())
		resp, err := http.Post(url, "", nil)
		if err != nil {
			return fmt.Errorf("error toggling device %s: %w", device.Id, err)
		}
		resp.Body.Close()
	}

	return nil
}
