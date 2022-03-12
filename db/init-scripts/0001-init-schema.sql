CREATE SCHEMA mailmandb;
SET search_path TO mailmandb;

create table customer
(
    id    serial primary key,
    email varchar(255) not null check (email <> ''),

    constraint unique_email unique (email)
);

create table mailing_entry
(
    customer_id int          not null,
    mailing_id  int          not null,
    title       varchar(255) not null check (title <> ''),
    content     text,
    insert_time timestamp    not null,

    primary key (customer_id, mailing_id),
    constraint fk_customer foreign key (customer_id) references customer (id)
);
