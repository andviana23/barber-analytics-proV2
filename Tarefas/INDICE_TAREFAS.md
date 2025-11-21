# 📋 Índice de Tarefas Pendentes — Barber Analytics Pro v2.0

**Atualização:** 20/11/2025
**Responsável:** Tech Lead / PMO

> As fases 0 a 4 foram concluídas 100% e arquivadas. Este diretório agora mantém apenas as frentes que ainda precisam de acompanhamento.

---

## 🔎 Snapshot Atual

| Fase   | Nome                                | Progresso             | Status               | Observações                                                                                                                                       |
| ------ | ----------------------------------- | --------------------- | -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| **5**  | Preparação Produção (V2 Standalone) | ██████░░ 50% (2/4)    | 🟡 Em Progresso      | Seeds e validação de integridade concluídos. Onboarding Flow (signup + wizard) e Guia de Deploy seguem pendentes.                                 |
| **6**  | Hardening                           | ████████░ 77% (10/13) | 🟡 Em Progresso      | Segurança, observabilidade e performance fechados. Restam **T-LGPD-001** e **T-OPS-005** (Sentry permaneceu como skipped por decisão do produto). |
| **7**  | Lançamento & Go-Live                | 0% (0/6)              | ⏳ Bloqueado pela F6 | Checklist só inicia após LGPD + Backup/DR concluídos; Go-Live segue previsto para 26/12/2025.                                                     |
| **8**  | Monitoramento & Estabilização       | 0%                    | ⏳ Planejado         | Execução nas 4 semanas pós Go-Live.                                                                                                               |
| **9**  | Evolução & Novas Funcionalidades    | 0%                    | ⏳ Planejado         | Roadmap estratégico (relatórios PDF, gráficos, notificações, etc.).                                                                               |
| **10** | Módulo de Agendamentos              | 0%                    | 📋 Em Planejamento   | Documento de requisitos completo aguardando slot de implementação.                                                                                |

---

## 🎯 Prioridades Imediatas

1. **T-PROD-003 — Onboarding Flow (Fase 5)**
   Backend + Frontend `/signup`, onboarding wizard e tutorial de primeiro acesso.
2. **T-PROD-004 — Documentação de Deploy (Fase 5)**
   `docs/DEPLOY_PRODUCTION.md`, scripts `deploy-backend.sh`/`deploy-frontend.sh` e workflow de aprovação.
3. **T-LGPD-001 & T-OPS-005 (Fase 6)**
   Endpoints /me (delete/export), banner de consentimento, política pública, backup automático + teste de restore documentado.

Concluir os itens acima destrava **T-LAUNCH-001** na Fase 7 e mantém a janela de Go-Live em dezembro.

---

## 📂 Arquivos Ativos

- `FASE_5_MIGRACAO.md` — Seeds, validações, onboarding e guia de deploy.
- `FASE_6_HARDENING.md` — Segurança, observabilidade, performance, LGPD e backup.
- `FASE_7_LANCAMENTO.md` — Checklist de Go-Live e plano de comunicação.
- `FASE_8_MONITORING.md` — Operação assistida pós-lançamento.
- `FASE_9_EVOLUCAO.md` — Roadmap evolutivo (Q1/Q2 2026).
- `FASE_10_AGENDAMENTOS.md` — Planejamento detalhado do módulo de agendamentos.
- `INTEGRACAO_ASAAS_PLANO.md` — Guia caso a integração volte para o roadmap.
- `INDICE_TAREFAS_OLD.md` — Registro histórico (read-only).

---

## 🧭 Como Atualizar

1. Abra o arquivo da fase em `Tarefas/FASE_X.md`.
2. Atualize checklists `[ ] → [x]`, porcentagens e observações dentro da fase.
3. Volte a este índice e ajuste a tabela de snapshot/prioridades conforme necessário.

---

## 🗓️ Linha do Tempo Prevista

- **Nov/25:** concluir Fase 5 e pendências da Fase 6.
- **Dez/25:** rodar checklist da Fase 7 e executar o Go-Live.
- **Jan/26:** iniciar Fase 8 (monitoramento + feedback).
- **Fev/26 em diante:** Fase 9 (evolução) e encaixar Fase 10 (agendamentos).

---

## ✅ Histórico

- Fases 0 → 4 concluídas entre 14 e 19/11/2025 (Fundamentos, DevOps, Backend Core, Módulos Backend e Frontend).
- Documentação completa permanece acessível no histórico do repositório para consulta futura.

**Última revisão:** 20/11/2025 — após limpeza da pasta `Tarefas/` e consolidação das pendências reais.
