-- +goose Up
alter table if exists "user"
    drop column face_embedding,
    add column face_embedding jsonb;

-- +goose Down
alter table if exists "user"
    drop column face_embedding,
    add column face_embedding text;
