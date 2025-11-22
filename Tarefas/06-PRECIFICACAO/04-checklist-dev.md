# ✅ Checklist Dev — Precificação

- [ ] Domínio criado com validação de ranges (margem 5-100%, markup >=1, impostos/comissão >=0).
- [ ] Repositórios com operações de upsert para config e listagem de simulações por item/tenant.
- [ ] Use case de simulação considera custos diretos, comissões, impostos e margem alvo; suporta presets.
- [ ] Endpoints `/pricing/config` e `/pricing/simulations` com RBAC (owner/manager) e rate limiting.
- [ ] UI do simulador integrada com React Query; export/import de presets.
- [ ] Métricas/logs para simulações executadas.
