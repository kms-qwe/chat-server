-- +goose Up
create table chat (
    chat_id integer generated always as identity primary key
);

-- +goose Down
drop table chat;

