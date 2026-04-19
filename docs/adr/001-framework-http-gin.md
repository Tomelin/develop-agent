# ADR-001: Framework HTTP GIN

## Status
Aceito

## Contexto
O projeto requer um servidor HTTP de alta performance para a camada de backend (desenvolvido em Golang). Existem várias opções maduras como `net/http` padrão, Echo, Fiber e GIN. A escolha precisa balancear performance, ecossistema e familiaridade para uso em APIs REST.

## Decisão
Adoção do framework **GIN** (`github.com/gin-gonic/gin`) para o servidor HTTP.

## Consequências
- **Prós**: Alta performance, roteamento rápido baseado em radix tree, middlewares padrão da comunidade integrados (como Recovery, CORS, Logger), e facilidade de manipulação de JSON.
- **Contras**: Leve overhead em comparação com o `net/http` nativo puro, porém a curva de produtividade compensa.
