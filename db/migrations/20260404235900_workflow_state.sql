-- +goose Up
alter table "workflow"
    add column state varchar(3) not null default 'off' check (state in ('on', 'off'));

-- +goose Down
alter table "workflow"
    drop column state;
