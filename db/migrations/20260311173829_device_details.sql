-- +goose Up
alter table if exists "device"
    add column name text,
    add column type text;

-- +goose Down
alter table if exists "device"
    drop column name,
    drop column name;
