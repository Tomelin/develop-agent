# ADR-002: Banco de Dados Principal MongoDB

## Status
Aceito

## Contexto
A plataforma de Agência de IA terá diferentes fluxos, tarefas, prompts e metadados que podem ter formatos variados dependendo do fluxo (Desenvolvimento, Marketing, Landing Pages). A persistência dos projetos e pipelines precisa ser flexível.

## Decisão
Adoção do **MongoDB** como banco de dados principal.

## Consequências
- **Prós**: Natureza schemaless que facilita o armazenamento de outputs heterogêneos de modelos de linguagem (JSONs dinâmicos). Fácil integração em Golang com o driver oficial. Boa escalabilidade horizontal.
- **Contras**: Ausência de validações estritas a nível de banco de dados (esquemas SQL), exigindo forte responsabilidade na camada de aplicação (Clean Architecture) para assegurar a consistência dos dados.
