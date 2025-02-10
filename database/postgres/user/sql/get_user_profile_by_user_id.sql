select id,
       user_id,
       personal_info,
       created_at,
       updated_at
from user_profiles
where user_id = $1