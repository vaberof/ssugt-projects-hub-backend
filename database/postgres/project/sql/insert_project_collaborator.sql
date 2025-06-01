insert into project_collaborators (project_id, user_id, role)
values (:project_id, :user_id, :role)
returning id;