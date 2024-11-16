-- +goose Up
create table chat_participants (
    chat_id integer,
    user_name text,

    constraint pk_chat_id_user_name primary key(chat_id, user_name),
    constraint fk_chat_id foreign key (chat_id) references chat(id) on delete cascade
);

-- +goose Down
drop table chat_participants;

