package postgres

import (
	"context"
	"database/sql"
	"decisionMaker/model"
	"decisionMaker/persistence"
	"errors"
	"net"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type DeviceRepository struct {
	database *sql.DB
}

func NewDeviceRepository(db *sql.DB) DeviceRepository {
	return DeviceRepository{
		database: db,
	}
}

func (dr DeviceRepository) Create(ctx context.Context, ip net.IP) (model.Device, error) {
	id, err := gonanoid.New(12)
	if err != nil {
		return model.Device{}, err
	}
	_, err = dr.database.ExecContext(ctx, `insert into "device" values ($1, $2)`, id, ip.String())
	if err != nil {
		return model.Device{}, err
	}
	return model.Device{
		Id: id,
		Ip: ip,
	}, err
}

func (dr DeviceRepository) Get(ctx context.Context, id string) (model.Device, error) {
	d := model.Device{}
	var i string
	err := dr.database.QueryRowContext(ctx, `select * from "device" where id = $1`, id).Scan(&d.Id, &i)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Device{}, persistence.ErrNotFound
	}
	if err != nil {
		return model.Device{}, err
	}
	d.Ip = net.ParseIP(i)
	return d, nil
}

func (dr DeviceRepository) Update() (model.Device, error) {
	//TODO implement me
	panic("implement me")
}

func (dr DeviceRepository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
