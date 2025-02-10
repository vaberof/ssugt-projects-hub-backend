-- +goose Up

create table if not exists project_collaborators
(
    id         serial primary key,
    project_id int  not null references projects (id),
    user_id    int  not null references users (id),
    role       text not null
);

create index if not exists idx_project_collaborators_project_id on project_collaborators (project_id);
create index if not exists idx_project_collaborators_user_id on project_collaborators (user_id);

-- +goose Down

drop index if exists idx_project_collaborators_project_id;
drop index if exists idx_project_collaborators_user_id;
drop table if exists project_collaborators;
