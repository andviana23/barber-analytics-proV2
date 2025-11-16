# 📡 API Reference

**Versão:** 2.0  
**Data:** 14/11/2025  
**Base URL:** `https://api.seudominio.com/v2`

---

## 📋 Índice

1. [Autenticação](#autenticação)
2. [Receitas](#receitas)
3. [Despesas](#despesas)
4. [Assinaturas](#assinaturas)
5. [Fluxo de Caixa](#fluxo-de-caixa)
6. [Erros](#erros)

---

## 🔐 Autenticação

### Login

```http
POST /auth/login
Content-Type: application/json

{
  \"email\": \"usuario@example.com\",
  \"password\": \"senha123\"
}

HTTP/1.1 200 OK
{
  \"access_token\": \"eyJ0eXAiOiJKV1QiLCJhbGc...\",
  \"refresh_token\": \"refresh_eyJ0eXAiOiJKV1QiLCJhbGc...\",
  \"expires_in\": 900,
  \"user\": {
    \"id\": \"user-123\",
    \"email\": \"usuario@example.com\",
    \"nome\": \"João Silva\",
    \"tenant_id\": \"tenant-abc\",
    \"role\": \"owner\"
  }
}
```

### Refresh Token

```http
POST /auth/refresh
Content-Type: application/json

{
  \"refresh_token\": \"refresh_eyJ0eXAiOiJKV1QiLCJhbGc...\"
}

HTTP/1.1 200 OK
{
  \"access_token\": \"eyJ0eXAiOiJKV1QiLCJhbGc...\",
  \"expires_in\": 900
}
```

### Headers Obrigatórios

```http
Authorization: Bearer {access_token}
Content-Type: application/json
```

---

## 💰 Receitas

### Criar Receita

```http
POST /financial/receitas
Authorization: Bearer {token}
Content-Type: application/json

{
  \"descricao\": \"Corte de cabelo\",
  \"valor\": \"50.00\",
  \"categoria_id\": \"cat-001\",
  \"metodo_pagamento\": \"PIX\",
  \"data\": \"2024-11-14\"
}

HTTP/1.1 201 Created
{
  \"id\": \"rcta-001\",
  \"descricao\": \"Corte de cabelo\",
  \"valor\": \"50.00\",
  \"status\": \"CONFIRMADO\",
  \"criado_em\": \"2024-11-14T10:30:00Z\"
}
```

### Listar Receitas

```http
GET /financial/receitas?from=2024-11-01&to=2024-11-30&categoria_id=cat-001&page=1&page_size=50
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  \"data\": [
    {
      \"id\": \"rcta-001\",
      \"descricao\": \"Corte\",
      \"valor\": \"50.00\",
      \"categoria_id\": \"cat-001\",
      \"data\": \"2024-11-14\",
      \"status\": \"CONFIRMADO\",
      \"criado_em\": \"2024-11-14T10:30:00Z\"
    }
  ],
  \"pagination\": {
    \"total\": 125,
    \"page\": 1,
    \"page_size\": 50,
    \"total_pages\": 3
  }
}
```

### Atualizar Receita

```http
PUT /financial/receitas/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  \"descricao\": \"Corte de cabelo (atualizado)\",
  \"valor\": \"55.00\"
}

HTTP/1.1 200 OK
```

### Deletar Receita

```http
DELETE /financial/receitas/{id}
Authorization: Bearer {token}

HTTP/1.1 204 No Content
```

---

## 💸 Despesas

### Criar Despesa

```http
POST /financial/despesas
Authorization: Bearer {token}
Content-Type: application/json

{
  \"descricao\": \"Aluguel\",
  \"valor\": \"1000.00\",
  \"categoria_id\": \"cat-expenses\",
  \"metodo_pagamento\": \"TRANSFERENCIA\",
  \"data\": \"2024-11-14\"
}

HTTP/1.1 201 Created
{
  \"id\": \"desp-001\",
  \"descricao\": \"Aluguel\",
  \"valor\": \"1000.00\",
  \"status\": \"PENDENTE\",
  \"criado_em\": \"2024-11-14T10:30:00Z\"
}
```

### Listar Despesas

```http
GET /financial/despesas?from=2024-11-01&to=2024-11-30&status=PENDENTE
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  \"data\": [...],
  \"pagination\": {...}
}
```

---

## 🎟️ Assinaturas

### Criar Assinatura

```http
POST /subscriptions
Authorization: Bearer {token}
Content-Type: application/json

{
  \"plan_id\": \"plan-123\",
  \"barbeiro_id\": \"user-456\",
  \"data_inicio\": \"2024-11-14\"
}

HTTP/1.1 201 Created
{
  \"id\": \"sub-001\",
  \"plan_id\": \"plan-123\",
  \"status\": \"ATIVA\",
  \"data_inicio\": \"2024-11-14\",
  \"proxima_fatura_data\": \"2024-12-14\",
  \"criado_em\": \"2024-11-14T10:30:00Z\"
}
```

### Listar Assinaturas

```http
GET /subscriptions?status=ATIVA
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  \"data\": [
    {
      \"id\": \"sub-001\",
      \"plan_id\": \"plan-123\",
      \"barbeiro_id\": \"user-456\",
      \"status\": \"ATIVA\",
      \"data_inicio\": \"2024-11-14\",
      \"proxima_fatura_data\": \"2024-12-14\"
    }
  ],
  \"pagination\": {...}
}
```

### Cancelar Assinatura

```http
POST /subscriptions/{id}/cancel
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  \"id\": \"sub-001\",
  \"status\": \"CANCELADA\",
  \"data_fim\": \"2024-11-14\"
}
```

---

## 💹 Fluxo de Caixa

### Obter Fluxo de Caixa

```http
GET /financial/cashflow?from=2024-11-01&to=2024-11-30
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  \"periodo\": {
    \"data_inicio\": \"2024-11-01\",
    \"data_fim\": \"2024-11-30\"
  },
  \"saldo_inicial\": \"5000.00\",
  \"entradas\": \"3250.00\",
  \"saidas\": \"1500.00\",
  \"saldo_final\": \"6750.00\"
}
```

### Projeção Fluxo de Caixa

```http
GET /financial/cashflow/projection
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  \"dias\": [
    {
      \"data\": \"2024-11-15\",
      \"saldo_projetado\": \"6800.00\",
      \"entradas_projetadas\": \"100.00\",
      \"saidas_projetadas\": \"50.00\"
    }
  ]
}
```

---

## ⚠️ Erros

### Estrutura de Erro

```json
{
  \"error\": {
    \"code\": \"INVALID_REQUEST\",
    \"message\": \"Descrição do erro\",
    \"details\": {
      \"field\": \"valor\",
      \"reason\": \"value must be greater than 0\"
    }
  },
  \"trace_id\": \"trace-123456\"
}
```

### Códigos de Erro

| Status | Código | Descrição |
|--------|--------|-----------|
| 400 | INVALID_REQUEST | Requisição inválida |
| 401 | UNAUTHORIZED | Não autenticado |
| 403 | FORBIDDEN | Sem permissão |
| 404 | NOT_FOUND | Recurso não encontrado |
| 422 | UNPROCESSABLE_ENTITY | Dados inválidos |
| 429 | RATE_LIMITED | Limite de requisições |
| 500 | INTERNAL_ERROR | Erro interno |

---

**Última atualização:** 14/11/2025
