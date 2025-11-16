# 📋 Fluxo Manual de Assinaturas

**Versão:** 1.0
**Data:** 15/11/2025
**Status:** Documentação Oficial

---

## 🎯 Objetivo

Documentar o fluxo manual completo para criação, acompanhamento e gestão de assinaturas de barbeiros no Barber Analytics Pro v2, **sem integração automática com Asaas**.

---

## 📌 Visão Geral

O módulo de assinaturas v2 opera em **modo manual**, onde:
- ✅ Planos são cadastrados internamente
- ✅ Assinaturas são criadas manualmente via sistema
- ✅ Faturas (invoices) são geradas manualmente
- ✅ Pagamentos são registrados manualmente
- ✅ Alertas automatizados notificam sobre vencimentos
- ❌ **Não há** integração automática com gateway de pagamento

---

## 🔄 Fluxo Completo

### **Etapa 1: Cadastro de Planos**

**Responsável:** Administrador do Sistema
**Ferramenta:** Frontend → Módulo Planos

#### Passos:
1. Acessar módulo "Planos de Assinatura"
2. Clicar em "Criar Novo Plano"
3. Preencher:
   - Nome do plano (ex: "Plano Barbeiro Mensal")
   - Descrição
   - Valor (R$)
   - Periodicidade (MENSAL, TRIMESTRAL, SEMESTRAL, ANUAL)
4. Salvar

#### Validações:
- ✅ Valor deve ser maior que zero
- ✅ Nome é obrigatório
- ✅ Periodicidade deve ser válida

#### Exemplo:
```json
{
  "nome": "Plano Barbeiro Mensal",
  "descricao": "Acesso completo ao sistema + 10% de comissão",
  "valor": 99.90,
  "periodicidade": "MENSAL"
}
```

---

### **Etapa 2: Criar Assinatura Manual**

**Responsável:** Administrador/Gestor
**Ferramenta:** Frontend → Módulo Assinaturas

#### Passos:
1. Acessar módulo "Assinaturas"
2. Clicar em "Nova Assinatura"
3. Selecionar:
   - Barbeiro (dropdown)
   - Plano (dropdown)
   - Data de início
   - Data da primeira fatura
4. Confirmar criação

#### O que acontece automaticamente:
- ✅ Assinatura criada com `status = ATIVA`
- ✅ Campo `origem_dado = 'manual'`
- ✅ Cálculo automático da `proxima_fatura_data` baseado na periodicidade
- ✅ Nenhuma cobrança é enviada para gateway

#### Exemplo de registro:
```json
{
  "id": "uuid-gerado",
  "tenant_id": "tenant-uuid",
  "plan_id": "plano-uuid",
  "barbeiro_id": "barbeiro-uuid",
  "status": "ATIVA",
  "data_inicio": "2025-11-15",
  "proxima_fatura_data": "2025-12-15",
  "origem_dado": "manual"
}
```

---

### **Etapa 3: Gerar Invoice Manual**

**Responsável:** Sistema (Cron Job) + Manual
**Ferramenta:** Cron Job `ValidateSubscriptions` (02:00) + Frontend

#### Opção A: Geração Automática (Cron)
O cron job `ValidateSubscriptions` executa diariamente às **02:00** e:
1. Busca assinaturas ativas com `proxima_fatura_data <= HOJE`
2. Gera invoices automaticamente com:
   - `status = PENDENTE`
   - `data_vencimento = proxima_fatura_data + 5 dias`
   - `valor = plano.valor`
3. Atualiza `proxima_fatura_data` da assinatura

#### Opção B: Geração Manual
Administrador pode gerar invoice manualmente:
1. Acessar assinatura
2. Clicar em "Gerar Nova Fatura"
3. Confirmar período de competência e valor
4. Sistema cria invoice com `manual = true`

#### Exemplo de invoice:
```json
{
  "id": "invoice-uuid",
  "tenant_id": "tenant-uuid",
  "assinatura_id": "assinatura-uuid",
  "valor": 99.90,
  "status": "PENDENTE",
  "data_vencimento": "2025-12-20",
  "competencia_inicio": "2025-12-15",
  "competencia_fim": "2026-01-14",
  "manual": false
}
```

---

### **Etapa 4: Registrar Pagamento Manual**

**Responsável:** Administrador/Gestor
**Ferramenta:** Frontend → Módulo Assinaturas → Invoice

#### Passos:
1. Acessar invoice pendente
2. Clicar em "Registrar Pagamento"
3. Informar:
   - Data do pagamento (padrão: hoje)
   - Observações (opcional)
4. Confirmar

#### O que acontece:
- ✅ Invoice passa para `status = PAGO`
- ✅ `data_pagamento` é registrada
- ✅ Uma **receita** é criada automaticamente no módulo financeiro:
  - `categoria = "Assinatura Barbeiro"`
  - `valor = invoice.valor`
  - `origem_dado = "assinatura_manual"`
- ✅ Fluxo de caixa é atualizado

#### Exemplo de receita gerada:
```json
{
  "descricao": "Pagamento Assinatura - João Silva (Dez/2025)",
  "valor": 99.90,
  "categoria_id": "uuid-categoria-assinatura",
  "data": "2025-12-18",
  "status": "RECEBIDO",
  "origem_dado": "assinatura_manual",
  "manual": false
}
```

---

### **Etapa 5: Monitoramento e Alertas**

**Responsável:** Sistema (Cron Job)
**Ferramenta:** Cron Job `AlertsJob` (08:00)

#### Alertas Automatizados:
1. **Invoices Vencidas (não pagas)**
   - Detecta invoices com `status = PENDENTE` e `data_vencimento < HOJE`
   - Marca como `VENCIDO`
   - Envia notificação para administrador

2. **Assinaturas Próximas do Vencimento**
   - Detecta assinaturas com `proxima_fatura_data` em até 5 dias
   - Notifica gestor para preparar cobrança

3. **Assinaturas Sem Pagamento Recorrente**
   - Detecta assinaturas com invoices vencidas há mais de 30 dias
   - Sugere suspensão ou cancelamento

---

### **Etapa 6: Cancelamento de Assinatura**

**Responsável:** Administrador/Gestor
**Ferramenta:** Frontend → Módulo Assinaturas

#### Passos:
1. Acessar assinatura ativa
2. Clicar em "Cancelar Assinatura"
3. Confirmar motivo (opcional)
4. Sistema marca assinatura como `CANCELADA`

#### O que acontece:
- ✅ `status = CANCELADA`
- ✅ `data_fim = HOJE`
- ✅ Invoices futuras não são mais geradas
- ✅ Invoices pendentes permanecem para cobrança

---

## 🛠️ Ferramentas de Apoio

### **1. Dashboard de Assinaturas**
- Total de assinaturas ativas/canceladas
- Receita mensal prevista de assinaturas
- Invoices pendentes e vencidas
- Gráficos de histórico

### **2. Relatórios**
- Relatório de pagamentos mensais
- Histórico de assinaturas por barbeiro
- Taxa de inadimplência

### **3. Notificações**
- Email/SMS para barbeiros sobre vencimento
- Alertas para gestores sobre inadimplência

---

## ✅ Checklist de Validação Manual

Após implantação, validar:

- [ ] Plano criado aparece no dropdown
- [ ] Assinatura criada com status ATIVA
- [ ] Primeira invoice gerada automaticamente (cron)
- [ ] Pagamento registrado gera receita financeira
- [ ] Alerta de vencimento notifica gestor
- [ ] Cancelamento impede geração de novas invoices
- [ ] Fluxo de caixa reflete receitas de assinaturas

---

## 🚀 Evolução Futura (Fase 5+)

- [ ] Integração com Asaas para cobrança automática
- [ ] Webhooks para notificação de pagamentos
- [ ] Cobrança recorrente automatizada
- [ ] Suporte a múltiplos gateways

---

**Última Atualização:** 15/11/2025
**Revisado por:** Andrey Viana
**Próxima Revisão:** Após 50 assinaturas cadastradas
