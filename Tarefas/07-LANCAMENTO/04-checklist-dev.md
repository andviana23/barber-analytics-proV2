# ✅ Checklist Dev — Lançamento

- [ ] Todas as migrations aplicadas e verificadas via `scripts/validate_schema.sh`.
- [ ] Suites de testes: backend (`go test ./...`), frontend (`pnpm test:*`), segurança (35/35), carga (k6) passando.
- [ ] Pipelines de deploy configurados com rollback e health checks.
- [ ] Feature flags configuradas para rollout gradual se necessário.
- [ ] Documentação de release e mudanças publicada.
