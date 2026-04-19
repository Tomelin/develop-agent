package prompt

import (
	"testing"
	"time"
)

func TestPromptAggregatorComposeOrder(t *testing.T) {
	agg := NewPromptAggregator()
	now := time.Now().UTC()
	messages := agg.Compose(AggregationInput{
		AgentSystemPrompts: []string{"agent identity", "agent style"},
		GlobalPrompts: []*UserPrompt{
			{Content: "global 2", Enabled: true, Priority: 2, CreatedAt: now.Add(2 * time.Second)},
			{Content: "global 1", Enabled: true, Priority: 1, CreatedAt: now.Add(1 * time.Second)},
		},
		GroupPrompts: []*UserPrompt{
			{Content: "group 1", Enabled: true, Priority: 1, CreatedAt: now.Add(3 * time.Second)},
		},
		RAGContext:       "spec previous phase",
		PhaseInstruction: "execute phase",
	})

	if len(messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(messages))
	}
	if messages[0].Role != "system" || messages[1].Role != "system" || messages[2].Role != "user" {
		t.Fatalf("unexpected roles order: %#v", messages)
	}
	expected := "agent identity\n\nagent style\n\nglobal 1\n\nglobal 2\n\ngroup 1"
	if messages[0].Content != expected {
		t.Fatalf("unexpected system composition:\n%s", messages[0].Content)
	}
}

func TestPromptAggregatorComposeEdgeCases(t *testing.T) {
	agg := NewPromptAggregator()

	cases := []struct {
		name  string
		input AggregationInput
		want  int
	}{
		{name: "without global", input: AggregationInput{AgentSystemPrompts: []string{"base"}, GroupPrompts: []*UserPrompt{{Content: "group", Enabled: true}}, PhaseInstruction: "go"}, want: 2},
		{name: "without group", input: AggregationInput{AgentSystemPrompts: []string{"base"}, GlobalPrompts: []*UserPrompt{{Content: "global", Enabled: true}}, PhaseInstruction: "go"}, want: 2},
		{name: "without rag", input: AggregationInput{AgentSystemPrompts: []string{"base"}, PhaseInstruction: "go"}, want: 2},
		{name: "without user prompts", input: AggregationInput{AgentSystemPrompts: []string{"base"}, PhaseInstruction: "go"}, want: 2},
		{name: "max prompt token estimate", input: AggregationInput{AgentSystemPrompts: []string{"base"}, GroupPrompts: []*UserPrompt{{Content: "a b c d", Enabled: true}}, PhaseInstruction: "go"}, want: 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := agg.Compose(tc.input)
			if len(got) != tc.want {
				t.Fatalf("expected %d messages, got %d", tc.want, len(got))
			}
		})
	}

	if tokens := EstimateTokens("a b c d e f g h"); tokens <= 0 {
		t.Fatalf("expected positive token estimate, got %d", tokens)
	}
}
