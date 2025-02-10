insert into users(email, password_hash, full_name, phone_number, created_at, updated_at)
values (:email, :password_hash, :full_name, :phone_number, :created_at, now())
returning id