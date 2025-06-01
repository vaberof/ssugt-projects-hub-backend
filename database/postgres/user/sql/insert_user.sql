insert into users(role_id, email, password_hash, full_name, created_at, updated_at)
values (:role_id, :email, :password_hash, :full_name, :created_at, now())
returning id