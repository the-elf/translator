create table language
(
    id   serial primary key,
    code varchar(10) not null unique
);