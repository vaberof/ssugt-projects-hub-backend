select id, project_id, user_id, role
from project_collaborators
where project_id = any($1)