# 📦 Módulo de Estoque

**Versão:** 1.0  
**Data:** 14/11/2025  
**Status:** Design (Fase Posterior)

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Entidades de Domínio](#entidades-de-domínio)
3. [Casos de Uso](#casos-de-uso)
4. [Schema do Banco](#schema-do-banco)
5. [Integrações](#integrações)

---

## 🎯 Visão Geral

Módulo de **gerenciamento de estoque de produtos** (shampoos, cremes, lâminas, etc.). Controla:

- Cadastro de produtos com SKU
- Movimentações de entrada (compra) e saída (uso/venda)
- Alertas de estoque baixo
- Controle de fornecedores
- Histórico e auditoria

**Prioridade:** Baixa (implementar após Módulo Financeiro + Assinaturas estarem estáveis)

---

## 🏛️ Entidades de Domínio

### Produto (Aggregate Root)

```go
type Produto struct {
    ID              string
    TenantID        string
    SKU             string  // Código único por tenant
    Nome            string
    Descricao       string
    Categoria       string  // Xampú, Creme, Lâmina, etc
    QuantidadeAtual int
    QuantidadeMinima int
    ValorUnitario   decimal.Decimal
    UnidadeMedida   string  // Unidade (mL, g, un)
    Ativa           bool
    CriadoEm        time.Time
    AtualizadoEm    time.Time
}

func (p *Produto) Validate() error {
    if p.SKU == "" {
        return ErrProdutoSKURequired
    }
    if p.QuantidadeAtual < 0 {
        return ErrQuantidadeInvalida
    }
    return nil
}

func (p *Produto) EstaBaixo() bool {
    return p.QuantidadeAtual <= p.QuantidadeMinima
}

func (p *Produto) Adicionar(quantidade int) error {
    if quantidade <= 0 {
        return ErrQuantidadeInvalida
    }
    p.QuantidadeAtual += quantidade
    p.AtualizadoEm = time.Now()
    return nil
}

func (p *Produto) Remover(quantidade int) error {
    if quantidade <= 0 {
        return ErrQuantidadeInvalida
    }
    if p.QuantidadeAtual < quantidade {
        return ErrQuantidadeInsuficiente
    }
    p.QuantidadeAtual -= quantidade
    p.AtualizadoEm = time.Now()
    return nil
}
```

### Movimentacao (Entity)

```go
type TipoMovimentacao string

const (
    MovEntrada   TipoMovimentacao = "ENTRADA"
    MovSaida     TipoMovimentacao = "SAIDA"
    MovAjuste    TipoMovimentacao = "AJUSTE"
    MovDevolucao TipoMovimentacao = "DEVOLUCAO"
)

type Movimentacao struct {
    ID              string
    TenantID        string
    ProdutoID       string
    Tipo            TipoMovimentacao
    Quantidade      int
    MotivacaoID     string  // FK para Fornecedor (entrada) ou Barbeiro (saída)
    Observacoes     string
    CriadoEm        time.Time
}
```

### Fornecedor (Entity)

```go
type Fornecedor struct {
    ID          string
    TenantID    string
    RazaoSocial string
    CNPJ        string
    Email       string
    Telefone    string
    Endereco    Endereco
    Ativo       bool
    CriadoEm    time.Time
}
```

---

## 💼 Casos de Uso Futuros

- `CreateProdutoUseCase`
- `ListProdutosUseCase`
- `AtualizarQuantidadeUseCase`
- `AlertarEstoqueBaixoUseCase` (Cron diário)
- `GerarRelatorioEstoqueUseCase`
- `CreateFornecedorUseCase`

---

## 🗄️ Schema do Banco

```sql
CREATE TABLE produtos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    categoria VARCHAR(100),
    quantidade_atual INT NOT NULL DEFAULT 0 CHECK (quantidade_atual >= 0),
    quantidade_minima INT DEFAULT 0,
    valor_unitario DECIMAL(10, 2) NOT NULL CHECK (valor_unitario > 0),
    unidade_medida VARCHAR(20),
    ativa BOOLEAN DEFAULT true,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, sku)
);

CREATE TABLE movimentacoes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    produto_id UUID NOT NULL REFERENCES produtos(id) ON DELETE RESTRICT,
    tipo VARCHAR(50) NOT NULL,
    quantidade INT NOT NULL CHECK (quantidade > 0),
    motivacao_id VARCHAR(255),
    observacoes TEXT,
    criado_em TIMESTAMP DEFAULT NOW()
);

CREATE TABLE fornecedores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    razao_social VARCHAR(255) NOT NULL,
    cnpj VARCHAR(14) UNIQUE,
    email VARCHAR(255),
    telefone VARCHAR(20),
    endereco_rua VARCHAR(255),
    endereco_numero INT,
    endereco_cidade VARCHAR(100),
    ativo BOOLEAN DEFAULT true,
    criado_em TIMESTAMP DEFAULT NOW()
);
```

---

**Status:** 📅 Agendado para Fase 5+
