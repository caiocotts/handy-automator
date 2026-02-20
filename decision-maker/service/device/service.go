package device

import (
	"context"
	"decisionMaker/model"
	"decisionMaker/persistence"
	"net"
)

type Service struct {
	deviceRepository persistence.DeviceRepository
}

func NewService(r persistence.DeviceRepository) *Service {
	return &Service{deviceRepository: r}
}

func (s Service) Create(ctx context.Context, ip net.IP) (model.Device, error) {
	return s.deviceRepository.Create(ctx, ip)
}

func (s Service) Delete(ctx context.Context, id string) error {
	return s.deviceRepository.Delete(ctx, id)
}

func (s Service) GetById(ctx context.Context, id string) (model.Device, error) {
	return s.deviceRepository.Get(ctx, id)
}

func (s Service) GetAll(ctx context.Context) ([]model.Device, error) {
	return s.deviceRepository.GetAll(ctx)
}
