package project

import "time"

type Phase13RunInput struct {
	BackendBaseURL string `json:"backend_base_url"`
	FrontendURL    string `json:"frontend_url"`
	IncludeDevOps  bool   `json:"include_devops"`
}

type Phase13DeliveryReport struct {
	GeneratedAt   time.Time `json:"generated_at"`
	ProjectID     string    `json:"project_id"`
	ProjectStatus string    `json:"project_status"`
	IncludeDevOps bool      `json:"include_devops"`
	Artifacts     []string  `json:"artifacts"`
	Warnings      []string  `json:"warnings,omitempty"`
}
