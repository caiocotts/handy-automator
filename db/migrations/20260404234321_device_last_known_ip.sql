-- +goose Up
alter table "device"
    add column last_known_ip varchar(15);

-- +goose Down
alter table "device"
    drop column last_known_ip;
