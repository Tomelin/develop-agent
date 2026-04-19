package agentsdk

import "context"

// BillingSink receives usage accounting events.
type BillingSink interface {
	Track(provider string, model string, usage Usage)
}

// BillingProvider decorates a provider and tracks usage.
type BillingProvider struct {
	Base Provider
	Sink BillingSink
}

func (b *BillingProvider) Initialize(ctx context.Context, cfg Config) error {
	return b.Base.Initialize(ctx, cfg)
}
func (b *BillingProvider) Name() string      { return b.Base.Name() }
func (b *BillingProvider) ModelName() string { return b.Base.ModelName() }
func (b *BillingProvider) Close() error      { return b.Base.Close() }

func (b *BillingProvider) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	resp, err := b.Base.Complete(ctx, req)
	if err != nil {
		return CompletionResponse{}, err
	}
	if b.Sink != nil {
		b.Sink.Track(b.Base.Name(), b.Base.ModelName(), resp.Usage)
	}
	return resp, nil
}

func (b *BillingProvider) Stream(ctx context.Context, req CompletionRequest) (<-chan StreamEvent, error) {
	return b.Base.Stream(ctx, req)
}
