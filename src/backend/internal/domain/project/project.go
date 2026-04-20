package project

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FlowType string

type ProjectStatus string

type PhaseStatus string

type Track string

type TaskType string

type TaskComplexity string

type TaskStatus string

type ExecutionMode string

const (
	FlowSoftware    FlowType = "SOFTWARE"
	FlowLandingPage FlowType = "LANDING_PAGE"
	FlowMarketing   FlowType = "MARKETING"
)

const (
	ProjectDraft      ProjectStatus = "DRAFT"
	ProjectInProgress ProjectStatus = "IN_PROGRESS"
	ProjectPaused     ProjectStatus = "PAUSED"
	ProjectCompleted  ProjectStatus = "COMPLETED"
	ProjectArchived   ProjectStatus = "ARCHIVED"
)

const (
	PhasePending    PhaseStatus = "PENDING"
	PhaseInProgress PhaseStatus = "IN_PROGRESS"
	PhaseReview     PhaseStatus = "REVIEW"
	PhaseCompleted  PhaseStatus = "COMPLETED"
	PhaseRejected   PhaseStatus = "REJECTED"
)

const (
	TrackFull     Track = "FULL"
	TrackFrontend Track = "FRONTEND"
	TrackBackend  Track = "BACKEND"
)

const (
	TaskTypeFrontend TaskType = "FRONTEND"
	TaskTypeBackend  TaskType = "BACKEND"
	TaskTypeInfra    TaskType = "INFRA"
	TaskTypeTest     TaskType = "TEST"
	TaskTypeDoc      TaskType = "DOC"
)

const (
	ComplexityLow      TaskComplexity = "LOW"
	ComplexityMedium   TaskComplexity = "MEDIUM"
	ComplexityHigh     TaskComplexity = "HIGH"
	ComplexityCritical TaskComplexity = "CRITICAL"
)

const (
	TaskTodo       TaskStatus = "TODO"
	TaskInProgress TaskStatus = "IN_PROGRESS"
	TaskDone       TaskStatus = "DONE"
	TaskBlocked    TaskStatus = "BLOCKED"
)

const (
	ExecutionModeAutomatic ExecutionMode = "AUTOMATIC"
	ExecutionModeManual    ExecutionMode = "MANUAL"
)

type AgentTriad struct {
	Producer  string `bson:"producer" json:"producer"`
	Reviewer  string `bson:"reviewer" json:"reviewer"`
	Refiner   string `bson:"refiner" json:"refiner"`
	Provider  string `bson:"provider,omitempty" json:"provider,omitempty"`
	ModelHint string `bson:"model_hint,omitempty" json:"model_hint,omitempty"`
}

type PhaseExecution struct {
	PhaseNumber   int              `bson:"phase_number" json:"phase_number"`
	PhaseName     string           `bson:"phase_name" json:"phase_name"`
	Status        PhaseStatus      `bson:"status" json:"status"`
	Track         Track            `bson:"track" json:"track"`
	Tracks        []TrackExecution `bson:"tracks,omitempty" json:"tracks,omitempty"`
	FeedbackCount int              `bson:"feedback_count" json:"feedback_count"`
	FeedbackLimit int              `bson:"feedback_limit" json:"feedback_limit"`
	Artifacts     []string         `bson:"artifacts,omitempty" json:"artifacts,omitempty"`
	AgentTriad    AgentTriad       `bson:"agent_triad" json:"agent_triad"`
	StartedAt     *time.Time       `bson:"started_at,omitempty" json:"started_at,omitempty"`
	CompletedAt   *time.Time       `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}

type TrackExecution struct {
	Track       Track       `bson:"track" json:"track"`
	Status      PhaseStatus `bson:"status" json:"status"`
	StartedAt   *time.Time  `bson:"started_at,omitempty" json:"started_at,omitempty"`
	CompletedAt *time.Time  `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}

type TransitionRecord struct {
	Kind      string      `bson:"kind" json:"kind"`
	From      string      `bson:"from" json:"from"`
	To        string      `bson:"to" json:"to"`
	At        time.Time   `bson:"at" json:"at"`
	Reason    string      `bson:"reason,omitempty" json:"reason,omitempty"`
	Phase     int         `bson:"phase,omitempty" json:"phase,omitempty"`
	Triggered string      `bson:"triggered_by,omitempty" json:"triggered_by,omitempty"`
	Meta      interface{} `bson:"meta,omitempty" json:"meta,omitempty"`
}

type Project struct {
	ID                 bson.ObjectID      `bson:"_id,omitempty" json:"id"`
	Name               string             `bson:"name" json:"name"`
	Description        string             `bson:"description" json:"description"`
	FlowType           FlowType           `bson:"flow_type" json:"flow_type"`
	Status             ProjectStatus      `bson:"status" json:"status"`
	CurrentPhaseNumber int                `bson:"current_phase_number" json:"current_phase_number"`
	Phases             []PhaseExecution   `bson:"phases" json:"phases"`
	LinkedProjectID    *bson.ObjectID     `bson:"linked_project_id,omitempty" json:"linked_project_id,omitempty"`
	OwnerUserID        bson.ObjectID      `bson:"owner_user_id" json:"owner_user_id"`
	DynamicModeEnabled bool               `bson:"dynamic_mode_enabled" json:"dynamic_mode_enabled"`
	SpecMD             string             `bson:"spec_md,omitempty" json:"spec_md,omitempty"`
	RoadmapJSON        string             `bson:"roadmap_json,omitempty" json:"roadmap_json,omitempty"`
	TotalTokensUsed    int64              `bson:"total_tokens_used" json:"total_tokens_used"`
	TotalCostUSD       float64            `bson:"total_cost_usd" json:"total_cost_usd"`
	BudgetUSD          float64            `bson:"budget_usd,omitempty" json:"budget_usd,omitempty"`
	BudgetAlerted80    bool               `bson:"budget_alerted_80,omitempty" json:"budget_alerted_80,omitempty"`
	TransitionHistory  []TransitionRecord `bson:"transition_history,omitempty" json:"transition_history,omitempty"`
	Phase5Mode         ExecutionMode      `bson:"phase_5_mode,omitempty" json:"phase_5_mode,omitempty"`
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"updated_at"`
	ArchivedAt         *time.Time         `bson:"archived_at,omitempty" json:"archived_at,omitempty"`
}

type Task struct {
	ID              bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	ProjectID       bson.ObjectID  `bson:"project_id" json:"project_id"`
	PhaseID         string         `bson:"phase_id,omitempty" json:"phase_id,omitempty"`
	EpicID          string         `bson:"epic_id,omitempty" json:"epic_id,omitempty"`
	Title           string         `bson:"title" json:"title"`
	Description     string         `bson:"description,omitempty" json:"description,omitempty"`
	Type            TaskType       `bson:"type" json:"type"`
	Complexity      TaskComplexity `bson:"complexity" json:"complexity"`
	EstimatedHours  float64        `bson:"estimated_hours" json:"estimated_hours"`
	Track           Track          `bson:"track,omitempty" json:"track,omitempty"`
	Dependencies    []string       `bson:"dependencies,omitempty" json:"dependencies,omitempty"`
	Status          TaskStatus     `bson:"status" json:"status"`
	AssignedAgentID string         `bson:"assigned_agent_id,omitempty" json:"assigned_agent_id,omitempty"`
	CreatedAt       time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `bson:"updated_at" json:"updated_at"`
}

func NewProject(name, description string, flowType FlowType, ownerUserID bson.ObjectID, dynamic bool, linkedProjectID *bson.ObjectID) (*Project, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	if !flowType.IsValid() {
		return nil, errors.New("invalid flow type")
	}

	now := time.Now().UTC()
	return &Project{
		ID:                 bson.NewObjectID(),
		Name:               name,
		Description:        strings.TrimSpace(description),
		FlowType:           flowType,
		Status:             ProjectDraft,
		CurrentPhaseNumber: 1,
		Phases:             defaultPhases(),
		LinkedProjectID:    linkedProjectID,
		OwnerUserID:        ownerUserID,
		DynamicModeEnabled: dynamic,
		Phase5Mode:         ExecutionModeManual,
		TotalTokensUsed:    0,
		TotalCostUSD:       0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func defaultPhases() []PhaseExecution {
	return []PhaseExecution{
		{PhaseNumber: 1, PhaseName: "Criação do Projeto", Status: PhasePending, Track: TrackFull, FeedbackLimit: 10},
		{PhaseNumber: 2, PhaseName: "Engenharia de Software", Status: PhasePending, Track: TrackFull, Tracks: defaultSplitTracks(), FeedbackLimit: 5},
		{PhaseNumber: 3, PhaseName: "Arquitetura de Software", Status: PhasePending, Track: TrackFull, Tracks: defaultSplitTracks(), FeedbackLimit: 5},
		{PhaseNumber: 4, PhaseName: "Planejamento", Status: PhasePending, Track: TrackFull, FeedbackLimit: 5},
		{PhaseNumber: 5, PhaseName: "Desenvolvimento", Status: PhasePending, Track: TrackFull, FeedbackLimit: 5},
		{PhaseNumber: 6, PhaseName: "Testes", Status: PhasePending, Track: TrackFull, FeedbackLimit: 5},
		{PhaseNumber: 7, PhaseName: "Segurança", Status: PhasePending, Track: TrackFull, FeedbackLimit: 5},
		{PhaseNumber: 8, PhaseName: "Documentação", Status: PhasePending, Track: TrackFull, FeedbackLimit: 5},
		{PhaseNumber: 9, PhaseName: "DevOps e Deploy", Status: PhasePending, Track: TrackFull, FeedbackLimit: 5},
	}
}

func defaultSplitTracks() []TrackExecution {
	return []TrackExecution{
		{Track: TrackFrontend, Status: PhasePending},
		{Track: TrackBackend, Status: PhasePending},
	}
}

func (f FlowType) IsValid() bool {
	switch f {
	case FlowSoftware, FlowLandingPage, FlowMarketing:
		return true
	default:
		return false
	}
}

func (s ProjectStatus) IsValid() bool {
	switch s {
	case ProjectDraft, ProjectInProgress, ProjectPaused, ProjectCompleted, ProjectArchived:
		return true
	default:
		return false
	}
}

func (s PhaseStatus) IsValid() bool {
	switch s {
	case PhasePending, PhaseInProgress, PhaseReview, PhaseCompleted, PhaseRejected:
		return true
	default:
		return false
	}
}

func (t TaskType) IsValid() bool {
	switch t {
	case TaskTypeFrontend, TaskTypeBackend, TaskTypeInfra, TaskTypeTest, TaskTypeDoc:
		return true
	default:
		return false
	}
}

func (c TaskComplexity) IsValid() bool {
	switch c {
	case ComplexityLow, ComplexityMedium, ComplexityHigh, ComplexityCritical:
		return true
	default:
		return false
	}
}

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskTodo, TaskInProgress, TaskDone, TaskBlocked:
		return true
	default:
		return false
	}
}

func (m ExecutionMode) IsValid() bool {
	switch m {
	case ExecutionModeAutomatic, ExecutionModeManual:
		return true
	default:
		return false
	}
}
