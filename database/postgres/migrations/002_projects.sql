-- +goose Up

create table if not exists project_types
(
    id   serial primary key,
    name text not null
);

insert into project_types(name)
values ('science'),
       ('laboratory');

create table if not exists projects
(
    id         serial primary key,
    user_id    int       not null references users (id),
    type_id    int       not null references project_types (id),
    status     text      not null,
    attributes jsonb     not null,
    is_deleted bool      not null default false,
    created_at timestamp not null,
    updated_at timestamp not null
);

create index if not exists idx_projects_user_id on projects (user_id);
create index if not exists idx_projects_type_id on projects (type_id);
create index if not exists idx_projects_attributes_title on projects (lower(attributes ->> 'title'));
create index if not exists idx_projects_attributes_tags on projects using GIN (attributes);

create table if not exists project_reviews
(
    id          serial primary key,
    project_id  int       not null references projects (id),
    reviewed_by int       not null references users (id),
    status      text      not null,
    comment     text      not null,
    created_at  timestamp not null
);

create index if not exists idx_project_reviews_project_id on project_reviews (project_id);
create index if not exists idx_project_reviews_reviewed_by on project_reviews (reviewed_by);

-- +goose Down

drop index if exists idx_project_reviews_project_id;
drop index if exists idx_project_reviews_reviewed_by;
drop table if exists project_reviews;

drop index if exists idx_projects_user_id;
drop index if exists idx_projects_type_id;
drop index if exists idx_projects_attributes_title;
drop index if exists idx_projects_attributes_tags;
drop table if exists projects;

drop table if exists project_types;