package models

import (
	"encoding/json"
	"time"
)

type Changes json.RawMessage

type ProjectReview struct {
	Id         int `json:"id"`
	ReviewedBy int `json:"reviewedBy"`
	//Changes     Changes       `json:"changes"`
	Status      ProjectStatus `json:"status"`
	Comment     string        `json:"comment"`
	CreatedAt   time.Time     `json:"createdAt"`
	ProjectData Project       `json:"projectData"`
}
