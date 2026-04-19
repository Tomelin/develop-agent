package interview

import "context"

type Repository interface {
	EnsureIndexes(ctx context.Context) error
	FindByProjectID(ctx context.Context, projectID string) (*InterviewSession, error)
	UpsertByProjectID(ctx context.Context, projectID string, session *InterviewSession) error
}
