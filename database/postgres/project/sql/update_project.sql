update projects
set user_id    = :user_id,
    type_id    = :type_id,
    status     = :status,
    attributes = :attributes,
    updated_at = now()
where id = :id
