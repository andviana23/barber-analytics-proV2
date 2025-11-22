# 🔎 Contexto — Hardening & OPS

## Estado atual
- Fase 6 ≈85%: métricas, dashboards, alertas, rate limiting e caching **concluídos** (`Tarefas/FASE_6_HARDENING.md`).
- LGPD: migrations já criaram `user_preferences` e `users.deleted_at`, mas endpoints `/me/preferences`, `/me/export`, `/me` (delete) e banner de consentimento não existem.
- Backup/DR: documentação pronta, mas não há pipeline executando backup nem teste de restore (`T-OPS-005` pendente).

## Riscos
- Sem endpoints LGPD, não atendemos requisições de titulares e violamos compliance.
- Sem rotina de backup + teste de restore, o risco de perda de dados inviabiliza go-live.

## Objetivo desta pasta
Fechar pendências críticas de LGPD e de resiliência operacional para liberar o roadmap a partir da Fase 7.

## Referências
- `Tarefas/FASE_6_HARDENING.md`
- `docs/COMPLIANCE_LGPD.md`
- `docs/BACKUP_DR.md`
- `docs/DIAGRAMA_DEPENDENCIAS_COMPLETO.md`
