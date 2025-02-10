insert into projects(user_id, type_id, status, attributes, created_at, updated_at)
values (:user_id, :type_id, :status, :attributes, :created_at, now())
returning id