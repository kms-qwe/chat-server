-- +goose Up
create table chat (
    id integer generated always as identity primary key
);

-- +goose Down
drop table chat;

