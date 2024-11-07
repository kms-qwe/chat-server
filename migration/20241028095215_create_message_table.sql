-- +goose Up
create table message (
    id integer generated always as identity primary key,
    user_name text not null,
    message_text text, 
    chat_id integer not null,
    message_time_send timestamp,

    constraint fk_chat_id foreign key (chat_id) references chat(chat_id) on delete cascade
);


-- +goose Down
drop table message;
