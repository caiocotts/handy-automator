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

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{
		database: db,
	}
}

func (r DeviceRepository) Create(ctx context.Context, ip net.IP) (model.Device, error) {
	id, err := gonanoid.New(12)
	if err != nil {
		return model.Device{}, err
	}
	_, err = r.database.ExecContext(ctx, `insert into "device" values ($1, $2)`, id, ip.String())
	if err != nil {
		return model.Device{}, err
	}
	return model.Device{
		Id: id,
		Ip: ip,
	}, nil
}

func (r DeviceRepository) Get(ctx context.Context, id string) (model.Device, error) {
	d := model.Device{}
	var i string
	err := r.database.QueryRowContext(ctx, `select * from "device" where id = $1`, id).Scan(&d.Id, &i)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Device{}, persistence.ErrNotFound
	}
	if err != nil {
		return model.Device{}, err
	}
	d.Ip = net.ParseIP(i)
	return d, nil
}

func (r DeviceRepository) GetAll(ctx context.Context) ([]model.Device, error) {
	rows, err := r.database.QueryContext(ctx, `select * from "device"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var devices []model.Device
	for rows.Next() {
		var d model.Device
		var i string
		if err := rows.Scan(&d.Id, &i); err != nil {
			return nil, err
		}
		d.Ip = net.ParseIP(i)
		devices = append(devices, d)
	}
	if err = rows.Err(); err != nil {
		return devices, err
	}
	return devices, nil
}

func (r DeviceRepository) Update() (model.Device, error) {
	//TODO implement me
	panic("implement me")
}

func (r DeviceRepository) Delete(ctx context.Context, id string) error {
	res, err := r.database.ExecContext(ctx, `delete from device where id = $1`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return persistence.ErrNotFound
	}
	return nil
}
