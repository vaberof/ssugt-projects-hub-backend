select distinct p.id,
                p.user_id,
                p.type_id,
                p.status,
                p.attributes,
                p.created_at,
                p.updated_at
from projects as p
         left join project_collaborators as pc on p.id = pc.project_id
where not p.is_deleted
  and (p.user_id = $1 or pc.user_id = $1)