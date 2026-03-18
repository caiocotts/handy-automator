package postgres

import (
	"context"
	"database/sql"
	"decisionMaker/model"
	"decisionMaker/persistence"
	"net"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type WorkflowRepository struct {
	database *sql.DB
}

func NewWorkflowRepository(db *sql.DB) *WorkflowRepository {
	return &WorkflowRepository{
		database: db,
	}
}

func (r WorkflowRepository) Create(ctx context.Context, name, userId string) (model.Workflow, error) {
	id, err := gonanoid.New(12)
	if err != nil {
		return model.Workflow{}, err
	}
	_, err = r.database.ExecContext(ctx, `insert into "workflow" values ($1, $2, $3)`, id, name, userId)
	if err != nil {
		return model.Workflow{}, persistence.ParseDBError(err)
	}
	return model.Workflow{
		Id:     id,
		Name:   name,
		UserId: userId,
	}, nil

}

func (r WorkflowRepository) Get(ctx context.Context, id string) (model.Workflow, error) {
	w := model.Workflow{}
	err := r.database.
		QueryRowContext(ctx, `select id, name, user_id from "workflow" where id = $1`, id).
		Scan(&w.Id, &w.Name, &w.UserId)
	if err != nil {
		return model.Workflow{}, persistence.ParseDBError(err)
	}
	q := `
select id, ip, name, type
from device
         join workflow_device wd on device.id = wd.device_id
where workflow_id = $1
`
	rows, err := r.database.QueryContext(ctx, q, id)
	if err != nil {
		return model.Workflow{}, persistence.ParseDBError(err)
	}
	defer rows.Close()
	var devices []model.Device
	for rows.Next() {
		d := model.Device{}
		i := ""
		if err := rows.Scan(&d.Id, &i, &d.Name, &d.Type); err != nil {
			return model.Workflow{}, err
		}
		d.Ip = net.ParseIP(i)
		devices = append(devices, d)
	}
	w.Devices = devices

	return w, nil
}

func (r WorkflowRepository) GetAll(ctx context.Context) ([]model.Workflow, error) {
	rows, err := r.database.QueryContext(ctx, `select * from "workflow"`)
	if err != nil {
		return nil, persistence.ParseDBError(err)
	}
	defer rows.Close()
	var workflows []model.Workflow
	for rows.Next() {
		w := model.Workflow{}
		if err := rows.Scan(&w.Id, &w.Name, &w.UserId); err != nil {
			return nil, persistence.ParseDBError(err)
		}
		workflows = append(workflows, w)
	}
	if err = rows.Err(); err != nil {
		return workflows, err
	}
	return workflows, nil
}

func (r WorkflowRepository) Update() (model.Workflow, error) {
	//TODO implement me
	panic("implement me")
}

func (r WorkflowRepository) Delete(context.Context, string) error {
	//TODO implement me
	panic("implement me")
}

func (r WorkflowRepository) AssociateDevices(ctx context.Context, workflowId string, deviceIds []string) ([]string, error) {
	tx, err := r.database.Begin()
	if err != nil {
		return nil, persistence.ParseDBError(err)
	}

	_, err = tx.ExecContext(ctx, `delete from workflow_device where workflow_id = $1`, workflowId)
	if err != nil {
		tx.Rollback()
		return nil, persistence.ParseDBError(err)
	}

	for _, deviceId := range deviceIds {
		_, err = tx.ExecContext(ctx, `insert into workflow_device (workflow_id, device_id) values ($1, $2)`, workflowId, deviceId)
		if err != nil {
			tx.Rollback()
			return nil, persistence.ParseDBError(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return deviceIds, nil
}
