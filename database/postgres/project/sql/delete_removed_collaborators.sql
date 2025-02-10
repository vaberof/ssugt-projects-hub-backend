delete
from project_collaborators
where project_id = $1
  and not id = any ($2)