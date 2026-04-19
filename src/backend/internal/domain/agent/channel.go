package agent

import "time"

// AgentMessage é o envelope JSON padrão para toda comunicação entre agentes.
// A comunicação direta entre agentes é proibida — todo dado deve passar por este tipo.
type AgentMessage struct {
	ID        string         `json:"id"`
	From      string         `json:"from"`
	To        string         `json:"to"`
	Message   string         `json:"message"`
	Status    string         `json:"status"`
	Timestamp time.Time      `json:"timestamp"`
	Meta      map[string]any `json:"meta,omitempty"`
}

// AgentChannel encapsula os canais de entrada e saída de um agente.
// Criado automaticamente quando o agente é instanciado;
// fechado e removido automaticamente quando o agente é encerrado.
type AgentChannel struct {
	In  <-chan AgentMessage
	Out chan<- AgentMessage
}
