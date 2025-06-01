package models

import "time"

type ProjectSearchFilters struct {
	BaseFilters      BaseFilters      `json:"baseFilters"`
	AttributeFilters AttributeFilters `json:"attributeFilters"`
}

type BaseFilters struct {
	Type   int       `json:"type"`
	Status string    `json:"status"`
	Date   time.Time `json:"date"`
}

type AttributeFilters struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}
