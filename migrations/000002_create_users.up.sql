create table "user"
(
    id          serial primary key,
    language_id int not null,
    chat_id     bigint not null unique,
    constraint fk_user_language foreign key (language_id) references language (id)
);