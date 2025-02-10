select r.id,
       r.name,
       ur.user_id
from roles as r
         join users_roles as ur on r.id = ur.role_id
where ur.user_id = $1