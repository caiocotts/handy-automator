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

func (ds Service) RegisterDevice(ctx context.Context, ip net.IP) (model.Device, error) {
	d, err := ds.deviceRepository.Create(ctx, ip)
	if err != nil {
		return model.Device{}, err
	}
	return d, nil
}

func (ds Service) GetDeviceById(ctx context.Context, id string) (model.Device, error) {
	d, err := ds.deviceRepository.Get(ctx, id)
	if err != nil {
		return model.Device{}, err
	}
	return d, nil
}
