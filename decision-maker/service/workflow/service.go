package workflow

import (
	"context"
	"decisionMaker/model"
	"decisionMaker/persistence"
	"decisionMaker/service/discovery"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

var ErrEmptyName = errors.New("workflow name is empty")
var ErrDeviceUnreachable = errors.New("device not found on the network")

type Service struct {
	workflowRepository persistence.WorkflowRepository
	discoveryService   *discovery.Service
}

func NewService(r persistence.WorkflowRepository, d *discovery.Service) *Service {
	return &Service{workflowRepository: r, discoveryService: d}
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

func (s Service) Trigger(ctx context.Context, gestureId int) ([]model.DeviceTriggerStatus, error) {
	uid := ctx.Value("userId").(string)
	if uid == "" {
		return nil, errors.New("error: user ID not present in context")
	}

	w, err := s.workflowRepository.GetByUserIdAndGestureId(ctx, uid, gestureId)
	if err != nil {
		return nil, err
	}

	targetState := "on"
	if w.State == "on" {
		targetState = "off"
	}

	results := make([]model.DeviceTriggerStatus, len(w.Devices))
	var wg sync.WaitGroup
	for i, device := range w.Devices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			status := model.DeviceTriggerStatus{DeviceId: device.Id, Ok: true}

			ip, ok := s.discoveryService.Resolve(device.Hostname)
			if !ok {
				if device.LastKnownIp == nil {
					status.Ok = false
					status.Error = fmt.Sprintf("%s: %s", ErrDeviceUnreachable.Error(), device.Hostname)
					results[i] = status
					return
				}
				ip = device.LastKnownIp
			}

			stateURL := fmt.Sprintf("http://%s/device/state", ip.String())
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, stateURL, nil)
			if err != nil {
				status.Ok = false
				status.Error = fmt.Sprintf("error creating state request for device %s: %s", device.Id, err)
				results[i] = status
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				status.Ok = false
				status.Error = fmt.Sprintf("error checking state of device %s: %s", device.Id, err)
				results[i] = status
				return
			}
			var stateResp struct {
				State string `json:"state"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&stateResp); err != nil {
				resp.Body.Close()
				status.Ok = false
				status.Error = fmt.Sprintf("error reading state of device %s: %s", device.Id, err)
				results[i] = status
				return
			}
			resp.Body.Close()

			if stateResp.State == targetState {
				results[i] = status
				return
			}

			toggleURL := fmt.Sprintf("http://%s/device/toggle", ip.String())
			req, err = http.NewRequestWithContext(ctx, http.MethodPost, toggleURL, nil)
			if err != nil {
				status.Ok = false
				status.Error = fmt.Sprintf("error creating toggle request for device %s: %s", device.Id, err)
				results[i] = status
				return
			}
			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				status.Ok = false
				status.Error = fmt.Sprintf("error toggling device %s: %s", device.Id, err)
				results[i] = status
				return
			}
			resp.Body.Close()
			results[i] = status
		}()
	}
	wg.Wait()

	if err := s.workflowRepository.UpdateState(ctx, w.Id, targetState); err != nil {
		return results, err
	}

	return results, nil
}
