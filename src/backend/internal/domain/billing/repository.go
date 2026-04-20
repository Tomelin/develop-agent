package billing

import "context"

type Repository interface {
	EnsureIndexes(ctx context.Context) error
	Create(ctx context.Context, record *BillingRecord) error
	Summary(ctx context.Context, filter QueryFilter) (*Summary, error)
	ProjectDetails(ctx context.Context, filter QueryFilter) (*ProjectDetails, error)
	ByModel(ctx context.Context, filter QueryFilter) ([]GroupedCostItem, error)
	ByPhase(ctx context.Context, filter QueryFilter) ([]GroupedCostItem, error)
	TopProjects(ctx context.Context, filter QueryFilter) ([]GroupedCostItem, error)
	List(ctx context.Context, filter QueryFilter) ([]BillingRecord, int64, error)
}
