package views

type StudentScientificConferenceProjectTemplate struct {
	Title            string  `json:"title"`
	Object           string  `json:"object"`
	Summary          string  `json:"summary"`
	Cost             float64 `json:"cost"`
	DevelopingStage  string  `json:"developing_stage"`
	RealizationTerm  string  `json:"realization_term"`
	ApplicationScope string  `json:"application_scope"`
}

type LaboratoryProjectTemplate struct {
	LaboratoryName   string   `json:"laboratory_name"`
	Title            string   `json:"title"`
	Object           string   `json:"object,omitempty"`
	Summary          string   `json:"summary,omitempty"`
	Problematic      string   `json:"problematic,omitempty"`
	Solution         string   `json:"solution,omitempty"`
	Functionality    string   `json:"functionality,omitempty"`
	TechnologyStack  string   `json:"technology_stack,omitempty"`
	Advantages       []string `json:"advantages,omitempty"`
	Cost             float64  `json:"cost,omitempty"`
	DevelopingStage  string   `json:"developing_stage,omitempty"`
	RealizationTerm  string   `json:"realization_term,omitempty"`
	ApplicationScope string   `json:"application_scope,omitempty"`
}
