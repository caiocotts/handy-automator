-- +goose Up
alter table if exists "user"
add unique (username);

-- +goose Down
alter table if exists "user"
drop constraint user_username_key;
