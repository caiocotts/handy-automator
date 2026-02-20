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

func NewService(dr persistence.DeviceRepository) *Service {
	return &Service{deviceRepository: dr}
}

func (ds Service) Create(ctx context.Context, ip net.IP) (model.Device, error) {
	return ds.deviceRepository.Create(ctx, ip)
}

func (ds Service) Delete(ctx context.Context, id string) error {
	return ds.deviceRepository.Delete(ctx, id)
}

func (ds Service) GetById(ctx context.Context, id string) (model.Device, error) {
	return ds.deviceRepository.Get(ctx, id)
}

func (ds Service) GetAll(ctx context.Context) ([]model.Device, error) {
	return ds.deviceRepository.GetAll(ctx)
}
