package agentsdk

import (
	"context"
	"time"
)

// RetryProvider decorates a provider with bounded retries.
type RetryProvider struct {
	Base        Provider
	MaxRetries  int
	Backoff     time.Duration
	ShouldRetry func(error) bool
}

func (r *RetryProvider) Initialize(ctx context.Context, cfg Config) error {
	return r.Base.Initialize(ctx, cfg)
}
func (r *RetryProvider) Name() string      { return r.Base.Name() }
func (r *RetryProvider) ModelName() string { return r.Base.ModelName() }
func (r *RetryProvider) Close() error      { return r.Base.Close() }

func (r *RetryProvider) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	if r.MaxRetries < 0 {
		r.MaxRetries = 0
	}
	if r.Backoff <= 0 {
		r.Backoff = 100 * time.Millisecond
	}

	var lastErr error
	for attempt := 0; attempt <= r.MaxRetries; attempt++ {
		resp, err := r.Base.Complete(ctx, req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		if r.ShouldRetry != nil && !r.ShouldRetry(err) {
			return CompletionResponse{}, err
		}
		if attempt < r.MaxRetries {
			select {
			case <-ctx.Done():
				return CompletionResponse{}, ctx.Err()
			case <-time.After(r.Backoff):
			}
		}
	}
	return CompletionResponse{}, lastErr
}

func (r *RetryProvider) Stream(ctx context.Context, req CompletionRequest) (<-chan StreamEvent, error) {
	return r.Base.Stream(ctx, req)
}
