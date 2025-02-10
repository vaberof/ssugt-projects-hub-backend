-- +goose Up

create table if not exists users
(
    id                        serial primary key,
    email                     text unique not null,
    password_hash             text        not null,
    full_name                 text        not null,
    phone_number              text unique not null,
    is_email_confirmed        bool        not null default false,
    is_phone_number_confirmed bool        not null default false,
    created_at                timestamp   not null,
    updated_at                timestamp   not null
);

create index if not exists idx_users_email on users (email);

create table if not exists user_profiles
(
    id            serial primary key,
    user_id       int       not null references users (id),
    personal_info jsonb     not null,
    settings      jsonb     not null,
    created_at    timestamp not null,
    updated_at    timestamp not null
);

create index if not exists idx_user_profiles_user_id on user_profiles (user_id);

create table if not exists roles
(
    id   serial primary key,
    name text not null
);

insert into roles(name)
values ('user'),
       ('moderator'),
       ('admin');

create table if not exists users_roles
(
    user_id int references users (id),
    role_id int references roles (id)
);

-- +goose Down

drop table if exists users_roles;
drop table if exists roles;

drop index if exists idx_user_profiles_user_id;
drop table if exists user_profiles;

drop index if exists idx_users_email;
drop table if exists users;
