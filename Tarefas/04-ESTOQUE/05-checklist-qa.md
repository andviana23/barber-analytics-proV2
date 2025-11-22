# 🧪 Checklist QA — Estoque

- [ ] Criar entrada com múltiplos itens e verificar atualização de saldos e custo médio.
- [ ] Registrar saída com motivo e validar bloqueio quando saldo insuficiente.
- [ ] Executar consumo automático ao registrar serviço e conferir baixas corretas.
- [ ] Rodar inventário (contagem + ajuste) e validar auditoria dos ajustes.
- [ ] Disparar alerta de estoque mínimo ao simular item abaixo do limiar.
- [ ] Gerar relatório Curva ABC e conferir classificação A/B/C.
- [ ] Multi-tenant: movimentações e inventário isolados.
- [ ] Regressão de integração com financeiro (payable opcional gerado apenas quando configurado).
