-- +goose Up
-- +goose StatementBegin
alter table if exists "gesture"
    alter column user_id set not null;

alter table if exists "workflow"
    alter column user_id set not null;

alter table if exists "workflow_device"
    alter column workflow_id set not null,
    alter column device_id set not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists "workflow_device"
    alter column device_id drop not null,
    alter column workflow_id drop not null;

alter table if exists "workflow"
    alter column user_id drop not null;

alter table if exists "gesture"
    alter column user_id drop not null;
-- +goose StatementEnd
