create table if not exists users_test
(
    id             uuid                     not null
        constraint users_pk primary key,
    web3_address   text                     not null,
    first_name     text                              default null,
    last_name      text                              default null,
    username       varchar(255)                      default null,
    deactivated_at timestamp with time zone          default null,
    created_at     timestamp with time zone not null default now(),
    updated_at     timestamp with time zone not null default now()
);

create unique index web3Address_uix on users_test (web3_address);
create unique index username_uix on users_test (username);