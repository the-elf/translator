create table button
(
    id          serial primary key,
    language_id int          not null,
    group_name  varchar(100) not null,
    data        varchar(64)  not null,
    title       text         not null,
    constraint fk_button_language foreign key (language_id) references language (id),
    constraint uq_button_name_data_lang unique (group_name, data, language_id)
);