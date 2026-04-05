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

func (r DeviceRepository) Create(ctx context.Context, hostname string) (model.Device, error) {
	id, err := gonanoid.New(12)
	if err != nil {
		return model.Device{}, err
	}
	_, err = r.database.ExecContext(ctx, `insert into "device" (id, hostname) values ($1, $2)`, id, hostname)
	if err != nil {
		return model.Device{}, persistence.ParseDBError(err)
	}
	return model.Device{
		Id:       id,
		Hostname: hostname,
	}, nil
}

func (r DeviceRepository) Upsert(ctx context.Context, hostname string, ip net.IP) (persistence.UpsertDeviceResult, error) {
	var d model.Device
	id, err := gonanoid.New(12)
	if err != nil {
		return persistence.UpsertDeviceResult{}, err
	}
	var rawIp sql.NullString
	var prevRawIp sql.NullString
	var isNew bool
	err = r.database.QueryRowContext(ctx, `
		with prev as (
			select id, last_known_ip from "device" where hostname = $2
		)
		insert into "device" (id, hostname, last_known_ip)
		values ($1, $2, $3)
		on conflict (hostname) do update set last_known_ip = excluded.last_known_ip
		returning
			id, hostname, last_known_ip, name, type,
			(select id from prev) is null as is_new,
			(select last_known_ip from prev) as prev_ip`,
		id, hostname, ip.String(),
	).Scan(&d.Id, &d.Hostname, &rawIp, &d.Name, &d.Type, &isNew, &prevRawIp)
	if err != nil {
		return persistence.UpsertDeviceResult{}, persistence.ParseDBError(err)
	}
	if rawIp.Valid {
		d.LastKnownIp = net.ParseIP(rawIp.String)
	}
	var prevIP net.IP
	if prevRawIp.Valid {
		prevIP = net.ParseIP(prevRawIp.String)
	}
	return persistence.UpsertDeviceResult{
		Device:     d,
		IsNew:      isNew,
		PreviousIP: prevIP,
	}, nil
}

func (r DeviceRepository) Get(ctx context.Context, id string) (model.Device, error) {
	d := model.Device{}
	var rawIp sql.NullString
	err := r.database.QueryRowContext(ctx, `select id, hostname, last_known_ip, name, type from "device" where id = $1`, id).Scan(&d.Id, &d.Hostname, &rawIp, &d.Name, &d.Type)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Device{}, persistence.ErrNotFound
	}
	if err != nil {
		return model.Device{}, persistence.ParseDBError(err)
	}
	if rawIp.Valid {
		d.LastKnownIp = net.ParseIP(rawIp.String)
	}
	return d, nil
}

func (r DeviceRepository) GetAll(ctx context.Context) ([]model.Device, error) {
	rows, err := r.database.QueryContext(ctx, `select id, hostname, last_known_ip, name, type from "device"`)
	if err != nil {
		return nil, persistence.ParseDBError(err)
	}
	defer rows.Close()
	var devices []model.Device
	for rows.Next() {
		var d model.Device
		var rawIp sql.NullString
		if err := rows.Scan(&d.Id, &d.Hostname, &rawIp, &d.Name, &d.Type); err != nil {
			return nil, persistence.ParseDBError(err)
		}
		if rawIp.Valid {
			d.LastKnownIp = net.ParseIP(rawIp.String)
		}
		devices = append(devices, d)
	}
	if err = rows.Err(); err != nil {
		return devices, persistence.ParseDBError(err)
	}
	return devices, nil
}

func (r DeviceRepository) Delete(ctx context.Context, id string) error {
	res, err := r.database.ExecContext(ctx, `delete from device where id = $1`, id)
	if err != nil {
		return persistence.ParseDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return persistence.ErrNotFound
	}
	return nil
}
