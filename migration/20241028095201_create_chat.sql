-- +goose Up
create table chat (
    id integer generated always as identity primary key,
    created_at timestamp not null default now()
);

-- +goose Down
drop table chat;

