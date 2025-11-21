> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 📝 Próximos Passos — Implementação Completa

## ✅ Etapas Concluídas

### 1. Infraestrutura de Testes E2E
- [x] Script `seed_test_data.go` para dados de teste
- [x] Configuração `.env.test` no frontend
- [x] Global setup (aguarda backend + valida seed)
- [x] Global teardown (cleanup automático)
- [x] Auth fixture reutilizável
- [x] Playwright config atualizado
- [x] CI/CD workflow atualizado (GitHub Actions)

### 2. Testes E2E Implementados
- [x] Login flow (3 testes)
- [x] Receitas CRUD (6 testes)
- [x] Despesas CRUD (8 testes) — **NOVO**
- [x] Assinaturas management (9 testes) — **NOVO**

**Total:** 26 cenários E2E implementados

### 3. Documentação
- [x] E2E Testing Guide (guia completo)
- [x] E2E Quickstart (referência rápida)
- [x] GitHub Secrets Setup (configuração CI/CD)
- [x] Script automatizado `run-e2e-tests.sh`
- [x] README atualizado

---

## 🚀 Execução Local (Passo a Passo)

### Método 1: Script Automatizado (Recomendado)

```bash
# Tudo em um comando!
./scripts/run-e2e-tests.sh

# Com navegador visível
./scripts/run-e2e-tests.sh --headed

# Interface interativa
./scripts/run-e2e-tests.sh --ui

# Pular seed (usar dados existentes)
./scripts/run-e2e-tests.sh --skip-seed
```

### Método 2: Manual

#### Passo 1: Backend

```bash
# Terminal 1: Backend
cd backend

# Configurar DATABASE_URL
export DATABASE_URL="postgresql://neondb_owner:***@ep-winter-leaf-*.neon.tech/neondb?sslmode=require"

# Gerar JWT keys (primeira vez)
mkdir -p keys
openssl genrsa -out keys/private.pem 2048
openssl rsa -in keys/private.pem -pubout -out keys/public.pem

# Aplicar migrations
make migrate-up

# Seed de dados de teste
go run scripts/seed_test_data.go

# Iniciar backend
make run
# Ou: go run cmd/api/main.go
```

**Verificar backend:**
```bash
curl http://localhost:8080/health
# Resposta: {"status":"healthy","timestamp":"..."}
```

#### Passo 2: Frontend - Testes Unitários

```bash
# Terminal 2: Frontend
cd frontend

# Testes unitários
pnpm test:unit
# Resultado esperado: 67/67 passing

# Testes de acessibilidade
pnpm test:a11y
# Resultado esperado: 25/25 passing
```

#### Passo 3: Frontend - Testes E2E

```bash
# Executar todos os testes E2E
pnpm test:e2e

# Ou testes específicos
pnpm test:e2e -- e2e/login.spec.ts
pnpm test:e2e -- e2e/receitas.spec.ts
pnpm test:e2e -- e2e/despesas.spec.ts
pnpm test:e2e -- e2e/assinaturas.spec.ts

# Com navegador visível
pnpm test:e2e:headed

# Interface interativa
pnpm test:e2e:ui
```

**Resultados esperados:**
```
Running 26 tests using 3 workers

✓ Login Flow (3 testes)
✓ Receitas CRUD (6 testes)
✓ Despesas CRUD (8 testes)
✓ Assinaturas Management (9 testes)

26 passed (90-120s)
```

---

## 🔐 Configuração GitHub Secrets

### Secrets Obrigatórios

1. **`NEON_DATABASE_URL`** ⭐ (Obrigatório)
   ```
   postgresql://user:password@ep-xxx.neon.tech/neondb?sslmode=require
   ```

2. **`E2E_USER_EMAIL`** (Opcional - default: `qa@barberpro.dev`)

3. **`E2E_USER_PASSWORD`** (Opcional - default: `Test@1234`)

### Como Configurar

#### Via GitHub Web

1. Repositório → **Settings**
2. **Secrets and variables** → **Actions**
3. **New repository secret**
4. Nome: `NEON_DATABASE_URL`
5. Value: Cole a URL do Neon
6. **Add secret**

#### Via GitHub CLI

```bash
gh secret set NEON_DATABASE_URL -b "postgresql://..."
gh secret set E2E_USER_EMAIL -b "qa@barberpro.dev"
gh secret set E2E_USER_PASSWORD -b "Test@1234"

# Verificar
gh secret list
```

**Documentação completa:** [`docs/GITHUB_SECRETS_SETUP.md`](../docs/GITHUB_SECRETS_SETUP.md)

---

## 🔄 Executar Workflow CI

### Disparar Manualmente

1. GitHub → **Actions** → **Frontend Tests**
2. **Run workflow** → **Run workflow**
3. Aguardar execução

### Disparar via Push/PR

```bash
# Qualquer push para develop/main dispara workflow
git add .
git commit -m "feat: add E2E tests for despesas and assinaturas"
git push origin develop
```

### Monitorar Execução

**Jobs esperados:**
1. ✅ `unit-tests` → lint + unit + a11y (67 + 25 testes)
2. ✅ `e2e-tests` → backend + seed + 26 testes E2E
3. ✅ `coverage-report` → upload Codecov

**Duração estimada:** 8-12 minutos

---

## 🧹 Cleanup Automático

### Configuração

**Arquivo:** `frontend/.env.test`

```env
# Habilitar cleanup automático após testes
E2E_AUTO_CLEANUP=true  # false para desabilitar
```

**Como funciona:**

1. **Global teardown** executa após todos os testes
2. Busca registros com prefixo "E2E Test"
3. Deleta receitas, despesas e assinaturas de teste
4. Remove arquivos temporários antigos

**Logs esperados:**
```
🏁 Playwright Global Teardown

🧹 Limpando dados de teste...
   ✅ Deletadas 5 receitas de teste
   ✅ Deletadas 3 despesas de teste
✅ Cleanup de dados concluído

🗑️  Removendo arquivos temporários...
   ✅ Removido storage state antigo (>24h)
   ✅ Removidos 2 diretórios de screenshots antigos
✅ Limpeza de arquivos concluída

✅ Global Teardown concluído com sucesso!
```

### Desabilitar Cleanup

Para manter dados de teste no banco (útil para debug):

```env
E2E_AUTO_CLEANUP=false
```

---

## 📊 Sumário de Testes

### Testes Unitários
| Componente | Testes | Status |
|------------|--------|--------|
| Button | 30 | ✅ |
| AccessibleInput | 18 | ✅ |
| Modal | 19 | ✅ |
| **Total** | **67** | **✅** |

### Testes de Acessibilidade
| Componente | Testes | Violations |
|------------|--------|------------|
| Button | 8 | 0 |
| AccessibleInput | 9 | 0 |
| Modal | 8 | 0 |
| **Total** | **25** | **0** |

### Testes E2E
| Módulo | Cenários | Status |
|--------|----------|--------|
| Login | 3 | ✅ |
| Receitas | 6 | ✅ |
| Despesas | 8 | ✅ |
| Assinaturas | 9 | ✅ |
| **Total** | **26** | **✅** |

**Total Geral:** 118 testes automatizados

---

## 🎯 Melhorias Futuras

### Curto Prazo
- [ ] Adicionar data-testid em componentes faltantes
- [ ] Implementar endpoint `/tests/cleanup` no backend
- [ ] Adicionar testes de performance (Lighthouse)
- [ ] Configurar testes cross-browser (Firefox, Safari)

### Médio Prazo
- [ ] Testes E2E de agendamentos
- [ ] Testes E2E de estoque
- [ ] Testes de integração Asaas (webhook)
- [ ] Testes visuais com Percy

### Longo Prazo
- [ ] Testes de carga (k6)
- [ ] Testes de segurança (OWASP ZAP)
- [ ] Cobertura E2E > 80%
- [ ] Integração com Sonar Cloud

---

## 📚 Documentação de Referência

| Documento | Descrição |
|-----------|-----------|
| [E2E Testing Guide](../frontend/docs/E2E_TESTING_GUIDE.md) | Guia completo (300+ linhas) |
| [E2E Quickstart](../E2E_QUICKSTART.md) | Referência rápida |
| [GitHub Secrets Setup](../docs/GITHUB_SECRETS_SETUP.md) | Configuração CI/CD |
| [Frontend README](../frontend/README.md) | Setup geral |
| [FASE_4_FRONTEND](../Tarefas/FASE_4_FRONTEND.md) | Progresso de tarefas |

---

## 🐛 Troubleshooting Comum

### Backend não inicia
```bash
# Verificar DATABASE_URL
echo $DATABASE_URL

# Testar conexão
psql "$DATABASE_URL" -c "SELECT 1"

# Ver logs
tail -f /tmp/backend-e2e.log
```

### Testes falham com timeout
```bash
# Aumentar timeout em .env.test
BACKEND_HEALTH_TIMEOUT=60000

# Ou aumentar timeout no teste
await page.waitForURL('/dashboard', { timeout: 30000 });
```

### Seed falha
```bash
# Re-executar seed
cd backend
go run scripts/seed_test_data.go

# Verificar no banco
psql "$DATABASE_URL" -c "SELECT * FROM users WHERE email = 'qa@barberpro.dev'"
```

---

## ✅ Checklist Final

Antes de considerar E2E completo, verificar:

- [x] Backend inicia localmente sem erros
- [x] Seed cria dados idempotentemente
- [x] Testes unitários passam (67/67)
- [x] Testes a11y passam (25/25)
- [x] Testes E2E passam localmente (26/26)
- [x] Secrets configurados no GitHub
- [x] Workflow CI executa sem erros
- [x] Cleanup automático funciona
- [x] Documentação completa e atualizada

---

**Status:** ✅ 100% Completo
**Data:** 15/11/2025
**Autor:** Andrey Viana
