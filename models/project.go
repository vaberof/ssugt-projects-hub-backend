package models

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"time"
)

type ProjectType int

const (
	ScienceProjectType ProjectType = iota + 1
	LaboratoryProjectType
)

type ProjectStatus string

const (
	InProcessProjectStatus ProjectStatus = "В обработке"
	RefusedProjectStatus   ProjectStatus = "Отклонено"
	CompletedProjectStatus ProjectStatus = "Подтверждено"
)

type Project struct {
	Id            int             `json:"id"`
	UserId        int             `json:"userId"`
	Type          ProjectType     `json:"type"`
	Status        ProjectStatus   `json:"status"`
	Attributes    json.RawMessage `json:"attributes"`
	Collaborators []Collaborator  `json:"collaborators"`
	Files         []ProjectFile   `json:"files"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

type DevelopingStage string

const (
	ConceptDevelopingStage     DevelopingStage = "Концепт"
	ResearchDevelopingStage    DevelopingStage = "Исследование"
	DevelopmentDevelopingStage DevelopingStage = "Разработка"
	TestingDevelopingStage     DevelopingStage = "Тестирование"
	DeploymentDevelopingStage  DevelopingStage = "Внедрение"
	CompletedDevelopingStage   DevelopingStage = "Завершён"
	OnHoldDevelopingStage      DevelopingStage = "Приостановлен"
	CancelledDevelopingStage   DevelopingStage = "Отменён"
)

type FundingSource string

const (
	GovernmentGrantFundingSource FundingSource = "Государственный грант"
	PrivateFundingSource         FundingSource = "Частные инвесторы/компании"
	PersonalFundingSource        FundingSource = "Личные средства"

	OtherFundingSource FundingSource = "Другое"
)

type BaseProjectAttributes struct {
	Title            string          `json:"title"`            // Название проекта
	Object           string          `json:"object"`           // Объект исследования/разработки
	Summary          string          `json:"summary"`          // Краткое описание проекта
	Cost             decimal.Decimal `json:"cost"`             // Стоимость проекта
	DevelopingStage  DevelopingStage `json:"developingStage"`  // Текущий этап разработки
	RealizationTerm  string          `json:"realizationTerm"`  // Срок реализации
	ApplicationScope string          `json:"applicationScope"` // Область применения
	Tags             []string        `json:"tags"`             // Теги проекта
	FundingSource    FundingSource   `json:"fundingSource"`    // Источник финансирования
	TeamSize         int             `json:"teamSize"`         // Размер команды
}

type ScienceProjectAttributes struct {
	BaseProjectAttributes
	ResearchGoals   string `json:"researchGoals"`   // Цели исследования
	Methodology     string `json:"methodology"`     // Методология исследования
	PotentialImpact string `json:"potentialImpact"` // Потенциальное влияние
}

type LaboratoryProjectAttributes struct {
	BaseProjectAttributes
	LaboratoryName  string   `json:"laboratoryName"`  // Название лаборатории
	Problematic     string   `json:"problematic"`     // Проблематика проекта
	Solution        string   `json:"solution"`        // Предлагаемое решение
	Functionality   string   `json:"functionality"`   // Основная функциональность
	TechnologyStack string   `json:"technologyStack"` // Стек технологий
	Advantages      []string `json:"advantages"`      // Преимущества решения
	TestResults     string   `json:"testResults"`     // Результаты тестирования
	DeploymentPlan  string   `json:"deploymentPlan"`  // План внедрения
}
