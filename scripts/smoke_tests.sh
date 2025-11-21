#!/usr/bin/env bash

################################################################################
# smoke_tests.sh - Smoke Tests End-to-End para Barber Analytics Pro
#
# Descrição:
#   Executa testes de ponta-a-ponta simulando um fluxo de usuário real.
#   Testa as funcionalidades críticas da API para garantir que o sistema
#   está operacional.
#
# Uso:
#   ./scripts/smoke_tests.sh [API_URL]
#
# Exemplo:
#   ./scripts/smoke_tests.sh http://localhost:8080
#   ./scripts/smoke_tests.sh https://api.barber-analytics.com
#
# Fluxo Testado:
#   1. Health Check
#   2. Criar Tenant
#   3. Criar Usuário
#   4. Login (obter JWT)
#   5. Criar Receita
#   6. Listar Receitas
#   7. Cleanup (opcional)
#
# Requisitos:
#   - curl
#   - jq
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
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configurações
API_URL="${1:-http://localhost:8080}"
TIMEOUT=10

# Contadores
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_TOTAL=0

# Variáveis globais para armazenar IDs e tokens
TENANT_ID=""
USER_ID=""
JWT_TOKEN=""
RECEITA_ID=""

# Timestamp único para evitar conflitos
TIMESTAMP=$(date +%s)
UNIQUE_SUFFIX="test_${TIMESTAMP}"

# Função de log
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_test() {
    echo -e "${CYAN}[TEST]${NC} $1"
    ((TESTS_TOTAL++))
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
    ((TESTS_PASSED++))
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
    ((TESTS_FAILED++))
}

log_warning() {
    echo -e "${YELLOW}[⚠]${NC} $1"
}

# Função para fazer requisições HTTP
http_request() {
    local method="$1"
    local endpoint="$2"
    local data="${3:-}"
    local auth_header="${4:-}"

    local full_url="${API_URL}${endpoint}"

    if [ -n "$data" ]; then
        if [ -n "$auth_header" ]; then
            curl -s -w "\n%{http_code}" -X "$method" "$full_url" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $auth_header" \
                -H "X-Tenant-ID: $TENANT_ID" \
                -d "$data" \
                --max-time "$TIMEOUT"
        else
            curl -s -w "\n%{http_code}" -X "$method" "$full_url" \
                -H "Content-Type: application/json" \
                -d "$data" \
                --max-time "$TIMEOUT"
        fi
    else
        if [ -n "$auth_header" ]; then
            curl -s -w "\n%{http_code}" -X "$method" "$full_url" \
                -H "Authorization: Bearer $auth_header" \
                -H "X-Tenant-ID: $TENANT_ID" \
                --max-time "$TIMEOUT"
        else
            curl -s -w "\n%{http_code}" -X "$method" "$full_url" \
                --max-time "$TIMEOUT"
        fi
    fi
}

# Função para extrair código HTTP da resposta
get_http_code() {
    echo "$1" | tail -n1
}

# Função para extrair body da resposta
get_response_body() {
    echo "$1" | sed '$d'
}

# Banner
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  🧪 Barber Analytics Pro - Smoke Tests E2E"
echo "═══════════════════════════════════════════════════════════════"
echo ""
log_info "API URL: $API_URL"
echo ""

# Verificar dependências
log_info "Verificando dependências..."
if ! command -v curl &> /dev/null; then
    log_error "curl não está instalado"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    log_warning "jq não está instalado (parsing JSON será limitado)"
fi

log_success "Dependências verificadas"
echo ""

# ============================================================================
# TESTE 1: Health Check
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "1. Health Check"
echo "───────────────────────────────────────────────────────────────"

RESPONSE=$(http_request "GET" "/health")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 200 ]; then
    log_success "Health check passou (HTTP 200)"

    if command -v jq &> /dev/null; then
        STATUS=$(echo "$BODY" | jq -r '.status // empty')
        DB_STATUS=$(echo "$BODY" | jq -r '.database.connected // empty')

        if [ "$STATUS" = "healthy" ] || [ "$STATUS" = "degraded" ]; then
            log_success "Status do sistema: $STATUS"
        else
            log_warning "Status inesperado: $STATUS"
        fi

        if [ "$DB_STATUS" = "true" ]; then
            log_success "Banco de dados conectado"
        else
            log_error "Banco de dados não conectado"
        fi
    fi
else
    log_error "Health check falhou (HTTP $HTTP_CODE)"
    log_error "Response: $BODY"
    exit 1
fi

echo ""

# ============================================================================
# TESTE 2: Criar Tenant
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "2. Criar Tenant"
echo "───────────────────────────────────────────────────────────────"

TENANT_DATA=$(cat <<EOF
{
    "nome": "Barbearia ${UNIQUE_SUFFIX}",
    "cnpj": "12345678000199",
    "plano": "pro"
}
EOF
)

RESPONSE=$(http_request "POST" "/api/v1/tenants" "$TENANT_DATA")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 201 ] || [ "$HTTP_CODE" -eq 200 ]; then
    log_success "Tenant criado (HTTP $HTTP_CODE)"

    if command -v jq &> /dev/null; then
        TENANT_ID=$(echo "$BODY" | jq -r '.data.id // .id // empty')

        if [ -n "$TENANT_ID" ]; then
            log_success "Tenant ID: $TENANT_ID"
        else
            log_error "Não foi possível extrair tenant_id"
            log_error "Response: $BODY"
            exit 1
        fi
    else
        log_warning "jq não disponível, não é possível extrair tenant_id automaticamente"
        exit 1
    fi
else
    log_error "Falha ao criar tenant (HTTP $HTTP_CODE)"
    log_error "Response: $BODY"
    exit 1
fi

echo ""

# ============================================================================
# TESTE 3: Criar Usuário
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "3. Criar Usuário"
echo "───────────────────────────────────────────────────────────────"

USER_DATA=$(cat <<EOF
{
    "tenant_id": "$TENANT_ID",
    "nome": "Usuario Teste ${UNIQUE_SUFFIX}",
    "email": "teste_${UNIQUE_SUFFIX}@example.com",
    "password": "SenhaSegura123!",
    "role": "owner"
}
EOF
)

RESPONSE=$(http_request "POST" "/api/v1/users" "$USER_DATA")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 201 ] || [ "$HTTP_CODE" -eq 200 ]; then
    log_success "Usuário criado (HTTP $HTTP_CODE)"

    if command -v jq &> /dev/null; then
        USER_ID=$(echo "$BODY" | jq -r '.data.id // .id // empty')

        if [ -n "$USER_ID" ]; then
            log_success "User ID: $USER_ID"
        else
            log_warning "Não foi possível extrair user_id"
        fi
    fi
else
    log_error "Falha ao criar usuário (HTTP $HTTP_CODE)"
    log_error "Response: $BODY"
    # Continuamos mesmo se falhar (usuário pode já existir)
fi

echo ""

# ============================================================================
# TESTE 4: Login
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "4. Login (obter JWT)"
echo "───────────────────────────────────────────────────────────────"

LOGIN_DATA=$(cat <<EOF
{
    "email": "teste_${UNIQUE_SUFFIX}@example.com",
    "password": "SenhaSegura123!"
}
EOF
)

RESPONSE=$(http_request "POST" "/api/v1/auth/login" "$LOGIN_DATA")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 200 ]; then
    log_success "Login realizado com sucesso (HTTP 200)"

    if command -v jq &> /dev/null; then
        JWT_TOKEN=$(echo "$BODY" | jq -r '.data.access_token // .access_token // .token // empty')

        if [ -n "$JWT_TOKEN" ] && [ "$JWT_TOKEN" != "null" ]; then
            log_success "JWT Token obtido"
            # Mostra apenas primeiros caracteres do token
            TOKEN_PREVIEW="${JWT_TOKEN:0:20}..."
            log_info "Token: $TOKEN_PREVIEW"
        else
            log_error "Não foi possível extrair JWT token"
            log_error "Response: $BODY"
            exit 1
        fi
    fi
else
    log_error "Falha no login (HTTP $HTTP_CODE)"
    log_error "Response: $BODY"
    exit 1
fi

echo ""

# ============================================================================
# TESTE 5: Criar Receita
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "5. Criar Receita (autenticado)"
echo "───────────────────────────────────────────────────────────────"

RECEITA_DATA=$(cat <<EOF
{
    "descricao": "Corte de cabelo - Smoke Test",
    "valor": 50.00,
    "categoria_id": "00000000-0000-0000-0000-000000000001",
    "metodo_pagamento": "DINHEIRO",
    "data": "$(date +%Y-%m-%d)",
    "status": "CONFIRMADO"
}
EOF
)

RESPONSE=$(http_request "POST" "/api/v1/receitas" "$RECEITA_DATA" "$JWT_TOKEN")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 201 ] || [ "$HTTP_CODE" -eq 200 ]; then
    log_success "Receita criada (HTTP $HTTP_CODE)"

    if command -v jq &> /dev/null; then
        RECEITA_ID=$(echo "$BODY" | jq -r '.data.id // .id // empty')

        if [ -n "$RECEITA_ID" ]; then
            log_success "Receita ID: $RECEITA_ID"
        fi
    fi
else
    log_warning "Falha ao criar receita (HTTP $HTTP_CODE) - pode ser esperado se categoria não existe"
    log_info "Response: $BODY"
fi

echo ""

# ============================================================================
# TESTE 6: Listar Receitas
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "6. Listar Receitas (autenticado)"
echo "───────────────────────────────────────────────────────────────"

RESPONSE=$(http_request "GET" "/api/v1/receitas" "" "$JWT_TOKEN")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 200 ]; then
    log_success "Receitas listadas com sucesso (HTTP 200)"

    if command -v jq &> /dev/null; then
        RECEITAS_COUNT=$(echo "$BODY" | jq '.data | length // 0')
        log_info "Total de receitas retornadas: $RECEITAS_COUNT"
    fi
else
    log_error "Falha ao listar receitas (HTTP $HTTP_CODE)"
    log_error "Response: $BODY"
fi

echo ""

# ============================================================================
# Resumo Final
# ============================================================================
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 Resumo dos Smoke Tests"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo -e "  Total de testes:   $TESTS_TOTAL"
echo -e "  ${GREEN}Aprovados:${NC}         $TESTS_PASSED"
echo -e "  ${RED}Falhados:${NC}          $TESTS_FAILED"
echo ""

SUCCESS_RATE=$((TESTS_PASSED * 100 / TESTS_TOTAL))

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ Todos os smoke tests passaram! Sistema está operacional.${NC}"
    echo ""
    exit 0
elif [ $SUCCESS_RATE -ge 80 ]; then
    echo -e "${YELLOW}⚠ Alguns testes falharam, mas taxa de sucesso é ${SUCCESS_RATE}% (aceitável).${NC}"
    echo ""
    exit 0
else
    echo -e "${RED}✗ Muitos testes falharam. Sistema pode não estar operacional.${NC}"
    echo ""
    exit 1
fi
