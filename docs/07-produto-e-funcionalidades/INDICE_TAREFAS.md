> Criado em: 21/11/2025 18:00 (America/Sao_Paulo)

# 📋 Índice de Execução — Barber Analytics Pro v2.0 (Rebuild)

**Responsável:** Arquiteto-Chefe / Engineer-Lead / PMO  
**Status Geral:** Banco 100% migrado; Backend/Frontend **bloqueados** até concluir B1-B8 (pacote 01).  
**Fonte de verdade:** `DATABASE_MIGRATIONS_COMPLETED.md` + `docs/DIAGRAMA_DEPENDENCIAS_COMPLETO.md`

---

## ⚠️ Pré-requisitos obrigatórios
- Executar **pacote 01-BLOQUEIOS-BASE** antes de qualquer módulo (19 entidades novas, ports, repos, use cases, handlers, cron, services, hooks).
- Respeitar a ordem do diagrama: Bloqueadores → Hardening → Financeiro → Estoque → Metas → Precificação → Lançamento → Monitoramento → Evolução → Agendamentos.
- Validar aderência às migrations 026-038 (tabelas financeiras/metas/precificação + colunas LGPD).

---

## 🗺️ Ordem Oficial de Execução (Diagrama)
```mermaid
flowchart TB
    START([🚀 Início])

    subgraph BLOQ[01 - Bloqueios de Base]
        B1[Domínio 19 entidades]\nB2[Ports]\nB3[Repos sqlc]\nB4[Use Cases]\nB5[HTTP DTO/Handlers]\nB6[Cron Jobs]\nB7[Frontend Services]\nB8[Hooks]
    end

    subgraph HARD[02 - Hardening & OPS]
        H1[T-HAR-001 LGPD]\nH2[T-HAR-002 Backup/DR]
    end

    subgraph FIN[03 - Financeiro]
        F1[Payables]\nF2[Receivables]\nF3[Fluxo Compensado]\nF4[Comissões]\nF5[DRE]\nF6[Dashboard]
    end

    subgraph EST[04 - Estoque]
        E1[Entrada]\nE2[Saída]\nE3[Consumo Automático]\nE4[Inventário]\nE5[Estoque Mínimo]\nE6[Curva ABC]
    end

    subgraph MET[05 - Metas]
        M1[Meta Mensal]\nM2[Meta Barbeiro]\nM3[Meta Ticket]\nM4[Metas Automáticas]
    end

    subgraph PREC[06 - Precificação]
        P1[Simulador]
    end

    subgraph LAN[07 - Lançamento]
        L1[Checklist Pré-GoLive]
    end

    subgraph MON[08 - Monitoramento]
        MON1[Monitoramento + Suporte]
    end

    subgraph EVO[09 - Evolução]
        EV1[PMF/Crescimento]
    end

    subgraph AGE[10 - Agendamentos]
        AG1[Backend + UI + Notificações]
    end

    START --> BLOQ --> HARD --> FIN --> EST --> MET --> PREC --> LAN --> MON --> EVO --> AGE

    classDef blocker fill:#dc2626,stroke:#991b1b,color:#fff
    classDef module fill:#2563eb,stroke:#1d4ed8,color:#fff

    class B1,B2,B3,B4,B5,B6,B7,B8 blocker
    class F1,F2,F3,F4,F5,F6,E1,E2,E3,E4,E5,E6,M1,M2,M3,M4,P1,L1,MON1,EV1,AG1 module
```

---

## 📋 Sequência Detalhada (checklist obrigatório)

### 01 — Bloqueios de Base (B1-B8)
- **Pasta:** `Tarefas/01-BLOQUEIOS-BASE/`
- **Objetivo:** implementar domínio + ports + repos + use cases + HTTP + cron + services/hooks para as tabelas das migrations 026-038.
- **Entregas mínimas:**
  - [ ] 19 entidades novas + VOs (`Money`, `Percentage`, `DMais`, `MesAno`, enums de status).
  - [ ] Ports + repositórios sqlc com filtros por tenant/período/status.
  - [ ] Use cases para DRE, Fluxo, Payables/Receivables, Metas, Precificação, Estoque, Comissões.
  - [ ] DTOs/Handlers HTTP versionados (`/api/v1/...`) + RBAC + validação.
  - [ ] Cron jobs: DRE mensal, Fluxo diário, Compensações, Notificações payables, Estoque mínimo, Comissões.
  - [ ] Services frontend + hooks React Query (dre, fluxo, payables, receivables, metas, precificação, estoque).
- **Sprint sugerida:** 11-12.  
- **Referência:** `Tarefas/CONCLUIR/*`, `Tarefas/01-BLOQUEIOS-BASE/02-backlog.md`.

### 02 — Hardening & OPS (LGPD + Backup/DR)
- **Pasta:** `Tarefas/02-HARDENING-OPS/`
- **Objetivo:** concluir T-LGPD-001 e T-OPS-005 antes de abrir módulos.
- **Entregas mínimas:**
  - [ ] Endpoints LGPD (`/me/preferences`, `/me/export`, `/me` delete) + banner `/privacy` + auditoria.
  - [ ] Workflow de backup (GH Actions + S3 + PITR Neon) + teste de restore documentado.
  - [ ] Observabilidade aplicada aos novos endpoints (métricas, alertas, rate limit).
- **Sprint sugerida:** 12.

### 03 — Financeiro
- **Pasta:** `Tarefas/03-FINANCEIRO/`
- **Objetivo:** completar módulo financeiro avançado.
- **Ordem:** Payables → Receivables → Fluxo Compensado → Comissões → DRE → Dashboard.
- **Referências:** `Tarefas/FINANCEIRO/*.md`, `Tarefas/03-FINANCEIRO/02-backlog.md`.

### 04 — Estoque
- **Pasta:** `Tarefas/04-ESTOQUE/`
- **Objetivo:** implementar controle de estoque integrado ao financeiro.
- **Ordem:** Entrada → Saída → Consumo Automático → Inventário → Estoque Mínimo → Curva ABC.

### 05 — Metas
- **Pasta:** `Tarefas/05-METAS/`
- **Objetivo:** metas mensais, por barbeiro, ticket médio e metas automáticas (tabelas `metas_*`).
- **Ordem:** Meta Mensal → Meta Barbeiro → Meta Ticket → Metas Automáticas.

### 06 — Precificação
- **Pasta:** `Tarefas/06-PRECIFICACAO/`
- **Objetivo:** simulador de precificação com config por tenant (`precificacao_config`/`precificacao_simulacoes`).

### 07 — Lançamento (Fase 7)
- **Pasta:** `Tarefas/07-LANCAMENTO/`
- **Objetivo:** checklist pré-go-live, deploy e monitoramento inicial (T-LAUNCH-001..006).

### 08 — Monitoramento (Fase 8)
- **Pasta:** `Tarefas/08-MONITORAMENTO/`
- **Objetivo:** 4 semanas de estabilização pós-go-live (T-MON-001..008).

### 09 — Evolução (Fase 9)
- **Pasta:** `Tarefas/09-EVOLUCAO/`
- **Objetivo:** ciclos contínuos focados em PMF, crescimento e excelência operacional.

### 10 — Agendamentos (Fase 10)
- **Pasta:** `Tarefas/10-AGENDAMENTOS/`
- **Objetivo:** módulo completo de agenda/DayPilot/notificações conforme `FASE_10_AGENDAMENTOS.md`.

### 📂 Documentos de fase (reposicionados)
- `Tarefas/01-BLOQUEIOS-BASE/FASE_5_MIGRACAO.md`
- `Tarefas/02-HARDENING-OPS/FASE_6_HARDENING.md`
- `Tarefas/07-LANCAMENTO/FASE_7_LANCAMENTO.md`
- `Tarefas/08-MONITORAMENTO/FASE_8_MONITORING.md`
- `Tarefas/09-EVOLUCAO/FASE_9_EVOLUCAO.md`
- `Tarefas/10-AGENDAMENTOS/FASE_10_AGENDAMENTOS.md`

---

## ✅ Critérios de Aceitação do Índice
- Cada pacote (01-10) possui README + contexto + backlog + sprint-plan + checklists dev/qa.
- Nenhuma tarefa é iniciada fora da ordem do diagrama.
- Todas as entregas referenciam as migrations e documentos oficiais (PRD, ARQUITETURA, FLUXOS_CRITICOS, MODELO_DE_DADOS).
- Go-live apenas após: LGPD + Backup/DR validados; módulos Financeiro/Estoque/Metas/Precificação completos; checklist T-LAUNCH-002 aprovado.

---

## 🔗 Referências Rápidas
- `docs/DIAGRAMA_DEPENDENCIAS_COMPLETO.md`
- `ROADMAP_COMPLETO_V2.0.md`
- `CATALOGO_FUNCIONALIDADES.md`
- `PRD-BAP-v2.md`
- `docs/02-arquitetura/ARQUITETURA.md`
- `docs/02-arquitetura/FLUXOS_CRITICOS_SISTEMA.md`
- `docs/02-arquitetura/MODELO_DE_DADOS.md`
