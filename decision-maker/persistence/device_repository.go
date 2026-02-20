package persistence

import (
	"context"
	"decisionMaker/model"
	"net"
)

type DeviceRepository interface {
	Create(ctx context.Context, ip net.IP) (model.Device, error)
	Get(ctx context.Context, id string) (model.Device, error)
	GetAll(ctx context.Context) ([]model.Device, error)
	Update() (model.Device, error)
	Delete(ctx context.Context, id string) error
}
