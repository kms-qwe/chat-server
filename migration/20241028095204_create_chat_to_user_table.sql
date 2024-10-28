-- +goose Up
create table chatV1.chat_to_user (
    chat_id integer not null,
    user_name text not null,

    constraint unique_chat_user unique (chat_id, user_name),
    constraint fk_chat_id foreign key (chat_id) references chatV1.chat(id) on delete cascade
);

-- +goose Down
drop table chatV1.chat_to_user;

