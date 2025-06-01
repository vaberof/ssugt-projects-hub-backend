package project

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

type DbProject struct {
	Id            int            `db:"id"`
	UserId        int            `db:"user_id"`
	TypeId        int            `db:"type_id"`
	Status        string         `db:"status"`
	Attributes    types.JSONText `db:"attributes"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
	Collaborators []DbCollaborator
}

type DbProjectReview struct {
	Id          int       `db:"review_id"`
	ProjectId   int       `db:"review_project_id"`
	ReviewedBy  int       `db:"reviewed_by"`
	Status      string    `db:"review_status"`
	Comment     string    `db:"review_comment"`
	CreatedAt   time.Time `db:"review_created_at"`
	ProjectData DbProject
}

type DbCollaborator struct {
	Id        int    `db:"id"`
	ProjectId int    `db:"project_id"`
	UserId    int    `db:"user_id"`
	Role      string `db:"role"`
}
