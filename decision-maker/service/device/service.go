package device

import (
	"context"
	"decisionMaker/model"
	"decisionMaker/persistence"
	"errors"
	"net"
)

var ErrInvalidIPFormat = errors.New("malformed IP address")

type Service struct {
	deviceRepository persistence.DeviceRepository
}

func NewService(r persistence.DeviceRepository) *Service {
	return &Service{deviceRepository: r}
}

func (s Service) Create(ctx context.Context, ip string) (model.Device, error) {
	i := net.ParseIP(ip)
	if i == nil {
		return model.Device{}, ErrInvalidIPFormat
	}

	d, err := s.deviceRepository.Create(ctx, i)
	if err != nil {
		return model.Device{}, err
	}

	return d, nil
}

func (s Service) Delete(ctx context.Context, id string) error {
	return s.deviceRepository.Delete(ctx, id)
}

func (s Service) GetById(ctx context.Context, id string) (model.Device, error) {
	d, err := s.deviceRepository.Get(ctx, id)
	if err != nil {
		return model.Device{}, err
	}

	return d, nil
}

func (s Service) GetAll(ctx context.Context) ([]model.Device, error) {
	devices, err := s.deviceRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return devices, nil
}
