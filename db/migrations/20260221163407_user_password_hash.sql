-- +goose Up
-- +goose StatementBegin
alter table if exists "user"
    add column hash text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists "user"
    drop column hash;
-- +goose StatementEnd
