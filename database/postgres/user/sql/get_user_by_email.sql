select id,
       email,
       password_hash,
       full_name,
       phone_number,
       is_email_confirmed,
       is_phone_number_confirmed,
       created_at,
       updated_at
from users
where email = $1
