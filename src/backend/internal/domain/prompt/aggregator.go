package prompt

import "strings"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AggregationInput struct {
	AgentSystemPrompts []string
	GlobalPrompts      []*UserPrompt
	GroupPrompts       []*UserPrompt
	RAGContext         string
	PhaseInstruction   string
	UserInstruction    string
}

type PromptAggregator struct{}

func NewPromptAggregator() *PromptAggregator {
	return &PromptAggregator{}
}

func (a *PromptAggregator) Compose(input AggregationInput) []Message {
	messages := make([]Message, 0, 3)
	system := make([]string, 0)

	if len(input.AgentSystemPrompts) > 0 {
		system = append(system, strings.Join(input.AgentSystemPrompts, "\n\n"))
	}

	SortByPriority(input.GlobalPrompts)
	for _, p := range input.GlobalPrompts {
		if p != nil && p.Enabled {
			system = append(system, p.Content)
		}
	}

	SortByPriority(input.GroupPrompts)
	for _, p := range input.GroupPrompts {
		if p != nil && p.Enabled {
			system = append(system, p.Content)
		}
	}

	if len(system) > 0 {
		messages = append(messages, Message{Role: "system", Content: strings.Join(system, "\n\n")})
	}

	if strings.TrimSpace(input.RAGContext) != "" {
		messages = append(messages, Message{Role: "system", Content: "Contexto RAG (SPEC.md da fase anterior):\n" + strings.TrimSpace(input.RAGContext)})
	}

	userParts := make([]string, 0, 2)
	if strings.TrimSpace(input.PhaseInstruction) != "" {
		userParts = append(userParts, strings.TrimSpace(input.PhaseInstruction))
	}
	if strings.TrimSpace(input.UserInstruction) != "" {
		userParts = append(userParts, strings.TrimSpace(input.UserInstruction))
	}
	if len(userParts) > 0 {
		messages = append(messages, Message{Role: "user", Content: strings.Join(userParts, "\n\n")})
	}

	return messages
}
