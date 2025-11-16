# 🔗 Integrações Asaas

**Versão:** 1.0  
**Data:** 14/11/2025  
**Status:** Design Finalizado

---

## 📋 Índice

1. [Setup Asaas](#setup-asaas)
2. [Autenticação](#autenticação)
3. [APIs Utilizadas](#apis-utilizadas)
4. [Fluxos de Integração](#fluxos-de-integração)
5. [Error Handling](#error-handling)
6. [Webhooks (Futuro)](#webhooks-futuro)

---

## 🎯 Setup Asaas

### Pré-requisitos

1. **Conta Asaas**
   - [ ] Ir para https://asaas.com
   - [ ] Criar conta (modo produção)
   - [ ] Completar verificação de identidade

2. **API Key**
   - [ ] Acessar: Settings → API → Gerar API Key
   - [ ] Guardar chave segura: `ASAAS_API_KEY`
   - [ ] Nunca commitar no git

3. **Variáveis de Ambiente**
   ```bash
   ASAAS_API_KEY=sk_live_xxxxxxxxxxxxx
   ASAAS_BASE_URL=https://api.asaas.com/v3
   ASAAS_REQUEST_TIMEOUT=30s
   ```

---

## 🔐 Autenticação

Asaas usa API Key em header:

```http
GET https://api.asaas.com/v3/subscriptions
Authorization: Bearer sk_live_xxxxxxxxxxxxx
```

---

## 📡 APIs Utilizadas

### 1. Subscriptions (Assinaturas Recorrentes)

**Endpoint:** `POST /subscriptions`

**Objetivo:** Criar assinatura recorrente para um cliente.

```http
POST https://api.asaas.com/v3/subscriptions
Authorization: Bearer sk_live_xxxxxxxxxxxxx
Content-Type: application/json

{
  \"customer\": \"cus_000000000000000001\",
  \"billingType\": \"CREDIT_CARD\",
  \"value\": 49.90,
  \"nextDueDate\": \"2024-12-14\",
  \"cycle\": \"MONTHLY\",
  \"description\": \"Plano Premium - Barbeiro João\"
}

HTTP/1.1 200 OK
{
  \"object\": \"subscription\",
  \"id\": \"sub_000000000000000001\",
  \"dateCreated\": \"2024-11-14\",
  \"status\": \"ACTIVE\",
  \"customer\": \"cus_000000000000000001\",
  \"value\": 49.90,
  \"nextDueDate\": \"2024-12-14\",
  \"cycle\": \"MONTHLY\"
}
```

### 2. Invoices (Faturas)

**Endpoint:** `GET /invoices`

**Objetivo:** Listar faturas de uma assinatura.

```http
GET https://api.asaas.com/v3/invoices?subscription=sub_000000000000000001&limit=100
Authorization: Bearer sk_live_xxxxxxxxxxxxx

HTTP/1.1 200 OK
{
  \"object\": \"list\",
  \"hasMore\": false,
  \"data\": [
    {
      \"object\": \"invoice\",
      \"id\": \"inv_000000000000000001\",
      \"subscription\": \"sub_000000000000000001\",
      \"value\": 49.90,
      \"status\": \"PENDING\",
      \"dueDate\": \"2024-12-14\",
      \"invoiceUrl\": \"https://app.asaas.com/i/...\",
      \"dateCreated\": \"2024-11-14\"
    }
  ]
}
```

### 3. Get Invoice Details

**Endpoint:** `GET /invoices/{id}`

**Objetivo:** Obter detalhes de uma fatura específica.

```http
GET https://api.asaas.com/v3/invoices/inv_000000000000000001
Authorization: Bearer sk_live_xxxxxxxxxxxxx

HTTP/1.1 200 OK
{
  \"object\": \"invoice\",
  \"id\": \"inv_000000000000000001\",
  \"status\": \"CONFIRMED\",
  \"value\": 49.90,
  \"netValue\": 49.90,
  \"grossValue\": 49.90,
  \"dueDate\": \"2024-12-14\",
  \"originalDueDate\": \"2024-12-14\",
  \"paymentDate\": \"2024-12-10\",
  \"clientPaymentDate\": \"2024-12-10\",
  \"invoiceUrl\": \"https://app.asaas.com/i/...\",
  \"estimatedPaymentDate\": \"2024-12-10\"
}
```

### 4. Cancel Subscription

**Endpoint:** `DELETE /subscriptions/{id}`

**Objetivo:** Cancelar assinatura.

```http
DELETE https://api.asaas.com/v3/subscriptions/sub_000000000000000001
Authorization: Bearer sk_live_xxxxxxxxxxxxx

HTTP/1.1 200 OK
{
  \"object\": \"subscription\",
  \"id\": \"sub_000000000000000001\",
  \"status\": \"CANCELLED\"
}
```

---

## 🔄 Fluxos de Integração

### Fluxo 1: Criar Assinatura

```
User cria assinatura no front
        ↓
POST /subscriptions
        ↓
Backend: Validar dados
        ↓
Chamar Asaas: CreateSubscription()
        ↓
Asaas retorna ID
        ↓
Persistir: Assinatura(asaas_subscription_id)
        ↓
Response sucesso
```

**Código Go:**

```go
type AsaasClient struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

type CreateSubscriptionReq struct {
    Customer    string `json:\"customer\"`
    BillingType string `json:\"billingType\"`
    Value       float64 `json:\"value\"`
    NextDueDate time.Time `json:\"nextDueDate\"`
    Cycle       string `json:\"cycle\"` // MONTHLY, BIMONTHLY, QUARTERLY, SEMIANNUAL, YEARLY
}

type CreateSubscriptionResp struct {
    ID         string `json:\"id\"`
    Status     string `json:\"status\"`
    Value      float64 `json:\"value\"`
    NextDueDate time.Time `json:\"nextDueDate\"`
}

func (c *AsaasClient) CreateSubscription(
    req *CreateSubscriptionReq) (*CreateSubscriptionResp, error) {
    
    // Serializar request
    payload, _ := json.Marshal(req)
    
    // Fazer request
    httpReq, _ := http.NewRequest(
        \"POST\",
        c.baseURL + \"/subscriptions\",
        bytes.NewBuffer(payload),
    )
    httpReq.Header.Set(\"Authorization\", \"Bearer \" + c.apiKey)
    httpReq.Header.Set(\"Content-Type\", \"application/json\")
    
    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    // Parse response
    body, _ := ioutil.ReadAll(resp.Body)
    
    if resp.StatusCode != http.StatusOK {
        return nil, handleAsaasError(resp.StatusCode, body)
    }
    
    var result CreateSubscriptionResp
    json.Unmarshal(body, &result)
    
    return &result, nil
}
```

### Fluxo 2: Sincronizar Faturas (Cron Diário)

```
Cron executa 02:00
        ↓
Para cada tenant:
  Para cada assinatura ativa:
    Chamar Asaas: ListInvoices(subscription_id)
            ↓
    Para cada fatura:
      ├─ Se não existe no DB: Inserir
      └─ Se existe: Atualizar status
```

---

## ⚠️ Error Handling

### Mapeamento de Erros Asaas

```go
func handleAsaasError(statusCode int, body []byte) error {
    switch statusCode {
    case 400:
        return errors.New(\"Bad Request: dados inválidos\")
    
    case 401:
        return errors.New(\"Unauthorized: API key inválida\")
    
    case 403:
        return errors.New(\"Forbidden: sem permissão\")
    
    case 404:
        return errors.New(\"Not Found: recurso não existe\")
    
    case 422:
        // Erro de validação
        var errResp map[string]interface{}
        json.Unmarshal(body, &errResp)
        return fmt.Errorf(\"Validation error: %v\", errResp)
    
    case 429:
        return errors.New(\"Rate Limited: aguarde antes de retry\")
    
    case 500, 502, 503, 504:
        return errors.New(\"Server Error: tente novamente\")
    
    default:
        return fmt.Errorf(\"Unknown error: %d\", statusCode)
    }
}
```

### Retry Strategy

```go
func (c *AsaasClient) ListInvoicesWithRetry(
    subscriptionID string, maxRetries int) ([]*Invoice, error) {
    
    var lastErr error
    backoff := time.Second
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        invoices, err := c.ListInvoices(subscriptionID)
        
        if err == nil {
            return invoices, nil
        }
        
        lastErr = err
        
        // Retry apenas para erros transitórios
        if isTransient(err) {
            time.Sleep(backoff)
            backoff *= 2
        } else {
            return nil, err
        }
    }
    
    return nil, lastErr
}

func isTransient(err error) bool {
    // 429 (rate limit) e 5xx são transitórios
    return strings.Contains(err.Error(), \"Rate Limited\") ||
           strings.Contains(err.Error(), \"Server Error\")
}
```

---

## 🪝 Webhooks (Futuro)

### Configurar Webhook

1. **URL de Webhook**
   ```
   https://api.seudominio.com/v2/webhooks/asaas
   ```

2. **Eventos a monitorar**
   - `PAYMENT_RECEIVED` - Pagamento confirmado
   - `PAYMENT_PENDING` - Pagamento pendente
   - `PAYMENT_CONFIRMED` - Pagamento confirmado

3. **Payload esperado**
   ```json
   {
     \"event\": \"PAYMENT_RECEIVED\",
     \"data\": {
       \"id\": \"inv_000000000000000001\",
       \"subscription\": \"sub_000000000000000001\",
       \"value\": 49.90,
       \"status\": \"RECEIVED\",
       \"paymentDate\": \"2024-12-10\"
     }
   }
   ```

4. **Handler no backend**
   ```go
   func HandleAsaasWebhook(c echo.Context) error {
       var payload WebhookPayload
       c.Bind(&payload)
       
       switch payload.Event {
       case \"PAYMENT_RECEIVED\":
           // Atualizar fatura no DB
           // Criar receita no financeiro
           break
       }
       
       return c.JSON(200, map[string]string{\"status\": \"ok\"})
   }
   ```

---

## 📚 Referências

- [Documentação Asaas](https://docs.asaas.com)
- [API Asaas Subscriptions](https://docs.asaas.com/reference/crear-una-suscripcin)
- [Status de Faturas](https://docs.asaas.com/reference/obtener-facturas)

---

**Status:** ✅ Documentado
