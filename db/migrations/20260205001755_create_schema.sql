-- +goose Up
-- +goose StatementBegin
create table if not exists "user"
(
    id       varchar(12) primary key,
    username text
);

create table if not exists "workflow"
(
    id   varchar(12) primary key,
    name text
);

create table if not exists "device"
(
    id varchar(12) primary key,
    ip varchar(15)
);

create table if not exists "gesture"
(
    id varchar(12) primary key
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "gesture";
drop table if exists "device";
drop table if exists "workflow";
drop table if exists "user";
-- +goose StatementEnd
