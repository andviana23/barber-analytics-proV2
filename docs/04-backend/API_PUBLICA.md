> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 📡 API Reference

**Versão:** 2.0
**Data:** 14/11/2025
**Base URL:** `https://api.seudominio.com/v2`

---

## 📋 Índice

1. [Autenticação](#autenticação)
2. [Cadastro](#cadastro)
   - [Clientes](#clientes)
   - [Profissionais](#profissionais)
   - [Serviços](#serviços)
   - [Meios de Pagamento](#meios-de-pagamento)
3. [Receitas](#receitas)
4. [Despesas](#despesas)
5. [Assinaturas](#assinaturas)
6. [Fluxo de Caixa](#fluxo-de-caixa)
7. [Barber Turn (Lista da Vez)](#barber-turn-lista-da-vez)
8. [Erros](#erros)

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

## 📇 Cadastro

### Clientes

#### Criar Cliente

```http
POST /api/v1/cadastro/clientes
Authorization: Bearer {token}
Content-Type: application/json

{
  "nome": "João Silva",
  "email": "joao@example.com",
  "telefone": "+5511999887766",
  "cpf": "12345678901",
  "data_nascimento": "1990-05-15",
  "genero": "M",
  "endereco": {
    "logradouro": "Rua das Flores",
    "numero": "123",
    "complemento": "Apto 45",
    "bairro": "Centro",
    "cidade": "São Paulo",
    "uf": "SP",
    "cep": "01310100",
    "pais": "Brasil"
  },
  "tags": ["VIP", "Premium"],
  "observacoes": "Cliente preferencial"
}

HTTP/1.1 201 Created
{
  "id": "c1111111-1111-1111-1111-111111111111",
  "tenant_id": "tenant-abc",
  "nome": "João Silva",
  "email": "joao@example.com",
  "telefone": "+55 11 99988-7766",
  "cpf": "123.456.789-01",
  "data_nascimento": "1990-05-15",
  "genero": "M",
  "endereco": {
    "logradouro": "Rua das Flores",
    "numero": "123",
    "complemento": "Apto 45",
    "bairro": "Centro",
    "cidade": "São Paulo",
    "uf": "SP",
    "cep": "01310-100",
    "pais": "Brasil"
  },
  "tags": ["VIP", "Premium"],
  "observacoes": "Cliente preferencial",
  "ativo": true,
  "criado_em": "2024-11-14T10:30:00Z",
  "atualizado_em": "2024-11-14T10:30:00Z"
}
```

**Campos Obrigatórios:**

- `nome` (string, min 2 caracteres)
- `telefone` (string, formato +DDI DDD número)

**Campos Opcionais:**

- `email` (string, formato email válido, único por tenant)
- `cpf` (string, 11 dígitos, validação de dígitos verificadores)
- `data_nascimento` (string, formato YYYY-MM-DD)
- `genero` (string, "M" ou "F")
- `endereco` (objeto com logradouro, numero, bairro, cidade, uf, cep, pais)
- `tags` (array de strings)
- `observacoes` (string)

**Erros Comuns:**

- `400 BAD_REQUEST` - Campos obrigatórios faltando ou inválidos
- `409 CONFLICT` - Email ou CPF já cadastrado no tenant
- `403 FORBIDDEN` - Usuário sem permissão (requer owner ou manager)

---

#### Atualizar Cliente

```http
PUT /api/v1/cadastro/clientes/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "nome": "João Silva Atualizado",
  "telefone": "+5511999887766",
  "email": "joao.novo@example.com"
}

HTTP/1.1 200 OK
{
  "id": "c1111111-1111-1111-1111-111111111111",
  "nome": "João Silva Atualizado",
  "telefone": "+55 11 99988-7766",
  "email": "joao.novo@example.com",
  "atualizado_em": "2024-11-14T11:45:00Z"
}
```

**Validações:**

- Email não pode ser alterado para um já existente em outro cliente do mesmo tenant
- Campos não enviados mantêm valores anteriores
- Tenant isolation: não é possível atualizar cliente de outro tenant

**Erros Comuns:**

- `404 NOT_FOUND` - Cliente não encontrado ou pertence a outro tenant
- `409 CONFLICT` - Email já em uso por outro cliente

---

#### Buscar Cliente

```http
GET /api/v1/cadastro/clientes/{id}
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  "id": "c1111111-1111-1111-1111-111111111111",
  "tenant_id": "tenant-abc",
  "nome": "João Silva",
  "email": "joao@example.com",
  "telefone": "+55 11 99988-7766",
  "cpf": "123.456.789-01",
  "ativo": true,
  "criado_em": "2024-11-14T10:30:00Z"
}
```

**Erros Comuns:**

- `404 NOT_FOUND` - Cliente não encontrado ou pertence a outro tenant

---

#### Listar Clientes

```http
GET /api/v1/cadastro/clientes?nome=João&ativo=true&page=1&page_size=20
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  "items": [
    {
      "id": "c1111111-1111-1111-1111-111111111111",
      "nome": "João Silva",
      "email": "joao@example.com",
      "telefone": "+55 11 99988-7766",
      "ativo": true
    }
  ],
  "total": 15,
  "page": 1,
  "page_size": 20,
  "total_pages": 1
}
```

**Query Parameters:**

- `nome` (string) - Busca parcial por nome (ILIKE)
- `email` (string) - Busca exata por email
- `ativo` (boolean) - Filtrar por status ativo/inativo
- `tags` (string) - Filtrar por tag específica
- `page` (int, default: 1) - Número da página
- `page_size` (int, default: 50, max: 100) - Tamanho da página

**Isolamento Multi-Tenant:**

- Retorna apenas clientes do tenant autenticado
- Total count considera apenas registros do tenant

---

#### Deletar Cliente (Soft Delete)

```http
DELETE /api/v1/cadastro/clientes/{id}
Authorization: Bearer {token}

HTTP/1.1 204 No Content
```

**Comportamento:**

- Soft delete: marca `ativo = false` ao invés de remover registro
- Cliente inativo não aparece em listagens com filtro `ativo=true`
- Histórico de atendimentos é preservado

**Erros Comuns:**

- `404 NOT_FOUND` - Cliente não encontrado ou pertence a outro tenant

---

### Profissionais

#### Criar Profissional

```http
POST /api/v1/cadastro/profissionais
Authorization: Bearer {token}
Content-Type: application/json

{
  "nome": "Carlos Barbeiro",
  "email": "carlos@barbearia.com",
  "telefone": "+5511955443322",
  "cpf": "11122233344",
  "data_admissao": "2020-01-15",
  "comissao": 30.00,
  "tipo_comissao": "PERCENTUAL",
  "especialidades": ["Corte", "Barba", "Coloração"],
  "horario_trabalho": {
    "segunda": {"ativo": true, "entrada": "08:00", "saida": "18:00"},
    "terca": {"ativo": true, "entrada": "08:00", "saida": "18:00"},
    "quarta": {"ativo": true, "entrada": "08:00", "saida": "18:00"},
    "quinta": {"ativo": true, "entrada": "08:00", "saida": "18:00"},
    "sexta": {"ativo": true, "entrada": "08:00", "saida": "18:00"},
    "sabado": {"ativo": true, "entrada": "08:00", "saida": "14:00"},
    "domingo": {"ativo": false}
  }
}

HTTP/1.1 201 Created
{
  "id": "p1111111-1111-1111-1111-111111111111",
  "tenant_id": "tenant-abc",
  "nome": "Carlos Barbeiro",
  "email": "carlos@barbearia.com",
  "telefone": "+55 11 95544-3322",
  "cpf": "111.222.333-44",
  "data_admissao": "2020-01-15",
  "comissao": 30.00,
  "tipo_comissao": "PERCENTUAL",
  "especialidades": ["Corte", "Barba", "Coloração"],
  "status": "ATIVO",
  "criado_em": "2024-11-14T10:30:00Z"
}
```

**Campos Obrigatórios:**

- `nome` (string)
- `telefone` (string, formato +DDI DDD número)
- `comissao` (decimal)
- `tipo_comissao` (enum: "PERCENTUAL" ou "FIXO")

**Validações:**

- CPF único por tenant
- `comissao` entre 0-100 se `tipo_comissao = PERCENTUAL`
- `comissao` ≥ 0 se `tipo_comissao = FIXO`
- `status` iniciado como "ATIVO"

**Autorização:**

- ✅ Permitido: `owner`
- ❌ Bloqueado: `manager`, `accountant`, `employee`

**Erros Comuns:**

- `400 BAD_REQUEST` - Comissão inválida para o tipo
- `403 FORBIDDEN` - Apenas owner pode criar profissionais
- `409 CONFLICT` - CPF já cadastrado no tenant

---

#### Atualizar Profissional

```http
PUT /api/v1/cadastro/profissionais/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "nome": "Carlos Barbeiro Sênior",
  "comissao": 35.00,
  "status": "ATIVO"
}

HTTP/1.1 200 OK
{
  "id": "p1111111-1111-1111-1111-111111111111",
  "nome": "Carlos Barbeiro Sênior",
  "comissao": 35.00,
  "status": "ATIVO",
  "atualizado_em": "2024-11-14T11:45:00Z"
}
```

**Status Permitidos:**

- `ATIVO` - Profissional ativo e pode atender
- `INATIVO` - Temporariamente inativo
- `AFASTADO` - Afastado por motivo específico
- `DEMITIDO` - Desligado da empresa (soft delete)

**Autorização:**

- ✅ Permitido: `owner`
- ❌ Bloqueado: `manager`, `accountant`, `employee`

---

#### Buscar Profissional

```http
GET /api/v1/cadastro/profissionais/{id}
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  "id": "p1111111-1111-1111-1111-111111111111",
  "nome": "Carlos Barbeiro",
  "email": "carlos@barbearia.com",
  "comissao": 30.00,
  "tipo_comissao": "PERCENTUAL",
  "especialidades": ["Corte", "Barba"],
  "status": "ATIVO"
}
```

---

#### Listar Profissionais

```http
GET /api/v1/cadastro/profissionais?status=ATIVO&especialidade=Corte&page=1&page_size=20
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  "items": [
    {
      "id": "p1111111-1111-1111-1111-111111111111",
      "nome": "Carlos Barbeiro",
      "especialidades": ["Corte", "Barba"],
      "status": "ATIVO",
      "comissao": 30.00
    }
  ],
  "total": 5,
  "page": 1,
  "page_size": 20,
  "total_pages": 1
}
```

**Query Parameters:**

- `nome` (string) - Busca parcial por nome
- `status` (enum) - Filtrar por status (ATIVO, INATIVO, AFASTADO, DEMITIDO)
- `especialidade` (string) - Filtrar por especialidade (array contains)
- `page` (int, default: 1)
- `page_size` (int, default: 50, max: 100)

---

#### Deletar Profissional (Soft Delete)

```http
DELETE /api/v1/cadastro/profissionais/{id}
Authorization: Bearer {token}

HTTP/1.1 204 No Content
```

**Comportamento:**

- Soft delete: altera `status = DEMITIDO` e registra `data_demissao`
- Profissional demitido não pode ser atribuído a novos serviços
- Histórico de atendimentos/comissões preservado

**Autorização:**

- ✅ Permitido: `owner`
- ❌ Bloqueado: `manager`, `accountant`, `employee`

---

### Serviços

#### Criar Serviço

```http
POST /api/v1/cadastro/servicos
Authorization: Bearer {token}
Content-Type: application/json

{
  "nome": "Corte Masculino",
  "descricao": "Corte de cabelo masculino tradicional",
  "preco": 50.00,
  "duracao": 30,
  "comissao": 30.00,
  "categoria_id": "cat-001",
  "cor": "#FF5733",
  "profissionais_ids": ["p1111111-1111-1111-1111-111111111111"],
  "tags": ["Cabelo", "Masculino"]
}

HTTP/1.1 201 Created
{
  "id": "s1111111-1111-1111-1111-111111111111",
  "tenant_id": "tenant-abc",
  "nome": "Corte Masculino",
  "descricao": "Corte de cabelo masculino tradicional",
  "preco": 50.00,
  "duracao": 30,
  "comissao": 30.00,
  "categoria_id": "cat-001",
  "cor": "#FF5733",
  "profissionais_ids": ["p1111111-1111-1111-1111-111111111111"],
  "tags": ["Cabelo", "Masculino"],
  "ativo": true,
  "criado_em": "2024-11-14T10:30:00Z"
}
```

**Campos Obrigatórios:**

- `nome` (string)
- `preco` (decimal > 0)
- `duracao` (int, minutos, ≥ 5)
- `comissao` (decimal, 0-100%)

**Validações:**

- `duracao` mínima de 5 minutos
- `comissao` entre 0 e 100%
- `profissionais_ids` devem existir e pertencer ao tenant

**Erros Comuns:**

- `400 BAD_REQUEST` - Duração inválida ou comissão fora do range
- `404 NOT_FOUND` - Profissional não encontrado

---

#### Atualizar Serviço

```http
PUT /api/v1/cadastro/servicos/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "nome": "Corte Masculino Premium",
  "preco": 60.00,
  "duracao": 45,
  "profissionais_ids": ["p1111111-1111-1111-1111-111111111111", "p2222222-2222-2222-2222-222222222222"]
}

HTTP/1.1 200 OK
{
  "id": "s1111111-1111-1111-1111-111111111111",
  "nome": "Corte Masculino Premium",
  "preco": 60.00,
  "duracao": 45,
  "profissionais_ids": ["p1111111-1111-1111-1111-111111111111", "p2222222-2222-2222-2222-222222222222"],
  "atualizado_em": "2024-11-14T11:45:00Z"
}
```

---

#### Buscar Serviço

```http
GET /api/v1/cadastro/servicos/{id}
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  "id": "s1111111-1111-1111-1111-111111111111",
  "nome": "Corte Masculino",
  "preco": 50.00,
  "duracao": 30,
  "comissao": 30.00,
  "profissionais_ids": ["p1111111-1111-1111-1111-111111111111"],
  "ativo": true
}
```

---

#### Listar Serviços

```http
GET /api/v1/cadastro/servicos?nome=Corte&ativo=true&categoria_id=cat-001&page=1&page_size=20
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  "items": [
    {
      "id": "s1111111-1111-1111-1111-111111111111",
      "nome": "Corte Masculino",
      "preco": 50.00,
      "duracao": 30,
      "ativo": true
    }
  ],
  "total": 10,
  "page": 1,
  "page_size": 20,
  "total_pages": 1
}
```

**Query Parameters:**

- `nome` (string) - Busca parcial por nome
- `ativo` (boolean) - Filtrar por status ativo/inativo
- `categoria_id` (uuid) - Filtrar por categoria
- `tags` (string) - Filtrar por tag (array overlap)
- `page`, `page_size`

---

#### Deletar Serviço (Soft Delete)

```http
DELETE /api/v1/cadastro/servicos/{id}
Authorization: Bearer {token}

HTTP/1.1 204 No Content
```

**Comportamento:**

- Soft delete: marca `ativo = false`
- Serviço inativo não pode ser usado em novos agendamentos
- Histórico de atendimentos preservado

---

### Meios de Pagamento

#### Criar Meio de Pagamento

```http
POST /api/v1/cadastro/meios-pagamento
Authorization: Bearer {token}
Content-Type: application/json

{
  "nome": "Cartão de Crédito",
  "tipo": "CREDITO",
  "taxa": 2.5,
  "taxa_fixa": 0.39,
  "icone": "credit_card",
  "cor": "#4CAF50"
}

HTTP/1.1 201 Created
{
  "id": "m1111111-1111-1111-1111-111111111111",
  "tenant_id": "tenant-abc",
  "nome": "Cartão de Crédito",
  "tipo": "CREDITO",
  "taxa": 2.5,
  "taxa_fixa": 0.39,
  "icone": "credit_card",
  "cor": "#4CAF50",
  "ativo": true,
  "criado_em": "2024-11-14T10:30:00Z"
}
```

**Campos Obrigatórios:**

- `nome` (string)
- `tipo` (enum: "DINHEIRO", "PIX", "CREDITO", "DEBITO", "TRANSFERENCIA")

**Campos Opcionais:**

- `taxa` (decimal, percentual 0-100%)
- `taxa_fixa` (decimal, ≥ 0)
- `icone` (string)
- `cor` (string, formato hex)

**Validações:**

- `taxa` entre 0 e 100%
- `taxa_fixa` ≥ 0

---

#### Atualizar Meio de Pagamento

```http
PUT /api/v1/cadastro/meios-pagamento/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "nome": "Crédito Visa/Master",
  "taxa": 3.0,
  "taxa_fixa": 0.49
}

HTTP/1.1 200 OK
{
  "id": "m1111111-1111-1111-1111-111111111111",
  "nome": "Crédito Visa/Master",
  "taxa": 3.0,
  "taxa_fixa": 0.49,
  "atualizado_em": "2024-11-14T11:45:00Z"
}
```

---

#### Buscar Meio de Pagamento

```http
GET /api/v1/cadastro/meios-pagamento/{id}
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  "id": "m1111111-1111-1111-1111-111111111111",
  "nome": "Cartão de Crédito",
  "tipo": "CREDITO",
  "taxa": 2.5,
  "taxa_fixa": 0.39,
  "ativo": true
}
```

---

#### Listar Meios de Pagamento

```http
GET /api/v1/cadastro/meios-pagamento?tipo=CREDITO&ativo=true&page=1&page_size=20
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  "items": [
    {
      "id": "m1111111-1111-1111-1111-111111111111",
      "nome": "Cartão de Crédito",
      "tipo": "CREDITO",
      "taxa": 2.5,
      "ativo": true
    }
  ],
  "total": 5,
  "page": 1,
  "page_size": 20,
  "total_pages": 1
}
```

**Query Parameters:**

- `nome` (string) - Busca parcial por nome
- `tipo` (enum) - Filtrar por tipo
- `ativo` (boolean) - Filtrar por status
- `page`, `page_size`

---

#### Deletar Meio de Pagamento (Soft Delete)

```http
DELETE /api/v1/cadastro/meios-pagamento/{id}
Authorization: Bearer {token}

HTTP/1.1 204 No Content
```

**Comportamento:**

- Soft delete: marca `ativo = false`
- Meio inativo não pode ser usado em novos pagamentos
- Histórico de transações preservado

---

### Parâmetros de Query Comuns (Cadastro)

Todos os endpoints de listagem do módulo de cadastro suportam:

| Parâmetro   | Tipo    | Descrição                        | Default | Max |
| ----------- | ------- | -------------------------------- | ------- | --- |
| `page`      | int     | Número da página                 | 1       | -   |
| `page_size` | int     | Tamanho da página                | 50      | 100 |
| `nome`      | string  | Busca parcial (ILIKE)            | -       | -   |
| `ativo`     | boolean | Filtrar por status ativo/inativo | -       | -   |

**Respostas de Paginação:**

```json
{
  "items": [...],
  "total": 150,
  "page": 1,
  "page_size": 50,
  "total_pages": 3
}
```

---

### Autorização (Cadastro)

| Entidade        | CREATE         | UPDATE         | DELETE         | GET                                  | LIST                                 |
| --------------- | -------------- | -------------- | -------------- | ------------------------------------ | ------------------------------------ |
| Clientes        | Owner, Manager | Owner, Manager | Owner, Manager | Owner, Manager, Accountant, Employee | Owner, Manager, Accountant, Employee |
| Profissionais   | **Owner**      | **Owner**      | **Owner**      | Owner, Manager, Accountant, Employee | Owner, Manager, Accountant, Employee |
| Serviços        | Owner, Manager | Owner, Manager | Owner, Manager | Owner, Manager, Accountant, Employee | Owner, Manager, Accountant, Employee |
| Meios Pagamento | Owner, Manager | Owner, Manager | Owner, Manager | Owner, Manager, Accountant, Employee | Owner, Manager, Accountant, Employee |

**Importante:** Apenas `owner` pode criar/atualizar/deletar **profissionais**.

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

| Status | Código               | Descrição              |
| ------ | -------------------- | ---------------------- |
| 400    | INVALID_REQUEST      | Requisição inválida    |
| 401    | UNAUTHORIZED         | Não autenticado        |
| 403    | FORBIDDEN            | Sem permissão          |
| 404    | NOT_FOUND            | Recurso não encontrado |
| 422    | UNPROCESSABLE_ENTITY | Dados inválidos        |
| 429    | RATE_LIMITED         | Limite de requisições  |
| 500    | INTERNAL_ERROR       | Erro interno           |

---

**Última atualização:** 20/11/2025

## 🎯 Barber Turn (Lista da Vez)

Sistema de rodízio de profissionais baseado em pontos. Permite gerenciar a ordem de atendimento e garantir distribuição justa de clientes entre barbeiros.

### Adicionar Barbeiro à Lista

```http
POST /api/v1/barber-turn/add
Authorization: Bearer {token}
Content-Type: application/json

{
  "professional_id": "uuid-do-profissional"
}

HTTP/1.1 201 Created
{
  "id": "btl-001",
  "tenant_id": "tenant-abc",
  "professional_id": "prof-123",
  "professional_name": "João Silva",
  "current_points": 0,
  "last_turn_at": null,
  "is_active": true,
  "created_at": "2024-11-20T10:30:00Z",
  "updated_at": "2024-11-20T10:30:00Z"
}
```

**Validações:**

- ✅ Profissional deve existir no tenant
- ✅ Profissional deve estar **ATIVO**
- ✅ **Profissional deve ser do tipo BARBEIRO** (não pode ser MANICURE, RECEPCIONISTA, GERENTE ou OUTRO)
- ✅ Profissional não pode estar duplicado na lista

**Camadas de Validação:**

1. **Application Layer:** Validação explícita no use case antes de inserir
2. **Database Trigger:** `validate_professional_type_before_insert` bloqueia inserções inválidas
3. **Frontend UI:** Filtro de seleção exibe apenas profissionais tipo BARBEIRO

**Erros Comuns:**

- `400 BAD_REQUEST` - Profissional não é do tipo BARBEIRO
  ```json
  {
    "code": "BAD_REQUEST",
    "message": "Apenas profissionais do tipo BARBEIRO podem ser adicionados à lista da vez",
    "errors": null
  }
  ```
- `400 BAD_REQUEST` - Profissional já está na lista
- `404 NOT_FOUND` - Profissional não encontrado ou inativo

**Exemplo de Erro - Tipo Inválido:**

```json
{
  "error": "Apenas profissionais do tipo BARBEIRO podem ser adicionados à lista da vez"
}
```

---

### Listar Barbeiros na Fila

```http
GET /api/v1/barber-turn/list
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  "barbers": [
    {
      "id": "btl-001",
      "professional_id": "prof-123",
      "professional_name": "João Silva",
      "current_points": 0,
      "last_turn_at": null,
      "is_active": true,
      "position": 1
    },
    {
      "id": "btl-002",
      "professional_id": "prof-456",
      "professional_name": "Maria Santos",
      "current_points": 2,
      "last_turn_at": "2024-11-20T09:15:00Z",
      "is_active": true,
      "position": 2
    }
  ],
  "total": 2,
  "next_barber": {
    "professional_id": "prof-123",
    "professional_name": "João Silva",
    "current_points": 0
  }
}
```

**Ordenação:** Barbeiros são ordenados por `current_points` ASC, depois por `last_turn_at` ASC NULLS FIRST.

---

### Registrar Atendimento

```http
POST /api/v1/barber-turn/record
Authorization: Bearer {token}
Content-Type: application/json

{
  "professional_id": "prof-123"
}

HTTP/1.1 200 OK
{
  "professional_id": "prof-123",
  "professional_name": "João Silva",
  "previous_points": 0,
  "new_points": 1,
  "last_turn_at": "2024-11-20T11:00:00Z",
  "message": "Atendimento registrado com sucesso"
}
```

**Regras:**

- Incrementa `current_points` em +1
- Atualiza `last_turn_at` para NOW()
- Reordena automaticamente a fila

---

### Pausar/Ativar Barbeiro

```http
PUT /api/v1/barber-turn/{professional_id}/toggle-status
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  "professional_id": "prof-123",
  "professional_name": "João Silva",
  "is_active": false,
  "message": "Barbeiro pausado na fila"
}
```

**Estados:**

- `is_active: true` - Barbeiro recebe atendimentos normalmente
- `is_active: false` - Barbeiro pausado (não sai da lista, apenas para de receber)

---

### Remover Barbeiro da Lista

```http
DELETE /api/v1/barber-turn/{professional_id}
Authorization: Bearer {token}

HTTP/1.1 204 No Content
```

**Comportamento:** Remove completamente da lista (não preserva pontos acumulados).

---

### Reset Mensal (Histórico)

```http
POST /api/v1/barber-turn/reset
Authorization: Bearer {token}

HTTP/1.1 200 OK
{
  "message": "Reset mensal executado com sucesso",
  "snapshot": {
    "month_year": "2024-11",
    "total_barbers": 4,
    "total_points_reset": 25,
    "history_records_created": 4
  }
}
```

**Processo:**

1. Cria snapshot em `barber_turn_history` (mês anterior)
2. Reseta `current_points` para 0 em todos os barbeiros
3. Limpa `last_turn_at`

**Agendamento:** Executado automaticamente via CRON no dia 1º de cada mês às 00:00.

---

### Tipos de Profissionais Permitidos

| Tipo              | Pode Adicionar à Lista? | Motivo                           |
| ----------------- | ----------------------- | -------------------------------- |
| **BARBEIRO**      | ✅ Sim                  | Tipo padrão para rodízio         |
| **MANICURE**      | ❌ Não                  | Não atende clientes de barbearia |
| **RECEPCIONISTA** | ❌ Não                  | Não realiza atendimentos         |
| **GERENTE**       | ❌ Não                  | Função administrativa            |
| **OUTRO**         | ❌ Não                  | Função não definida              |

**Validação:** Implementada via trigger no banco de dados + validação no backend.

---

## ⚠️ Erros
