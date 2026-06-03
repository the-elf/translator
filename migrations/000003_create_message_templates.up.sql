create table message_template
(
    id          serial primary key,
    language_id int     not null,
    code        varchar not null,
    text        text    not null,
    constraint fk_message_template_language foreign key (language_id) references language (id),
    constraint uq_message_template_language_code unique (language_id, code)
);