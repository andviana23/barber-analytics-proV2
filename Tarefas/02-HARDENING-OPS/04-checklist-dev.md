# ✅ Checklist Dev — Hardening & OPS

- [ ] Endpoints LGPD criados com DTOs e validação; testam ownership do usuário logado (tenant + user_id).
- [ ] Exclusão lógica limpa tokens/sessions e agenda job para anonimizar dados relacionados (logs, assinaturas, invoices).
- [ ] Export retorna JSON completo e faz stream (evitar uso excessivo de memória).
- [ ] Banner/página `/privacy` consumindo preferências via hooks; consentimento armazenado em `user_preferences`.
- [ ] Backup workflow com variáveis seguras (DATABASE_URL, AWS creds) e artefatos versionados.
- [ ] Script/guide de restore validado com banco restaurado em staging.
- [ ] Alertas configurados para falha de backup, espaço S3, restore não testado (<30 dias).
- [ ] Métricas Prometheus para endpoints LGPD (latência, taxa de erro) adicionadas.
