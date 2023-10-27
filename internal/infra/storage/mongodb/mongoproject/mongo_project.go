package mongoproject

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Project struct {
	Id                        primitive.ObjectID         `bson:"_id,omitempty"`
	UserId                    primitive.ObjectID         `bson:"user_id"`
	ProjectType               string                     `bson:"project_type"`
	Authors                   []*Author                  `bson:"authors"`
	Organization              *Organization              `bson:"organization"`
	Director                  *Director                  `bson:"director"`
	SscProjectTemplate        *SSCProjectTemplate        `bson:"ssc_template,omitempty"`
	LaboratoryProjectTemplate *LaboratoryProjectTemplate `bson:"laboratory_template,omitempty"`
	Tags                      []string                   `bson:"tags,omitempty"`
	CreatedAt                 time.Time                  `bson:"created_at"`
	ModifiedAt                time.Time                  `bson:"modified_at"`
}

type SSCProjectTemplate struct {
	Title            string  `json:"title" bson:"title"`
	Object           string  `json:"object" bson:"object"`
	Summary          string  `json:"summary" bson:"summary"`
	Cost             float64 `json:"cost" bson:"cost"`
	DevelopingStage  string  `json:"developing_stage" bson:"developing_stage"`
	RealizationTerm  string  `json:"realization_term" bson:"realization_term"`
	ApplicationScope string  `json:"application_scope" bson:"application_scope"`
}

type LaboratoryProjectTemplate struct {
	LaboratoryName   string   `json:"laboratory_name" bson:"laboratory_name"`
	Title            string   `json:"title" bson:"title"`
	Problematic      string   `json:"problematic" bson:"problematic,omitempty"`
	Solution         string   `json:"solution" bson:"solution,omitempty"`
	Functionality    string   `json:"functionality" bson:"functionality,omitempty"`
	TechnologyStack  string   `json:"technology_stack" bson:"technology_stack,omitempty"`
	Advantages       []string `json:"advantages" bson:"advantages,omitempty"`
	Object           string   `json:"object" bson:"object,omitempty"`
	Summary          string   `json:"summary" bson:"summary,omitempty"`
	Cost             float64  `json:"cost" bson:"cost,omitempty"`
	DevelopingStage  string   `json:"developing_stage" bson:"developing_stage,omitempty"`
	RealizationTerm  string   `json:"realization_term" bson:"realization_term,omitempty"`
	ApplicationScope string   `json:"application_scope" bson:"application_scope,omitempty"`
}

type Author struct {
	FullName string `bson:"full_name"`
	Degree   string `bson:"degree,omitempty"`
	Course   int    `bson:"course,omitempty"`
	Group    string `bson:"group,omitempty"`
}

type Organization struct {
	Name    string `bson:"name,omitempty"`
	Address string `bson:"address,omitempty"`
}

type Director struct {
	FullName string `bson:"full_name"`
	Email    string `bson:"email,omitempty"`
	Phone    string `bson:"phone,omitempty"`
}
