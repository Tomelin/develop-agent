package billing

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TriadRole string

type Provider string

const (
	RoleProducer TriadRole = "PRODUCER"
	RoleReviewer TriadRole = "REVIEWER"
	RoleRefiner  TriadRole = "REFINER"
)

const (
	ProviderOpenAI    Provider = "OPENAI"
	ProviderAnthropic Provider = "ANTHROPIC"
	ProviderGoogle    Provider = "GOOGLE"
	ProviderOllama    Provider = "OLLAMA"
)

type BillingRecord struct {
	ID                             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID                      bson.ObjectID `bson:"project_id" json:"project_id"`
	UserID                         bson.ObjectID `bson:"user_id" json:"user_id"`
	PhaseNumber                    int           `bson:"phase_number" json:"phase_number"`
	PhaseName                      string        `bson:"phase_name" json:"phase_name"`
	TriadRole                      TriadRole     `bson:"triad_role" json:"triad_role"`
	AgentID                        string        `bson:"agent_id" json:"agent_id"`
	AgentName                      string        `bson:"agent_name" json:"agent_name"`
	Provider                       Provider      `bson:"provider" json:"provider"`
	Model                          string        `bson:"model" json:"model"`
	PromptTokens                   int64         `bson:"prompt_tokens" json:"prompt_tokens"`
	CompletionTokens               int64         `bson:"completion_tokens" json:"completion_tokens"`
	TotalTokens                    int64         `bson:"total_tokens" json:"total_tokens"`
	PricePerMillionPromptTokens    float64       `bson:"price_per_million_prompt_tokens" json:"price_per_million_prompt_tokens"`
	PricePerMillionCompletionToken float64       `bson:"price_per_million_completion_tokens" json:"price_per_million_completion_tokens"`
	EstimatedCostUSD               float64       `bson:"estimated_cost_usd" json:"estimated_cost_usd"`
	DurationMs                     int64         `bson:"duration_ms" json:"duration_ms"`
	IsAutoRejection                bool          `bson:"is_auto_rejection" json:"is_auto_rejection"`
	Timestamp                      time.Time     `bson:"timestamp" json:"timestamp"`
}

type PriceItem struct {
	Provider                Provider `json:"provider" yaml:"provider"`
	Model                   string   `json:"model" yaml:"model"`
	PromptPerMillionUSD     float64  `json:"prompt_price_per_million_tokens" yaml:"prompt_price_per_million_tokens"`
	CompletionPerMillionUSD float64  `json:"completion_price_per_million_tokens" yaml:"completion_price_per_million_tokens"`
}

type ModelPricingTable struct {
	LastUpdated string      `json:"last_updated" yaml:"last_updated"`
	Models      []PriceItem `json:"models" yaml:"models"`
}

type QueryFilter struct {
	UserID    string
	ProjectID string
	Provider  string
	From      *time.Time
	To        *time.Time
	Page      int64
	Limit     int64
}

type Summary struct {
	TotalCostUSD float64           `json:"total_cost_usd"`
	TotalTokens  int64             `json:"total_tokens"`
	ByProject    []GroupedCostItem `json:"by_project"`
	ByModel      []GroupedCostItem `json:"by_model"`
}

type GroupedCostItem struct {
	Key        string  `json:"key" bson:"_id"`
	CostUSD    float64 `json:"cost_usd" bson:"cost_usd"`
	Tokens     int64   `json:"tokens" bson:"tokens"`
	Executions int64   `json:"executions" bson:"executions"`
}

type ProjectDetails struct {
	ProjectID string            `json:"project_id"`
	ByPhase   []GroupedCostItem `json:"by_phase"`
	ByAgent   []GroupedCostItem `json:"by_agent"`
	ByModel   []GroupedCostItem `json:"by_model"`
	TotalUSD  float64           `json:"total_usd"`
}

func (r *BillingRecord) ComputeTotals() {
	if r.TotalTokens == 0 {
		r.TotalTokens = r.PromptTokens + r.CompletionTokens
	}
	r.EstimatedCostUSD = (float64(r.PromptTokens)/1_000_000.0)*r.PricePerMillionPromptTokens + (float64(r.CompletionTokens)/1_000_000.0)*r.PricePerMillionCompletionToken
	if r.Timestamp.IsZero() {
		r.Timestamp = time.Now().UTC()
	}
}

func PricingKey(provider Provider, model string) string {
	return fmt.Sprintf("%s::%s", strings.ToUpper(string(provider)), strings.ToLower(strings.TrimSpace(model)))
}
