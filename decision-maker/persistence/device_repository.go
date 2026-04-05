package persistence

import (
	"context"
	"decisionMaker/model"
	"net"
)

type UpsertDeviceResult struct {
	Device     model.Device
	IsNew      bool
	PreviousIP net.IP // nil if new device or IP was previously unknown
}

type DeviceRepository interface {
	Create(ctx context.Context, hostname string) (model.Device, error)
	Upsert(ctx context.Context, hostname string, ip net.IP) (UpsertDeviceResult, error)
	Get(ctx context.Context, id string) (model.Device, error)
	GetAll(ctx context.Context) ([]model.Device, error)
	Delete(ctx context.Context, id string) error
}
