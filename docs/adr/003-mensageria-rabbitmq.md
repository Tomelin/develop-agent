# ADR-003: Mensageria Assíncrona com RabbitMQ

## Status
Aceito

## Contexto
As Fases do sistema (1 ao 9) processam tarefas de maneira assíncrona com os Modelos de Linguagem (LLMs). As fases precisam se comunicar, e se algo falhar (ex: rejeição automática do LLM), uma fase pode acionar outra fase em background (DLQ e retries).

## Decisão
Adoção do **RabbitMQ** para a camada de mensageria da plataforma.

## Consequências
- **Prós**: Traz flexibilidade de roteamento (`topic` exchanges), fácil visualização com o Management Plugin, bom controle de ack/nack explícito. Ideal para orquestração de tarefas que demoram alguns minutos (requests para LLMs).
- **Contras**: Diferente do Kafka que atua como log append-only e retenção de eventos em alta escala; porém, a nossa prioridade é processamento de workers de tarefas ativas, para o que o RabbitMQ se sobressai.
