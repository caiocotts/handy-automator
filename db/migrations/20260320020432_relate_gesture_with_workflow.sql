-- +goose Up
drop table if exists "gesture";
create table "gesture"
(
    id    int primary key,
    label text
);

alter table if exists "workflow"
    add column gesture_id int references "gesture" (id);

-- +goose Down
alter table if exists "workflow"
    drop column gesture_id;

drop table if exists "gesture";
create table "gesture"
(
    id      text primary key,
    user_id text references "user" (id)
);

