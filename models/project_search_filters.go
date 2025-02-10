package models

type ProjectSearchFilters struct {
	BaseFilters      BaseFilters      `json:"baseFilters"`
	AttributeFilters AttributeFilters `json:"attributeFilters"`
}

type BaseFilters struct {
	Type   int    `json:"type"`
	Status string `json:"status"`
}

type AttributeFilters struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}
