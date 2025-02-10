select id,
       user_id,
       type_id,
       status,
       attributes,
       created_at,
       updated_at
from projects
where not is_deleted
  and id = $1
