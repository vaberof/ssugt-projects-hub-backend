update project_collaborators
set project_id = :project_id,
    user_id    = :user_id,
    role       = :role
where id = :id