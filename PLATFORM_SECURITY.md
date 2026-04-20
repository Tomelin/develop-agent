# Platform Security Review

## Escopo
Revisão de segurança da plataforma da Agência de IA (infra + backend), sem considerar aplicações geradas para clientes.

## Resultado Consolidado
- **Autenticação:** JWT (RS256) e rotação de refresh token validados.
- **Autorização:** Rotas administrativas protegidas por role `ADMIN`.
- **API Security:** Rate limiting no login e tratamento defensivo de payload inválido.
- **Infra Security:** Serviços de dados previstos para rede interna com autenticação.
- **Dependências:** Nenhuma vulnerabilidade crítica conhecida no baseline desta revisão.

## Checklist Operacional
- [x] Invalidar sessão no logout.
- [x] Restringir endpoints administrativos por middleware de role.
- [x] Evitar exposição de stack traces para usuários finais.
- [x] Restringir acesso inter-usuário a recursos por `owner_user_id`.

## Observações
Esta revisão é contínua e deve ser repetida a cada release de produção.
