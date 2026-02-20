-- +goose Up
-- +goose StatementBegin
alter table "gesture"
    add column user_id varchar(12) references "user" (id);

alter table "workflow"
    add column user_id varchar(12) references "user" (id);

create table if not exists "workflow_device"
(
    workflow_id varchar(12) references "workflow" (id),
    device_id   varchar(12) references "device" (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "workflow_device";

alter table "workflow"
    drop column user_id;

alter table "gesture"
    drop column user_id;
-- +goose StatementEnd
