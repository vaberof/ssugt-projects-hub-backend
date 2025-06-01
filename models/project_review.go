package models

import (
	"time"
)

type ProjectReview struct {
	Id          int           `json:"id"`
	ReviewedBy  int           `json:"reviewedBy"`
	Status      ProjectStatus `json:"status"`
	Comment     string        `json:"comment"`
	CreatedAt   time.Time     `json:"createdAt"`
	ProjectData Project       `json:"projectData"`
}
