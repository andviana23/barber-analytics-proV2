# 📚 Tarefas - Barber Analytics Pro v2.0

**Bem-vindo à Central de Tarefas do Projeto!**

---

## 🎯 COMECE AQUI

### 📘 **[00-GUIA_NAVEGACAO.md](./00-GUIA_NAVEGACAO.md)** ← **LEIA PRIMEIRO!**

Este é o **mapa completo** do projeto com:

- ✅ Sequência de execução obrigatória
- ✅ Estrutura de pastas explicada
- ✅ Estimativas de tempo
- ✅ Regras críticas
- ✅ Dashboard de progresso

**👉 Se você é novo, comece por este arquivo!**

---

## 📂 Estrutura Rápida

```
Tarefas/
│
├── 📘 00-GUIA_NAVEGACAO.md          ← VOCÊ DEVE LER ESTE PRIMEIRO!
├── 📋 INDICE_TAREFAS.md              ← Índice oficial + Diagrama Mermaid
├── ✅ DATABASE_MIGRATIONS_COMPLETED.md  ← Banco 100% pronto
├── 📖 INTEGRACAO_ASAAS_PLANO.md      ← Referência Asaas
│
├── 🔴 CONCLUIR/                      ← BLOQUEADOR #1 - Executar PRIMEIRO
├── 🔴 01-BLOQUEIOS-BASE/             ← BLOQUEADOR #2 - Após CONCLUIR
├── 🟡 02-HARDENING-OPS/              ← Após 01
├── 🟢 03-FINANCEIRO/                 ← Após 01 (paralelo com 04-06)
├── 🟢 04-ESTOQUE/                    ← Após 01 (paralelo com 03,05,06)
├── 🟢 05-METAS/                      ← Após 01 (paralelo com 03,04,06)
├── 🟢 06-PRECIFICACAO/               ← Após 01 (paralelo com 03-05)
├── 🔵 07-LANCAMENTO/                 ← Após 02-06
├── 🔵 08-MONITORAMENTO/              ← Após 07
├── 🔵 09-EVOLUCAO/                   ← Após 08
└── 🔵 10-AGENDAMENTOS/               ← Após 09
```

### Legenda de Status:

- 🔴 **BLOQUEADOR** - Deve ser executado ANTES de tudo
- 🟡 **SEQUENCIAL** - Aguarda etapa anterior
- 🟢 **PARALELO** - Pode ser feito em paralelo (após bloqueadores)
- 🔵 **FINAL** - Etapas finais sequenciais

---

## 🚦 Ordem de Execução (SIMPLIFICADA)

```
1. Leia:  00-GUIA_NAVEGACAO.md
          ↓
2. Leia:  INDICE_TAREFAS.md
          ↓
3. Leia:  DATABASE_MIGRATIONS_COMPLETED.md
          ↓
4. Execute: CONCLUIR/ (23 dias)
          ↓
5. Execute: 01-BLOQUEIOS-BASE/ (já incluído nos 23 dias)
          ↓
6. Execute: 02-HARDENING-OPS/ (5-7 dias)
          ↓
7. Execute em PARALELO: 03-FINANCEIRO/ + 04-ESTOQUE/ + 05-METAS/ + 06-PRECIFICACAO/
          ↓
8. Execute sequencialmente: 07 → 08 → 09 → 10
```

---

## 📊 Status Atual

| Componente         | Status         | Progresso |
| ------------------ | -------------- | --------- |
| Banco de Dados     | ✅ Completo    | 100%      |
| Backend (Go)       | 🟡 Parcial     | ~40%      |
| Frontend (Next.js) | 🟡 Parcial     | ~30%      |
| **Bloqueios**      | 🔴 **CRÍTICO** | **0%**    |

---

## 🎯 Próximos Passos

### Se você acabou de chegar:

1. ✅ **Leia** [`00-GUIA_NAVEGACAO.md`](./00-GUIA_NAVEGACAO.md)
2. ✅ **Leia** [`INDICE_TAREFAS.md`](./INDICE_TAREFAS.md)
3. ✅ **Leia** [`DATABASE_MIGRATIONS_COMPLETED.md`](./DATABASE_MIGRATIONS_COMPLETED.md)
4. 🔴 **Execute** [`CONCLUIR/`](./CONCLUIR/) - Tarefas bloqueadoras
5. 🔴 **Execute** [`01-BLOQUEIOS-BASE/`](./01-BLOQUEIOS-BASE/)
6. ✅ **Continue** com módulos 02-10 na ordem

### Se você já está trabalhando:

1. ✅ Verifique o progresso em [`00-GUIA_NAVEGACAO.md`](./00-GUIA_NAVEGACAO.md)
2. ✅ Consulte [`INDICE_TAREFAS.md`](./INDICE_TAREFAS.md) para ver o diagrama
3. ✅ Escolha a próxima tarefa na sequência correta

---

## 📖 Documentos de Referência

| Documento                                                                | Descrição                             |
| ------------------------------------------------------------------------ | ------------------------------------- |
| [`00-GUIA_NAVEGACAO.md`](./00-GUIA_NAVEGACAO.md)                         | **INÍCIO** - Mapa completo do projeto |
| [`INDICE_TAREFAS.md`](./INDICE_TAREFAS.md)                               | Índice oficial + Diagrama Mermaid     |
| [`DATABASE_MIGRATIONS_COMPLETED.md`](./DATABASE_MIGRATIONS_COMPLETED.md) | Status do banco (42 tabelas)          |
| [`INTEGRACAO_ASAAS_PLANO.md`](./INTEGRACAO_ASAAS_PLANO.md)               | Integração pagamentos Asaas           |

---

## 🗂️ Padrão de Arquivos (Cada Pasta)

Dentro de cada pasta `XX-NOME/` você encontrará:

| Arquivo               | Quando Usar                                 |
| --------------------- | ------------------------------------------- |
| `README.md`           | **SEMPRE LER PRIMEIRO** - Overview da etapa |
| `01-contexto.md`      | Antes de planejar - Estado atual            |
| `02-backlog.md`       | Antes de executar - Lista de tarefas        |
| `03-sprint-plan.md`   | Ao iniciar - Ordem de execução              |
| `04-checklist-dev.md` | Durante desenvolvimento                     |
| `05-checklist-qa.md`  | Antes de deploy - Testes                    |

---

## ⚠️ Regras CRÍTICAS

### ❌ NUNCA:

1. Pule `CONCLUIR/` ou `01-BLOQUEIOS-BASE/`
2. Execute módulos 03-10 antes de concluir 01
3. Ignore validações de `tenant_id`
4. Acesse repositório direto de cron (sempre use use case)

### ✅ SEMPRE:

1. Leia `02-backlog.md` antes de começar
2. Valide com `04-checklist-dev.md` antes de "pronto"
3. Execute testes com `05-checklist-qa.md`
4. Mantenha cobertura > 70%

---

## 🆘 Precisa de Ajuda?

- **Arquitetura**: `docs/02-arquitetura/ARQUITETURA.md`
- **Backend Go**: `docs/04-backend/GUIA_DEV_BACKEND.md`
- **Frontend**: `docs/03-frontend/GUIA_FRONTEND.md`
- **Design System**: `docs/03-frontend/DESIGN_SYSTEM.md`
- **IA (Copilot)**: `.github/copilot-instructions.md`

---

**Última Atualização:** 21/11/2025
**Versão:** 2.0

**BOA SORTE! 🚀**
