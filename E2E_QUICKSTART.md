# ⚡ Quick Start — E2E Tests

## 🚀 Execução Rápida (Recomendado)

```bash
# Tudo em um comando!
./scripts/run-e2e-tests.sh

# Com navegador visível
./scripts/run-e2e-tests.sh --headed

# Interface interativa
./scripts/run-e2e-tests.sh --ui
```

**O que o script faz automaticamente:**
1. ✅ Verifica pré-requisitos (Go, pnpm, DATABASE_URL)
2. ✅ Gera JWT keys se não existirem
3. ✅ Aplica migrations no Neon
4. ✅ Executa seed de dados de teste
5. ✅ Inicia backend em background
6. ✅ Aguarda backend estar pronto
7. ✅ Executa testes E2E
8. ✅ Para backend ao final (cleanup automático)

---

## 📝 Execução Manual (Passo a Passo)

### 1. Setup Inicial (uma vez)

```bash
# Backend - Gerar JWT keys
cd backend
mkdir -p keys
openssl genrsa -out keys/private.pem 2048
openssl rsa -in keys/private.pem -pubout -out keys/public.pem

# Backend - Aplicar migrations
export DATABASE_URL="postgresql://..."  # Neon URL
make migrate-up

# Backend - Criar dados de teste
go run scripts/seed_test_data.go
```

### 2. Executar Testes

```bash
# Terminal 1: Backend
cd backend
make run

# Terminal 2: Testes E2E (aguardar backend iniciar)
cd frontend
pnpm test:e2e
```

---

## 🔑 Credenciais de Teste

Após executar `seed_test_data.go`:

```
Email:    qa@barberpro.dev
Password: Test@1234
Tenant:   e2e-test-barbershop
```

---

## 🧪 Comandos de Teste

```bash
# Testes unitários
pnpm test:unit

# Testes de acessibilidade
pnpm test:a11y

# E2E - todos (headless)
pnpm test:e2e

# E2E - navegador visível
pnpm test:e2e:headed

# E2E - interface interativa
pnpm test:e2e:ui

# E2E - debug com breakpoints
pnpm test:e2e:debug

# Testes específicos
pnpm test:e2e -- e2e/login.spec.ts
pnpm test:e2e -- e2e/receitas.spec.ts
```

---

## 🐛 Troubleshooting

### Backend não inicia

```bash
# Verificar DATABASE_URL
echo $DATABASE_URL

# Testar conexão
curl http://localhost:8080/health
```

### Dados de teste não encontrados

```bash
# Re-executar seed
cd backend
go run scripts/seed_test_data.go
```

### Testes falham

```bash
# Ver logs do backend
cat /tmp/backend-e2e.log

# Ver relatório do Playwright
open frontend/playwright-report/index.html
```

---

## 📚 Documentação Completa

- [E2E Testing Guide](../frontend/docs/E2E_TESTING_GUIDE.md) — Guia completo
- [Frontend README](../frontend/README.md) — Setup geral
- [FASE_4_FRONTEND.md](../Tarefas/FASE_4_FRONTEND.md) — Tarefas e progresso

---

**Versão:** 1.0.0
**Data:** 15/11/2025
**Autor:** Andrey Viana
