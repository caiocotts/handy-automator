-- +goose Up
alter table if exists "user"
    add column face_embedding text;

-- +goose Down
alter table if exists "user"
    drop column face_embedding;
