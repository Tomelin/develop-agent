package project

import "time"

type SecuritySeverity string

type FindingStatus string

const (
	SecuritySeverityCritical SecuritySeverity = "CRITICAL"
	SecuritySeverityHigh     SecuritySeverity = "HIGH"
	SecuritySeverityMedium   SecuritySeverity = "MEDIUM"
	SecuritySeverityLow      SecuritySeverity = "LOW"
)

const (
	FindingStatusOpen     FindingStatus = "OPEN"
	FindingStatusFixed    FindingStatus = "FIXED"
	FindingStatusAccepted FindingStatus = "ACCEPTED"
)

type SecurityFinding struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Category    string           `json:"category"`
	Severity    SecuritySeverity `json:"severity"`
	CVSS        float64          `json:"cvss"`
	CVE         string           `json:"cve,omitempty"`
	Description string           `json:"description"`
	PoC         string           `json:"poc,omitempty"`
	File        string           `json:"file,omitempty"`
	Line        int              `json:"line,omitempty"`
	Remediation string           `json:"remediation,omitempty"`
	DetectedBy  string           `json:"detected_by"`
	Status      FindingStatus    `json:"status"`
}

type SecurityAuditSummary struct {
	Score         int `json:"score"`
	CriticalCount int `json:"critical_count"`
	HighCount     int `json:"high_count"`
	MediumCount   int `json:"medium_count"`
	LowCount      int `json:"low_count"`
	TotalFindings int `json:"total_findings"`
}

type SecurityAutoRejectionResult struct {
	Triggered      bool     `json:"triggered"`
	Reason         string   `json:"reason,omitempty"`
	RetryCount     int      `json:"retry_count"`
	ReturnedPhase5 bool     `json:"returned_phase5"`
	Findings       []string `json:"findings,omitempty"`
}

type SecurityAuditReport struct {
	GeneratedAt    time.Time                   `json:"generated_at"`
	Summary        SecurityAuditSummary        `json:"summary"`
	Status         string                      `json:"status"`
	Findings       []SecurityFinding           `json:"findings"`
	Dependencies   []SecurityFinding           `json:"dependencies,omitempty"`
	StaticAnalysis []SecurityFinding           `json:"static_analysis,omitempty"`
	AutoRejection  SecurityAutoRejectionResult `json:"auto_rejection"`
}
