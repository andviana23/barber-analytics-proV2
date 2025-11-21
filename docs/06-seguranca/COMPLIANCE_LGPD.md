> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🛡️ LGPD Compliance — Barber Analytics Pro

**Documento de Conformidade à Lei Geral de Proteção de Dados (LGPD)**
**Versão:** 1.0.0
**Data:** 15/11/2025
**Status:** 🟡 Em Implementação

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Bases Legais](#bases-legais)
3. [Dados Pessoais Tratados](#dados-pessoais-tratados)
4. [Fluxo de Consentimento](#fluxo-de-consentimento)
5. [Direitos do Titular](#direitos-do-titular)
6. [Política de Retenção](#política-de-retenção)
7. [Segurança e Controles](#segurança-e-controles)
8. [Contatos e DPO](#contatos-e-dpo)

---

## 🎯 Visão Geral

O **Barber Analytics Pro** é um sistema SaaS multi-tenant para gestão de barbearias que trata dados pessoais de:
- **Titulares primários:** Proprietários e funcionários de barbearias (users)
- **Titulares secundários:** Clientes finais das barbearias (via assinaturas)

Este documento descreve como garantimos conformidade com a LGPD (Lei nº 13.709/2018).

---

## ⚖️ Bases Legais

Tratamos dados pessoais com base nas seguintes hipóteses legais:

### 1. **Execução de Contrato (Art. 7º, V)**
- **Aplicável a:** Cadastro de usuários, tenants, assinaturas
- **Finalidade:** Prover o serviço SaaS contratado
- **Dados:** Nome, email, CNPJ, telefone, endereço

### 2. **Legítimo Interesse (Art. 7º, IX)**
- **Aplicável a:** Logs de auditoria, métricas de uso, analytics
- **Finalidade:** Segurança, prevenção de fraudes, melhoria do serviço
- **Dados:** IP address, user agent, timestamps, ações realizadas
- **Balanceamento:** Interesse legítimo não sobrepõe direitos do titular

### 3. **Consentimento (Art. 7º, I)**
- **Aplicável a:** Cookies não essenciais, error tracking (Sentry), marketing
- **Finalidade:** Analytics, personalização, comunicação comercial
- **Dados:** Preferências, comportamento de navegação
- **Forma:** Banner de consentimento com opção de aceitar/rejeitar

### 4. **Cumprimento de Obrigação Legal (Art. 7º, II)**
- **Aplicável a:** Dados fiscais, tributários, trabalhistas
- **Finalidade:** Conformidade com legislação brasileira
- **Dados:** CNPJ, notas fiscais, folha de pagamento

---

## 📊 Dados Pessoais Tratados

### Inventário de Dados

| Categoria | Dados Coletados | Finalidade | Base Legal | Retenção |
|-----------|----------------|------------|------------|----------|
| **Usuários** | Nome, email, senha (hash), role | Autenticação e autorização | Contrato | Até exclusão da conta |
| **Tenants** | Nome da barbearia, CNPJ, telefone | Identificação do tenant | Contrato | Até cancelamento do plano |
| **Logs** | IP, user agent, timestamp, ação | Segurança e auditoria | Legítimo interesse | 90 dias |
| **Audit Logs** | UserID, ação, old/new values | Rastreabilidade LGPD | Legítimo interesse | 90 dias |
| **Assinaturas** | Nome do cliente, email, telefone | Clube do Trato | Contrato | Até cancelamento |
| **Analytics** | Pageviews, cliques, tempo de sessão | Melhoria do produto | Consentimento | Enquanto consentir |
| **Error Tracking** | Stack traces, request context | Debugging | Consentimento | 30 dias |

### Dados Sensíveis

**NÃO tratamos dados sensíveis** conforme Art. 5º, II da LGPD:
- ❌ Origem racial ou étnica
- ❌ Convicção religiosa
- ❌ Opinião política
- ❌ Filiação sindical
- ❌ Dados genéticos/biométricos
- ❌ Dados de saúde
- ❌ Vida sexual

---

## 🍪 Fluxo de Consentimento

### 1. Banner de Consentimento (Frontend)

**Implementação:** Modal no primeiro acesso

```typescript
// Exemplo: Componente CookieConsent.tsx
{
  "necessarios": {
    "enabled": true,
    "description": "Cookies essenciais para funcionamento (auth, sessão)",
    "optional": false
  },
  "analytics": {
    "enabled": false,
    "description": "Google Analytics para melhorar experiência",
    "optional": true
  },
  "error_tracking": {
    "enabled": false,
    "description": "Sentry para detectar erros e bugs",
    "optional": true
  }
}
```

**Opções:**
- ✅ **Aceitar todos** — Habilita todos os cookies
- ❌ **Rejeitar opcionais** — Apenas cookies essenciais
- ⚙️ **Gerenciar preferências** — Granularidade por categoria

### 2. Persistência de Preferências

**Frontend:**
```javascript
// localStorage
localStorage.setItem('cookie_preferences', JSON.stringify({
  version: '1.0',
  timestamp: Date.now(),
  analytics: true,
  error_tracking: false
}));
```

**Backend (opcional):**
```sql
-- Tabela user_preferences
CREATE TABLE user_preferences (
  user_id UUID PRIMARY KEY REFERENCES users(id),
  analytics_enabled BOOLEAN DEFAULT false,
  error_tracking_enabled BOOLEAN DEFAULT false,
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 3. Respeitar Consentimento

**Exemplo: Sentry**
```typescript
// Inicializar apenas se consentimento
if (preferences.error_tracking) {
  Sentry.init({ dsn: '...' });
}
```

**Exemplo: Google Analytics**
```html
<!-- Não carregar script se não consentir -->
<script>
  if (localStorage.getItem('analytics_enabled') === 'true') {
    // Carregar gtag.js
  }
</script>
```

---

## 👤 Direitos do Titular

### 1. **Acesso aos Dados (Art. 18, II)**

**Endpoint:** `GET /api/v1/me`

**Resposta:**
```json
{
  "user": {
    "id": "uuid",
    "nome": "João Silva",
    "email": "joao@barberpro.dev",
    "role": "owner",
    "criado_em": "2025-01-15T10:00:00Z"
  },
  "tenant": {
    "id": "uuid",
    "nome": "Barbearia Silva",
    "cnpj": "12345678000190"
  },
  "preferences": {
    "analytics_enabled": true,
    "error_tracking_enabled": false
  }
}
```

---

### 2. **Portabilidade (Art. 18, V)**

**Endpoint:** `GET /api/v1/me/export`

**Implementação:**
```go
// internal/application/usecase/user/export_data_usecase.go
type ExportDataUseCase struct {
    userRepo       domain.UserRepository
    tenantRepo     domain.TenantRepository
    receitaRepo    domain.ReceitaRepository
    despesaRepo    domain.DespesaRepository
    assinaturaRepo domain.AssinaturaRepository
}

func (uc *ExportDataUseCase) Execute(ctx context.Context, userID string) (*ExportDataResponse, error) {
    // Buscar todos os dados do usuário
    user, _ := uc.userRepo.FindByID(ctx, userID)
    tenant, _ := uc.tenantRepo.FindByID(ctx, user.TenantID)
    receitas, _ := uc.receitaRepo.FindByTenant(ctx, user.TenantID)
    despesas, _ := uc.despesaRepo.FindByTenant(ctx, user.TenantID)
    assinaturas, _ := uc.assinaturaRepo.FindByTenant(ctx, user.TenantID)

    return &ExportDataResponse{
        User:         user,
        Tenant:       tenant,
        Receitas:     receitas,
        Despesas:     despesas,
        Assinaturas:  assinaturas,
        ExportedAt:   time.Now(),
    }, nil
}
```

**Handler:**
```go
// GET /api/v1/me/export
func (h *UserHandler) ExportData(c echo.Context) error {
    userID := middleware.GetUserIDFromContext(c)

    data, err := h.exportDataUseCase.Execute(c.Request().Context(), userID)
    if err != nil {
        return echo.NewHTTPError(500, "Erro ao exportar dados")
    }

    // Retornar JSON ou ZIP se grande
    c.Response().Header().Set("Content-Disposition", "attachment; filename=meus-dados.json")
    return c.JSON(200, data)
}
```

**Proteções:**
- ✅ Rate limiting: 1 export/dia por usuário
- ✅ Logs de auditoria: Registrar exports
- ✅ Excluir segredos: Senhas, tokens, chaves API

---

### 3. **Exclusão/Esquecimento (Art. 18, VI)**

**Endpoint:** `DELETE /api/v1/me`

**Implementação:**
```go
// internal/application/usecase/user/delete_account_usecase.go
type DeleteAccountUseCase struct {
    userRepo      domain.UserRepository
    jwtService    domain.JWTService
    auditService  *audit.AuditService
}

func (uc *DeleteAccountUseCase) Execute(ctx context.Context, userID string) error {
    // 1. Soft delete do usuário
    user, _ := uc.userRepo.FindByID(ctx, userID)
    user.Ativo = false
    user.DeletedAt = time.Now()

    // 2. Anonimizar dados pessoais
    user.Nome = "[USUÁRIO REMOVIDO]"
    user.Email = fmt.Sprintf("deleted-%s@anonimizado.local", user.ID[:8])
    user.PasswordHash = ""

    uc.userRepo.Update(ctx, user)

    // 3. Revogar tokens JWT (blacklist ou invalidar refresh tokens)
    uc.jwtService.RevokeAllTokens(userID)

    // 4. Anonimizar audit_logs (opcional, se não quebrar integridade)
    // Substituir user_id por "DELETED" em logs antigos

    // 5. Registrar ação de exclusão
    uc.auditService.RecordDelete(ctx, user.TenantID, userID, "User", userID, audit.ActionDeleteAccount)

    return nil
}
```

**Handler:**
```go
// DELETE /api/v1/me
func (h *UserHandler) DeleteAccount(c echo.Context) error {
    userID := middleware.GetUserIDFromContext(c)

    // Confirmar senha antes de deletar (segurança)
    var req struct {
        Password string `json:"password" validate:"required"`
    }
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Senha necessária para confirmar exclusão")
    }

    // Validar senha
    user, _ := h.userRepo.FindByID(c.Request().Context(), userID)
    if !h.passwordHasher.Compare(user.PasswordHash, req.Password) {
        return echo.NewHTTPError(401, "Senha incorreta")
    }

    // Executar exclusão
    if err := h.deleteAccountUseCase.Execute(c.Request().Context(), userID); err != nil {
        return echo.NewHTTPError(500, "Erro ao deletar conta")
    }

    return c.JSON(200, map[string]string{
        "message": "Conta excluída com sucesso. Seus dados foram anonimizados."
    })
}
```

**Importante:**
- ❌ **NÃO deletar** dados necessários para obrigações legais (notas fiscais)
- ✅ **Soft delete** com flag `deleted_at`
- ✅ **Anonimizar** PII (nome, email, telefone)
- ✅ **Revogar** todos os tokens de acesso

---

### 4. **Correção (Art. 18, III)**

**Endpoint:** `PUT /api/v1/me`

**Implementação:**
```go
// Permitir usuário atualizar seus próprios dados
func (h *UserHandler) UpdateProfile(c echo.Context) error {
    userID := middleware.GetUserIDFromContext(c)

    var req dto.UpdateUserRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, err.Error())
    }

    // Atualizar apenas campos permitidos
    user, _ := h.userRepo.FindByID(c.Request().Context(), userID)
    user.Nome = req.Nome
    user.Email = req.Email // Validar se email não está em uso

    if err := h.userRepo.Update(c.Request().Context(), user); err != nil {
        return echo.NewHTTPError(500, "Erro ao atualizar perfil")
    }

    return c.JSON(200, user)
}
```

---

### 5. **Revogação de Consentimento (Art. 18, IX)**

**Endpoint:** `PUT /api/v1/me/preferences`

**Implementação:**
```go
func (h *UserHandler) UpdatePreferences(c echo.Context) error {
    userID := middleware.GetUserIDFromContext(c)

    var req struct {
        AnalyticsEnabled     bool `json:"analytics_enabled"`
        ErrorTrackingEnabled bool `json:"error_tracking_enabled"`
    }
    c.Bind(&req)

    // Salvar preferências
    prefs := &domain.UserPreferences{
        UserID:               userID,
        AnalyticsEnabled:     req.AnalyticsEnabled,
        ErrorTrackingEnabled: req.ErrorTrackingEnabled,
        UpdatedAt:            time.Now(),
    }

    h.preferencesRepo.Save(c.Request().Context(), prefs)

    return c.JSON(200, map[string]string{"message": "Preferências atualizadas"})
}
```

---

## ⏳ Política de Retenção

### Tabela de Retenção

| Tipo de Dado | Período de Retenção | Justificativa | Ação Pós-Retenção |
|--------------|---------------------|---------------|-------------------|
| **Users (ativos)** | Enquanto conta ativa | Contrato | N/A |
| **Users (deletados)** | 90 dias (anonimizado) | Suporte e fraude | Hard delete |
| **Audit Logs** | 90 dias | Segurança | Auto-delete |
| **Error Logs** | 30 dias | Debugging | Auto-delete |
| **Assinaturas (ativas)** | Enquanto ativa | Contrato | N/A |
| **Assinaturas (canceladas)** | 5 anos | Tributário/fiscal | Anonimizar PII |
| **Receitas/Despesas** | 5 anos | Tributário/fiscal | Manter |
| **Analytics** | Enquanto consentir | Melhoria produto | Delete se revogar |

### Jobs de Limpeza

**Cron: Limpar dados expirados**
```go
// internal/infrastructure/scheduler/cleanup_job.go
func (j *CleanupJob) Run() {
    ctx := context.Background()

    // 1. Hard delete de usuários soft-deleted há >90 dias
    j.userRepo.HardDeleteOlderThan(ctx, 90*24*time.Hour)

    // 2. Delete audit_logs >90 dias
    j.auditRepo.DeleteOlderThan(ctx, 90*24*time.Hour)

    // 3. Anonimizar assinaturas canceladas há >5 anos
    j.assinaturaRepo.AnonymizeOlderThan(ctx, 5*365*24*time.Hour)
}
```

**Schedule:** Diário às 03:00 UTC

---

## 🔒 Segurança e Controles

### Medidas Técnicas Implementadas

| Controle | Implementação | Status |
|----------|---------------|--------|
| **Criptografia em trânsito** | TLS 1.3 via NGINX | ✅ |
| **Criptografia em repouso** | Neon PostgreSQL (AES-256) | ✅ |
| **Hashing de senhas** | Bcrypt (cost 12) | ✅ |
| **Autenticação** | JWT RS256 assimétrico | ✅ |
| **Autorização** | RBAC com 4 roles | ✅ |
| **Rate limiting** | NGINX + backend (50 req/min) | ✅ |
| **Logs de auditoria** | Tabela audit_logs (90 dias) | ✅ |
| **Backup** | Neon PITR + pg_dump diário | ⏳ |
| **Monitoramento** | Prometheus + Grafana + Alertas | ✅ |
| **Testes de segurança** | SQL injection, XSS, CSRF | ✅ |

### Medidas Organizacionais

| Controle | Status |
|----------|--------|
| **Treinamento LGPD** | ⏳ Agendar |
| **DPO designado** | ⏳ Definir |
| **Termo de Confidencialidade** | ⏳ Criar |
| **Privacy by Design** | ✅ Aplicado |
| **Privacy by Default** | ✅ Aplicado |

---

## 📞 Contatos e DPO

### Encarregado de Dados (DPO)

**Nome:** [A definir]
**Email:** dpo@barberpro.dev
**Telefone:** [A definir]

**Responsabilidades:**
- Orientar sobre conformidade com LGPD
- Receber comunicações da ANPD
- Atender solicitações de titulares
- Realizar avaliações de impacto (DPIA)

### Canal de Atendimento ao Titular

**Email:** privacidade@barberpro.dev
**Prazo de resposta:** 15 dias úteis (conforme LGPD)

**Solicitações aceitas:**
- ✅ Acesso aos dados (GET /me)
- ✅ Correção de dados (PUT /me)
- ✅ Portabilidade (GET /me/export)
- ✅ Exclusão (DELETE /me)
- ✅ Revogação de consentimento
- ✅ Informações sobre tratamento

---

## 📄 Privacy Policy (Política de Privacidade)

**Local:** `https://barberpro.dev/privacy`

**Conteúdo (resumo):**

```markdown
# Política de Privacidade - Barber Analytics Pro

Última atualização: 15/11/2025

## 1. Quem somos
Barber Analytics Pro é um sistema SaaS de gestão para barbearias...

## 2. Quais dados coletamos
- Nome, email, senha (hash)
- CNPJ, telefone, endereço da barbearia
- Logs de acesso (IP, user agent)
- Dados de uso (analytics, com consentimento)

## 3. Por que coletamos
- Execução do contrato (prover o serviço)
- Legítimo interesse (segurança, prevenção de fraudes)
- Consentimento (analytics, error tracking)

## 4. Com quem compartilhamos
- Asaas (processamento de pagamentos)
- Neon (hospedagem de banco de dados)
- Sentry (error tracking, se consentir)
- Google Analytics (se consentir)

## 5. Seus direitos
- Acessar seus dados
- Corrigir dados incorretos
- Solicitar exclusão (direito ao esquecimento)
- Portabilidade de dados
- Revogar consentimento

## 6. Como exercer direitos
Email: privacidade@barberpro.dev

## 7. Retenção de dados
- Dados ativos: Enquanto conta ativa
- Dados deletados: 90 dias (anonimizado)
- Dados fiscais: 5 anos (obrigação legal)

## 8. Segurança
- TLS 1.3
- Senhas com bcrypt
- JWT RS256
- Backups criptografados

## 9. Contato
DPO: dpo@barberpro.dev
```

---

## ✅ Checklist de Conformidade

### Documentação
- [x] Política de Privacidade criada
- [x] Inventário de dados mapeado
- [ ] Termo de Consentimento redigido
- [ ] DPIA (Avaliação de Impacto) realizada
- [ ] Registro de operações de tratamento

### Técnico
- [ ] Banner de consentimento implementado (frontend)
- [ ] Endpoint GET /me (acesso)
- [ ] Endpoint GET /me/export (portabilidade)
- [ ] Endpoint DELETE /me (exclusão)
- [ ] Endpoint PUT /me/preferences (revogação)
- [ ] Job de limpeza automática (retenção)
- [ ] Logs de auditoria de acessos a dados

### Organizacional
- [ ] DPO designado
- [ ] Treinamento da equipe
- [ ] Procedimento de resposta a incidentes
- [ ] Canal de atendimento ao titular

---

## 📚 Referências

- [Lei nº 13.709/2018 (LGPD)](http://www.planalto.gov.br/ccivil_03/_ato2015-2018/2018/lei/l13709.htm)
- [Guia ANPD — Agentes de Tratamento](https://www.gov.br/anpd/)
- [GDPR (referência internacional)](https://gdpr.eu/)
- [ISO 27701 (Privacy Information Management)](https://www.iso.org/standard/71670.html)

---

**Última Atualização:** 15/11/2025
**Versão:** 1.0.0
**Responsável:** Equipe Barber Analytics Pro
**Revisão:** A cada 6 meses ou quando houver mudanças significativas
