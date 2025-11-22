# 🧪 Checklist QA — Hardening & OPS

- [ ] Testar `/me/preferences` com roles diferentes (owner/employee) e verificar isolamento por tenant.
- [ ] Solicitar exportação e validar JSON completo sem campos vazios/corrompidos.
- [ ] Solicitar deleção e confirmar `users.deleted_at` preenchido + remoção/anonimização nos demais registros.
- [ ] Banner de consentimento respeita escolhas e permite revogação; preferências persistem após reload.
- [ ] Executar pipeline de backup manual e checar artefato no S3 (tamanho, checksum).
- [ ] Restaurar backup em staging e rodar `scripts/validate_schema.sh` + smoke tests.
- [ ] Verificar alertas disparando para falha de backup (simular) e ausência de restore (>30 dias).
- [ ] Regressão de segurança: SQLi/XSS/CSRF/RBAC continuam passando (35/35 testes).
