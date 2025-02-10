insert into project_reviews(project_id, reviewed_by, status, comment, created_at)
values (:user_id, :type_id, :status, :attributes, :created_at)
returning id