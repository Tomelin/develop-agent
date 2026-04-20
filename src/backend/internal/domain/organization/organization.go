package organization

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Plan string

const (
	PlanFree       Plan = "FREE"
	PlanStarter    Plan = "STARTER"
	PlanPro        Plan = "PRO"
	PlanEnterprise Plan = "ENTERPRISE"
)

type Organization struct {
	ID                  bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name                string        `bson:"name" json:"name"`
	Slug                string        `bson:"slug" json:"slug"`
	Plan                Plan          `bson:"plan" json:"plan"`
	MaxUsers            int           `bson:"max_users" json:"max_users"`
	MaxProjectsPerMonth int           `bson:"max_projects_per_month" json:"max_projects_per_month"`
	MaxTokensPerMonth   int64         `bson:"max_tokens_per_month" json:"max_tokens_per_month"`
	BillingEmail        string        `bson:"billing_email" json:"billing_email"`
	CreatedAt           time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time     `bson:"updated_at" json:"updated_at"`
}

func New(name, slug, billingEmail string, plan Plan) (*Organization, error) {
	name = strings.TrimSpace(name)
	slug = strings.ToLower(strings.TrimSpace(slug))
	if name == "" {
		return nil, errors.New("organization name is required")
	}
	if slug == "" {
		return nil, errors.New("organization slug is required")
	}
	if !plan.IsValid() {
		return nil, errors.New("invalid organization plan")
	}
	now := time.Now().UTC()
	return &Organization{
		ID:                  bson.NewObjectID(),
		Name:                name,
		Slug:                slug,
		Plan:                plan,
		MaxUsers:            1,
		MaxProjectsPerMonth: 2,
		MaxTokensPerMonth:   10_000,
		BillingEmail:        strings.ToLower(strings.TrimSpace(billingEmail)),
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

func (p Plan) IsValid() bool {
	switch p {
	case PlanFree, PlanStarter, PlanPro, PlanEnterprise:
		return true
	default:
		return false
	}
}
