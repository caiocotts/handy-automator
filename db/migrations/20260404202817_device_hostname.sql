-- +goose Up
alter table "device"
    drop column ip;
alter table "device"
    add column hostname text not null unique;

-- +goose Down
alter table "device"
    drop column hostname;
alter table "device"
    add column ip varchar(15);
