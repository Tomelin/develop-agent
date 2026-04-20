package organization

import "context"

type Repository interface {
	Create(ctx context.Context, org *Organization) error
	FindByID(ctx context.Context, id string) (*Organization, error)
	FindBySlug(ctx context.Context, slug string) (*Organization, error)
	Update(ctx context.Context, org *Organization) error
}
