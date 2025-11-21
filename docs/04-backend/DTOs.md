> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# Guia de DTOs (Backend Go)

Documenta onde vivem os DTOs, convenções de entrada/saída e como mapear domínio ↔ transporte. Use este guia ao criar novos endpoints ou refatorar validações.

## Onde ficam e escopo
- Código: `backend/internal/application/dto/*.go`.
- Mapeamento domínio ↔ DTO: `backend/internal/application/mapper/*.go`.
- Uso em casos de uso: `backend/internal/application/usecase/**`.
- Handlers HTTP apenas recebem DTO de request, chamam use case e retornam DTO de response.

## Arquivos e grupos principais
| Arquivo (dto) | Escopo | Exemplos de structs |
| --- | --- | --- |
| `auth_dto.go` | Login/refresh/criação de usuário | `LoginRequest`, `LoginResponse`, `UserResponse`, `ErrorResponse` |
| `financial_dto.go` | Receitas/Despesas | `CreateReceitaRequest`, `ReceitaResponse`, `ListReceitasRequest` (query params) |
| `subscription_dto.go` | Planos, assinaturas, invoices | `CreatePlanoAssinaturaRequest`, `CreateAssinaturaRequest`, `CreateInvoiceRequest`, `RegistrarPagamentoRequest`, `AssinaturaResponse` |
| `cadastro_dto.go` | Clientes, profissionais, serviços, produtos, cupons | `CreateClienteRequest`, `UpdateClienteRequest`, `ListClientesRequest`, `HorarioSemanalDTO`, `CreateProdutoRequest`, `AtualizarEstoqueRequest` |
| `barber_turn_dto.go` | Lista da Vez / fila de barbeiros | Requests de atualização de status e responses de fila/histórico |
| `feature_flag_dto.go` | Feature flags por tenant | `SetFeatureFlagRequest`, `FeatureFlagResponse` |
| `flexible_date.go` | Tipagem auxiliar | `FlexibleDate` (aceita formatos de data configuráveis) |

## Convenções de modelagem
- Sufixos: `Request` para entrada, `Response` para saída. Listagens usam `List*Request` com tags `query:"campo"` e defaults.
- Tags: sempre `json:"snake_case"` para body e `validate:"..."` para regras do `go-playground/validator`.
- Datas: inputs que aceitam múltiplos formatos usam `FlexibleDate`; responses retornam `time.Time` ou string formatada conforme domínio.
- Valores monetários: strings no DTO (`Valor string`) para evitar `float64`; parse ocorre no use case.
- Identidade/tenant: `tenant_id` nunca vem do cliente exceto em contextos estritamente controlados; use middleware/claims para preencher.

## Validação e erros
- Validação de entrada feita no use case com `validator/v10`; handlers desalocam apenas requests válidos.
- Erros padronizados: `ErrorResponse` (`code`, `message`, `details`, `trace_id`). Prefira códigos semânticos (ex.: `validation_error`, `not_found`).

## Mapeamento domínio ↔ DTO
- Converters ficam em `backend/internal/application/mapper/` (ex.: `cadastro_mapper.go`, `barber_turn_mapper.go`).
- Padrão: `ToXDTO(domain)` para saída, `FromXDTO(dto)` para entrada que exige value objects.
- Gaps conhecidos: conversão completa de `HorarioSemanalDTO` ainda não implementada (ver TODOs nos use cases de profissional).

## Respostas paginadas e listas
- Use cases de listagem retornam slices de `Response` DTO e podem incorporar paginação via query params (`page`, `page_size` com defaults).
- Para adicionar paginação: padronizar campos `total`, `page`, `page_size`, `data` (seguir exemplos em `financial` e `barber_turn`).

## Checklist para criar/alterar um DTO
- [ ] Nomear com sufixo `Request`/`Response`.
- [ ] Definir tags `json`/`query` + `validate`.
- [ ] Usar `FlexibleDate` para datas de entrada que aceitam string.
- [ ] Evitar `float64` para dinheiro; usar string ou decimal dentro do domínio.
- [ ] Mapear no `mapper`: domínio → DTO (saída) e, se necessário, DTO → domínio (entrada).
- [ ] Cobrir com testes do use case e ajustar handlers para bind/validation.

## Referências
- `docs/02-arquitetura/ARQUITETURA.md` (seção DTO).
- `docs/04-backend/GUIA_DEV_BACKEND.md` (exemplo de criação de use case com DTO).
- `backend/internal/application/dto/` (implementação atual).
- `backend/internal/application/mapper/` (conversores ativos e TODOs).
