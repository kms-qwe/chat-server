-- +goose Up
create table log(
    id integer generated always as identity primary key,
    message text not null,
    log_time timestamp not null default now()
);

-- +goose Down
drop table log;