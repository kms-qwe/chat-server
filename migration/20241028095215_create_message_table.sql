-- +goose Up
create table chatV1.message (
    id integer generated always as identity primary key,
    user_name text not null,
    message_text text, 
    chat_id integer not null,
    time_stamp timestamp,

    constraint fk_chat_id foreign key (chat_id) references chatV1.chat(id) on delete cascade
);


-- +goose Down
drop table chatV1.message;
