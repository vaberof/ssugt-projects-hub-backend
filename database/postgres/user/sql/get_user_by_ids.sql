select id,
       role_id,
       email,
       password_hash,
       full_name,
       created_at,
       updated_at
from users
where id = any($1)
