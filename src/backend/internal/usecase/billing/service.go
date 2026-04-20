package billing

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	domain "github.com/develop-agent/backend/internal/domain/billing"
	"gopkg.in/yaml.v3"
)

type Service struct {
	repo    domain.Repository
	pricing domain.ModelPricingTable
	lookup  map[string]domain.PriceItem
}

func NewService(repo domain.Repository, pricingFile string) (*Service, error) {
	table, err := loadPricingTable(pricingFile)
	if err != nil {
		return nil, err
	}
	lookup := make(map[string]domain.PriceItem, len(table.Models))
	for _, m := range table.Models {
		lookup[domain.PricingKey(m.Provider, m.Model)] = m
	}
	return &Service{repo: repo, pricing: table, lookup: lookup}, nil
}

func (s *Service) Pricing(_ context.Context) domain.ModelPricingTable {
	return s.pricing
}

func (s *Service) CreateRecord(ctx context.Context, rec *domain.BillingRecord) error {
	if rec == nil {
		return errors.New("billing record is required")
	}
	if rec.Model == "" || rec.Provider == "" {
		return errors.New("provider and model are required")
	}
	if p, ok := s.lookup[domain.PricingKey(rec.Provider, rec.Model)]; ok {
		rec.PricePerMillionPromptTokens = p.PromptPerMillionUSD
		rec.PricePerMillionCompletionToken = p.CompletionPerMillionUSD
	}
	rec.ComputeTotals()
	return s.repo.Create(ctx, rec)
}

func (s *Service) Summary(ctx context.Context, filter domain.QueryFilter) (*domain.Summary, error) {
	return s.repo.Summary(ctx, filter)
}

func (s *Service) ProjectDetails(ctx context.Context, filter domain.QueryFilter) (*domain.ProjectDetails, error) {
	return s.repo.ProjectDetails(ctx, filter)
}

func (s *Service) ByModel(ctx context.Context, filter domain.QueryFilter) ([]domain.GroupedCostItem, error) {
	return s.repo.ByModel(ctx, filter)
}

func (s *Service) ByPhase(ctx context.Context, filter domain.QueryFilter) ([]domain.GroupedCostItem, error) {
	return s.repo.ByPhase(ctx, filter)
}

func (s *Service) TopProjects(ctx context.Context, filter domain.QueryFilter) ([]domain.GroupedCostItem, error) {
	return s.repo.TopProjects(ctx, filter)
}

func (s *Service) Export(ctx context.Context, filter domain.QueryFilter, format string) ([]byte, string, error) {
	records, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, "", err
	}
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		out, err := json.MarshalIndent(records, "", "  ")
		return out, "application/json", err
	case "csv":
		return encodeCSV(records)
	default:
		return nil, "", errors.New("unsupported format (use csv or json)")
	}
}

func ParseTimeRange(fromRaw, toRaw string) (*time.Time, *time.Time, error) {
	var from, to *time.Time
	if strings.TrimSpace(fromRaw) != "" {
		parsed, err := time.Parse(time.RFC3339, fromRaw)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid from")
		}
		from = &parsed
	}
	if strings.TrimSpace(toRaw) != "" {
		parsed, err := time.Parse(time.RFC3339, toRaw)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid to")
		}
		to = &parsed
	}
	return from, to, nil
}

func loadPricingTable(path string) (domain.ModelPricingTable, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return domain.ModelPricingTable{}, err
	}
	var table domain.ModelPricingTable
	if err := yaml.Unmarshal(data, &table); err != nil {
		return domain.ModelPricingTable{}, err
	}
	return table, nil
}

func encodeCSV(records []domain.BillingRecord) ([]byte, string, error) {
	rows := [][]string{{"timestamp", "project_id", "phase_number", "phase_name", "triad_role", "agent_id", "agent_name", "provider", "model", "prompt_tokens", "completion_tokens", "total_tokens", "estimated_cost_usd", "duration_ms", "is_auto_rejection"}}
	for _, r := range records {
		rows = append(rows, []string{r.Timestamp.Format(time.RFC3339), r.ProjectID.Hex(), strconv.Itoa(r.PhaseNumber), r.PhaseName, string(r.TriadRole), r.AgentID, r.AgentName, string(r.Provider), r.Model, strconv.FormatInt(r.PromptTokens, 10), strconv.FormatInt(r.CompletionTokens, 10), strconv.FormatInt(r.TotalTokens, 10), strconv.FormatFloat(r.EstimatedCostUSD, 'f', 8, 64), strconv.FormatInt(r.DurationMs, 10), strconv.FormatBool(r.IsAutoRejection)})
	}
	var b strings.Builder
	w := csv.NewWriter(&b)
	if err := w.WriteAll(rows); err != nil {
		return nil, "", err
	}
	return []byte(b.String()), "text/csv", nil
}
