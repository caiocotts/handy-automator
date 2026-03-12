-- +goose Up
-- +goose Statement Begin
alter table if exists "user"
    ADD column face_embedding text;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
alter table if exists "user"
    drop column face_embedding;
SELECT 'down SQL query';
