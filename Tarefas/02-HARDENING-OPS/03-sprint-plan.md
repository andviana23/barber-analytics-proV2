# 🗓️ Plano de Sprint — Hardening & OPS

1. [ ] Implementar stack LGPD (T-HAR-001) — 8h
   - Backend endpoints + ajustes no tenant/user context.
   - Frontend banner + página `/privacy`.
2. [ ] Configurar Backup/DR (T-HAR-002) — 6h
   - Workflow GH Actions, S3, PITR, teste de restore.
3. [ ] Regressão/observabilidade (T-HAR-003) — 4h
   - Segurança, métricas, alertas, runbook atualizado.

**Gates:**
- Não iniciar Financeiro avançado sem `DELETE /me` + `GET /me/export` funcionando e com auditoria.
- Falha em restore bloqueia go-live; repetir até sucesso documentado.
