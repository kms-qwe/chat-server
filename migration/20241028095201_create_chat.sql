-- +goose Up
create table chatV1.chat (
    id integer generated always as identity primary key
);

-- +goose Down
drop table chatV1.chat;

