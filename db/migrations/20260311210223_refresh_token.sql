-- +goose Up
alter table if exists "user"
    add column refresh_token text;

-- +goose Down
alter table if exists "user"
    drop column refresh_token;
