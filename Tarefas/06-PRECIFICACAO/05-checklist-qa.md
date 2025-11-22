# 🧪 Checklist QA — Precificação

- [ ] Criar/atualizar configuração de precificação e validar restrições de range.
- [ ] Executar simulação com diferentes custos/impostos/comissões e conferir preço sugerido.
- [ ] Verificar persistência do histórico em `precificacao_simulacoes` e filtros por item/tenant/data.
- [ ] UI: simulador salva presets e exibe resultados consistentes com backend.
- [ ] Multi-tenant: configs e históricos não vazam.
- [ ] Performance: simulação <200ms em dados seed.
