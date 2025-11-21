> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# Arquitetura de Segurança

- Auth: JWT RS256, refresh tokens.
- RBAC: roles Owner/Manager/Accountant/Employee, enforcement em middleware.
- Auditoria: audit logs para CREATE/UPDATE/DELETE, retenção e índices.
- Rate limiting: middleware + NGINX.
- Banco: Neon PostgreSQL com RLS (tenant_id), backups automáticos.
- Referências: RBAC.md, AUDIT_LOGS.md, SECURITY_TESTING.md, COMPLIANCE_LGPD.md.
