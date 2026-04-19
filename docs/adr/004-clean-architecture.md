# ADR-004: Clean Architecture no Backend

## Status
Aceito

## Contexto
O projeto precisa ser de fácil manutenção e testável, e ser agnóstico aos detalhes de banco de dados e APIs externas (OpenAI, Anthropic). O código não pode misturar regras de negócio da "Agência de IA" com a infraestrutura.

## Decisão
Adoção de **Clean Architecture** (Arquitetura Limpa). O projeto Go terá a estrutura dividida em:
- `domain`: Regras de negócio, interfaces e tipos isolados.
- `usecase`: Casos de uso da aplicação e orquestração.
- `infra`: Adapters e integrações com bancos de dados, redis, AMQP, LLMs.
- `api`: Controladores HTTP, GIN.

## Consequências
- **Prós**: Desacoplamento, testabilidade alta através de interfaces (Mocks), facilidade de substituir componentes futuramente (ex: trocar MongoDB por PostgreSQL, ou mudar o provedor de LLM sem afetar o core business).
- **Contras**: Maior complexidade inicial com múltiplas camadas e arquivos para a mesma entidade. Curva de aprendizado inicial maior.
