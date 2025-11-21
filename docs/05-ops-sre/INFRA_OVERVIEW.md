> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# Infra Overview

- Banco principal: PostgreSQL no Neon (prod/staging), migrations via golang-migrate.
- App backend: Go 1.24.0 rodando via systemd em VPS Ubuntu + NGINX reverse proxy.
- Frontend: Next.js 16.0.3 (MUI) servindo via Vercel/VPS (dependendo do ambiente).
- Observabilidade: Prometheus/Grafana, Sentry (a integrar), logs JSON.
- Referências: GUIA_DEVOPS.md, BACKUP_DR.md, observability/.
