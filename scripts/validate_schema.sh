#!/usr/bin/env bash

################################################################################
# validate_schema.sh - Script de Validação do Schema PostgreSQL
#
# Descrição:
#   Audita o schema do banco de dados PostgreSQL para garantir que todas
#   as tabelas, índices, constraints e políticas RLS estão no lugar.
#
# Uso:
#   ./scripts/validate_schema.sh <DATABASE_URL>
#
# Exemplo:
#   ./scripts/validate_schema.sh "postgresql://user:pass@host:5432/dbname"
#
# Requisitos:
#   - psql (cliente PostgreSQL)
#   - jq (opcional, para output formatado)
#
# Autor: Andrey Viana
# Versão: 1.0.0
################################################################################

set -euo pipefail

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Contadores
PASSED=0
FAILED=0
WARNINGS=0

# Função de log
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
    ((PASSED++))
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
    ((FAILED++))
}

log_warning() {
    echo -e "${YELLOW}[⚠]${NC} $1"
    ((WARNINGS++))
}

# Função para executar query
run_query() {
    local query="$1"
    psql "$DATABASE_URL" -t -c "$query" 2>/dev/null || echo ""
}

# Função para verificar se tabela existe
check_table() {
    local table_name="$1"
    local result=$(run_query "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = '$table_name');")

    if [[ "$result" =~ "t" ]]; then
        log_success "Tabela '$table_name' existe"
        return 0
    else
        log_error "Tabela '$table_name' NÃO existe"
        return 1
    fi
}

# Função para verificar se índice existe
check_index() {
    local index_name="$1"
    local result=$(run_query "SELECT EXISTS (SELECT FROM pg_indexes WHERE schemaname = 'public' AND indexname = '$index_name');")

    if [[ "$result" =~ "t" ]]; then
        log_success "Índice '$index_name' existe"
        return 0
    else
        log_error "Índice '$index_name' NÃO existe"
        return 1
    fi
}

# Função para verificar RLS
check_rls() {
    local table_name="$1"
    local result=$(run_query "SELECT relrowsecurity FROM pg_class WHERE relname = '$table_name';")

    if [[ "$result" =~ "t" ]]; then
        log_success "RLS habilitado na tabela '$table_name'"
        return 0
    else
        log_warning "RLS NÃO habilitado na tabela '$table_name'"
        return 1
    fi
}

# Função para verificar constraint
check_constraint() {
    local table_name="$1"
    local constraint_name="$2"
    local result=$(run_query "SELECT EXISTS (SELECT FROM information_schema.table_constraints WHERE table_schema = 'public' AND table_name = '$table_name' AND constraint_name = '$constraint_name');")

    if [[ "$result" =~ "t" ]]; then
        log_success "Constraint '$constraint_name' existe na tabela '$table_name'"
        return 0
    else
        log_error "Constraint '$constraint_name' NÃO existe na tabela '$table_name'"
        return 1
    fi
}

# Função para verificar coluna obrigatória
check_column() {
    local table_name="$1"
    local column_name="$2"
    local result=$(run_query "SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_schema = 'public' AND table_name = '$table_name' AND column_name = '$column_name');")

    if [[ "$result" =~ "t" ]]; then
        log_success "Coluna '$column_name' existe na tabela '$table_name'"
        return 0
    else
        log_error "Coluna '$column_name' NÃO existe na tabela '$table_name'"
        return 1
    fi
}

# Validação de argumentos
if [ $# -eq 0 ]; then
    echo "Uso: $0 <DATABASE_URL>"
    echo ""
    echo "Exemplo:"
    echo "  $0 \"postgresql://user:pass@localhost:5432/barber\""
    exit 1
fi

DATABASE_URL="$1"

# Banner
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  🔍 Barber Analytics Pro - Validação de Schema PostgreSQL"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Verificar se psql está instalado
if ! command -v psql &> /dev/null; then
    log_error "psql não está instalado. Instale o PostgreSQL client."
    exit 1
fi

# Testar conexão
log_info "Testando conexão com o banco de dados..."
if ! psql "$DATABASE_URL" -c "SELECT 1;" &> /dev/null; then
    log_error "Falha ao conectar ao banco de dados"
    exit 1
fi
log_success "Conexão com banco de dados estabelecida"

echo ""
echo "───────────────────────────────────────────────────────────────"
echo "  📊 Validando Tabelas Core"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Lista de tabelas esperadas
TABLES=(
    "tenants"
    "users"
    "categorias"
    "receitas"
    "despesas"
    "planos_assinatura"
    "assinaturas"
    "assinatura_invoices"
    "audit_logs"
    "feature_flags"
    "schema_migrations"
)

for table in "${TABLES[@]}"; do
    check_table "$table"
done

echo ""
echo "───────────────────────────────────────────────────────────────"
echo "  🔑 Validando Colunas Obrigatórias (Multi-tenant)"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Verificar colunas tenant_id nas tabelas principais
TENANT_TABLES=(
    "users"
    "categorias"
    "receitas"
    "despesas"
    "planos_assinatura"
    "assinaturas"
    "assinatura_invoices"
    "audit_logs"
)

for table in "${TENANT_TABLES[@]}"; do
    check_column "$table" "tenant_id"
    check_column "$table" "criado_em"
    check_column "$table" "atualizado_em" || check_column "$table" "criado_em"  # Algumas tabelas podem não ter atualizado_em
done

echo ""
echo "───────────────────────────────────────────────────────────────"
echo "  📇 Validando Índices de Performance"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Lista de índices esperados
INDEXES=(
    "idx_users_tenant_id"
    "idx_users_email"
    "idx_categorias_tenant_tipo"
    "idx_receitas_tenant_id"
    "idx_receitas_tenant_data"
    "idx_receitas_tenant_categoria"
    "idx_receitas_tenant_status"
    "idx_despesas_tenant_id"
    "idx_despesas_tenant_data"
    "idx_despesas_tenant_status"
    "idx_assinaturas_tenant"
    "idx_assinaturas_status"
    "idx_assinaturas_barbeiro"
    "idx_invoices_tenant"
    "idx_invoices_status"
    "idx_invoices_vencimento"
    "idx_audit_logs_tenant_timestamp"
)

for index in "${INDEXES[@]}"; do
    check_index "$index" || true  # Não falhamos por falta de índice
done

echo ""
echo "───────────────────────────────────────────────────────────────"
echo "  🔒 Validando Row-Level Security (RLS)"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Tabelas críticas que devem ter RLS
RLS_TABLES=(
    "users"
    "receitas"
    "despesas"
    "assinaturas"
    "assinatura_invoices"
)

for table in "${RLS_TABLES[@]}"; do
    check_rls "$table" || true  # Não falhamos, apenas warning
done

echo ""
echo "───────────────────────────────────────────────────────────────"
echo "  🔗 Validando Constraints e Foreign Keys"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Verificar algumas constraints críticas
check_constraint "users" "users_tenant_id_email_key" || log_warning "Constraint UNIQUE(tenant_id, email) pode estar com nome diferente"
check_constraint "categorias" "categorias_tenant_id_nome_key" || log_warning "Constraint UNIQUE(tenant_id, nome) pode estar com nome diferente"

echo ""
echo "───────────────────────────────────────────────────────────────"
echo "  📈 Validando Migrations"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Verificar tabela de migrations
if check_table "schema_migrations"; then
    MIGRATION_COUNT=$(run_query "SELECT COUNT(*) FROM schema_migrations;" | tr -d ' ')
    log_info "Total de migrations aplicadas: $MIGRATION_COUNT"

    if [ "$MIGRATION_COUNT" -gt 0 ]; then
        LATEST_VERSION=$(run_query "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;" | tr -d ' ')
        log_info "Última versão de migration: $LATEST_VERSION"
    fi
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  📋 Resumo da Validação"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo -e "  ${GREEN}Aprovados:${NC} $PASSED"
echo -e "  ${RED}Falhados:${NC}  $FAILED"
echo -e "  ${YELLOW}Avisos:${NC}    $WARNINGS"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ Schema do banco de dados está válido!${NC}"
    exit 0
else
    echo -e "${RED}✗ Schema do banco de dados tem problemas que precisam ser corrigidos.${NC}"
    exit 1
fi
