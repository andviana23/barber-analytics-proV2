# 📌 Backlog — Hardening & OPS

## 🔴 Obrigatórios
1. [ ] **T-HAR-001 — LGPD Compliance End-to-End** — ref. `Tarefas/FASE_6_HARDENING.md`
   - Endpoints: `GET/PUT /me/preferences`, `GET /me/export`, `DELETE /me` com deleção lógica (`users.deleted_at`) + scrub de PII.
   - Banner/página `/privacy` no frontend + registro de consentimento granular (necessário vs opcional) em `user_preferences`.
   - Logs de auditoria em toda operação LGPD e runbook para requisições de titulares.
2. [ ] **T-HAR-002 — Backup & Disaster Recovery (T-OPS-005)**
   - Workflow GitHub Actions: `pg_dump` do Neon, upload para S3 com versionamento, retenção e criptografia.
   - PITR configurado no Neon + teste de restore em staging documentado.
   - Alertas no Prometheus/Alertmanager para falha de backup e storage.
3. [ ] **T-HAR-003 — Validação final de segurança/observabilidade**
   - Revisar que novos endpoints LGPD possuem rate limiting, RBAC, métricas e alertas.
   - Documentar decisão de manter Sentry como skip (T-OPS-003) e garantir que stack Prometheus/Grafana cobre erros críticos.

## 🧭 Dependências
- Requer domínio e handlers prontos (`01-BLOQUEIOS-BASE`) para publicar endpoints.
- Usar `DATABASE_MIGRATIONS_COMPLETED.md` para validar colunas (`deleted_at`, `user_preferences`).
