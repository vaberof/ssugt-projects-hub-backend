-- +goose Up

create table if not exists roles
(
    id   serial primary key,
    name text not null
);

insert into roles(name)
values ('user'),
       ('admin');

create table if not exists users
(
    id            serial primary key,
    role_id       int references roles (id),
    email         text unique not null,
    password_hash text        not null,
    full_name     text        not null,
    created_at    timestamp   not null,
    updated_at    timestamp   not null
);

create index if not exists idx_users_email on users (email);
create index if not exists idx_users_full_name on users (full_name);

create table if not exists user_profiles
(
    id            serial primary key,
    user_id       int       not null references users (id),
    personal_info jsonb     not null,
    created_at    timestamp not null,
    updated_at    timestamp not null
);

create index if not exists idx_user_profiles_user_id on user_profiles (user_id);

-- +goose Down

drop index if exists idx_user_profiles_user_id;
drop table if exists user_profiles;

drop index if exists idx_users_email;
drop index if exists idx_users_full_name;
drop table if exists users;

drop table if exists roles;
