#!/bin/bash

# Este script não inicia mais backend/frontend automaticamente.
# Ele apenas mostra os comandos recomendados para rodar cada serviço
# de forma independente em terminais separados.

PROJECT_ROOT="/home/andrey/projetos/barber-Analytic-proV2"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

cat <<EOF
╔════════════════════════════════════════════════════════════╗
║        Dev Mode — Backend e Frontend iniciados manualmente  ║
╚════════════════════════════════════════════════════════════╝

Use os comandos abaixo em terminais separados.

1) Backend (Go)
------------------------------------------------------------
cd "$BACKEND_DIR"

# Ajuste DATABASE_URL conforme o ambiente (Neon, local etc.)
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/barber_db?sslmode=disable"
export DEV_MODE="true"
export PORT="8080"

go run cmd/api/main.go


2) Frontend (Next.js)
------------------------------------------------------------
cd "$FRONTEND_DIR"

export NEXT_PUBLIC_API_URL="http://localhost:8080/api/v1"

pnpm install
pnpm dev


3) Testar login
------------------------------------------------------------
Acesse: http://localhost:3000/login
Use as credenciais documentadas em LOGIN_IMPLEMENTADO.md.

EOF
