# ADR-005: Gerenciamento de Configuração com Viper

## Status
Aceito

## Contexto
A plataforma exige múltiplos parâmetros de conexão (Bancos, API Keys de LLMs, portas), devendo lidar tanto com valores de arquivos YAML para desenvolvimento, quanto variáveis de ambiente para produção via Docker.

## Decisão
Adoção da biblioteca **Viper** (`github.com/spf13/viper`).

## Consequências
- **Prós**: Suporte nativo a múltiplos formatos (YAML), substituição unificada via variáveis de ambiente, fácil leitura em Golang. Suporta defaults declarados.
- **Contras**: Uso de strings mágicas para o carregamento e reflexão, que não é "type-safe" durante a declaração estática, precisando que a aplicação carregue um struct validado (`config.go`) para contornar problemas.
