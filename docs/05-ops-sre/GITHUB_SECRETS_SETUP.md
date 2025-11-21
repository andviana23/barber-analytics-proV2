> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🔐 GitHub Secrets Configuration

## Secrets Necessários para CI/CD

Configure os seguintes secrets em **Settings → Secrets and variables → Actions → New repository secret**:

### 1. `NEON_DATABASE_URL` ⭐ (Obrigatório)

**Descrição:** URL de conexão com banco de dados Neon PostgreSQL

**Formato:**
```
postgresql://user:password@ep-xxx-xxx.us-east-1.aws.neon.tech/neondb?sslmode=require
```

**Como obter:**
1. Acesse [Neon Console](https://console.neon.tech)
2. Selecione seu projeto
3. Vá em **Connection Details**
4. Copie a **Connection String** (Pooled connection)
5. Remova o prefixo `psql '` e o sufixo `'` se houver

**Exemplo:**
```
postgresql://neondb_owner:npg_bH5euQYkf3iE@ep-winter-leaf-adhqz08p-pooler.us-east-1.aws.neon.tech/neondb?sslmode=require
```

**Verificação local:**
```bash
# Testar conexão
psql "postgresql://..." -c "SELECT version();"
```

---

### 2. `E2E_USER_EMAIL` (Opcional - tem default)

**Descrição:** Email do usuário de teste para E2E

**Default:** `qa@barberpro.dev`

**Uso:** Credenciais do seed criado por `scripts/seed_test_data.go`

---

### 3. `E2E_USER_PASSWORD` (Opcional - tem default)

**Descrição:** Senha do usuário de teste para E2E

**Default:** `Test@1234`

**Uso:** Deve corresponder ao seed de dados

---

### 4. `JWT_PRIVATE_KEY` e `JWT_PUBLIC_KEY` (Gerados automaticamente no CI)

**Descrição:** Chaves RSA para autenticação JWT

**Geração automática:** O workflow CI gera essas chaves automaticamente via OpenSSL:

```yaml
- name: Generate JWT keys
  run: |
    mkdir -p backend/keys
    openssl genrsa -out backend/keys/private.pem 2048
    openssl rsa -in backend/keys/private.pem -pubout -out backend/keys/public.pem
```

**Não é necessário configurar como secrets**, a menos que você queira usar chaves específicas.

---

## 📝 Como Configurar Secrets

### Via GitHub Web UI

1. Vá para o repositório no GitHub
2. **Settings** → **Secrets and variables** → **Actions**
3. Clique em **New repository secret**
4. Preencha:
   - **Name:** `NEON_DATABASE_URL`
   - **Secret:** Cole a URL do Neon
5. Clique em **Add secret**
6. Repita para outros secrets (se necessário)

### Via GitHub CLI

```bash
# Instalar gh CLI
brew install gh  # macOS
# ou: sudo apt install gh  # Ubuntu

# Autenticar
gh auth login

# Adicionar secrets
gh secret set NEON_DATABASE_URL -b "postgresql://..."
gh secret set E2E_USER_EMAIL -b "qa@barberpro.dev"
gh secret set E2E_USER_PASSWORD -b "Test@1234"

# Verificar secrets configurados
gh secret list
```

---

## ✅ Verificação

Após configurar os secrets, verifique no workflow:

1. Vá para **Actions** → **Frontend Tests**
2. Clique em **Run workflow** → **Run workflow**
3. Aguarde execução
4. Verifique logs do job `e2e-tests`:

```
✅ Setup Backend
✅ Running database migrations
✅ Seeding test data
✅ Starting backend server
✅ Backend is ready!
✅ Running E2E tests
```

---

## 🔍 Troubleshooting

### Secret não aparece no workflow

**Problema:** Workflow não consegue acessar `${{ secrets.NEON_DATABASE_URL }}`

**Soluções:**
1. Verificar se secret foi criado corretamente (Settings → Secrets)
2. Nome do secret deve ser **exatamente igual** (case-sensitive)
3. Re-executar workflow (pode ter cache)
4. Verificar logs do workflow para mensagens de erro

### Conexão com Neon falha

**Erro:**
```
error: failed to connect to database: FATAL: password authentication failed
```

**Soluções:**
1. Verificar se URL do Neon está correta
2. Testar conexão localmente: `psql "postgresql://..."`
3. Gerar nova senha no Neon Console
4. Atualizar secret com nova URL

### Migrations falham no CI

**Erro:**
```
error: migration failed: relation "users" already exists
```

**Soluções:**
1. Usar banco de teste separado (não produção!)
2. Adicionar flag `|| echo "Migrations already applied"` no workflow
3. Verificar se migrations estão na ordem correta

---

## 🚀 Workflow Completo

Após configurar secrets, o workflow executará:

```mermaid
graph TD
    A[Push/PR] --> B[unit-tests]
    B --> C[e2e-tests]
    C --> D{Backend Start}
    D -->|Success| E[Run E2E Tests]
    D -->|Fail| F[Upload Logs]
    E --> G[Upload Reports]
    B --> H[coverage-report]
```

---

## 📚 Referências

- [GitHub Actions Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Neon Documentation](https://neon.tech/docs/connect/connect-from-any-app)
- [Workflow File](./.github/workflows/frontend-tests.yml)

---

**Última atualização:** 15/11/2025
**Versão:** 1.0.0
