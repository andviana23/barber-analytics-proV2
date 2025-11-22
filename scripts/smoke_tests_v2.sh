#!/bin/bash
# Smoke Tests - Validação rápida dos endpoints críticos
# Uso: ./scripts/smoke_tests_v2.sh

# Removido set -e para não parar no primeiro erro
set -o pipefail

API_URL="${API_URL:-http://localhost:8080}"
TENANT_ID="${TENANT_ID:-00000000-0000-0000-0000-000000000001}"

echo "🔥 Iniciando Smoke Tests - Barber Analytics Pro v2.0"
echo "📍 API URL: $API_URL"
echo "🏢 Tenant ID: $TENANT_ID"
echo ""

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para testar endpoint
test_endpoint() {
    local method=$1
    local path=$2
    local expected_code=$3
    local description=$4
    local data=$5

    echo -n "Testing: $description ... "

    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -H "X-Tenant-ID: $TENANT_ID" \
            -d "$data" \
            "$API_URL$path")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "X-Tenant-ID: $TENANT_ID" \
            "$API_URL$path")
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" = "$expected_code" ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} (Expected $expected_code, got $http_code)"
        echo "Response: $body"
        return 1
    fi
}

# Contador de testes
total=0
passed=0
failed=0

# ==================== HEALTH CHECK ====================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 HEALTH CHECK"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

test_endpoint "GET" "/health" "200" "Health Check" && ((passed++)) || ((failed++))
((total++))
echo ""

# ==================== METAS ====================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎯 METAS (Mensal, Barbeiro, Ticket)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Meta Mensal - List
test_endpoint "GET" "/api/v1/metas/monthly?page=1&page_size=10" "200" "List Metas Mensais" && ((passed++)) || ((failed++))
((total++))

# Meta Barbeiro - List
test_endpoint "GET" "/api/v1/metas/barbers?page=1&page_size=10" "200" "List Metas Barbeiro" && ((passed++)) || ((failed++))
((total++))

# Meta Ticket - List
test_endpoint "GET" "/api/v1/metas/ticket?page=1&page_size=10" "200" "List Metas Ticket" && ((passed++)) || ((failed++))
((total++))

echo ""

# ==================== PRICING ====================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💰 PRECIFICAÇÃO"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Get Config (pode retornar 404 se não houver config ainda)
test_endpoint "GET" "/api/v1/pricing/config" "200" "Get Pricing Config" && ((passed++)) || {
    # Se retornar 404, não é erro crítico
    http_code=$(curl -s -w "%{http_code}" -X "GET" \
        -H "X-Tenant-ID: $TENANT_ID" \
        "$API_URL/api/v1/pricing/config" | tail -c 3)
    if [ "$http_code" = "404" ]; then
        echo -e "${YELLOW}⚠ WARN${NC} (Config not found - expected for new tenant)"
        ((passed++))
    else
        ((failed++))
    fi
}
((total++))

# List Simulações
test_endpoint "GET" "/api/v1/pricing/simulations?page=1&page_size=10" "200" "List Simulações" && ((passed++)) || ((failed++))
((total++))

echo ""

# ==================== FINANCIAL ====================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💵 FINANCEIRO"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Contas a Pagar - List
test_endpoint "GET" "/api/v1/financial/payables?page=1&page_size=10" "200" "List Contas a Pagar" && ((passed++)) || ((failed++))
((total++))

# Contas a Receber - List
test_endpoint "GET" "/api/v1/financial/receivables?page=1&page_size=10" "200" "List Contas a Receber" && ((passed++)) || ((failed++))
((total++))

# Compensações - List
test_endpoint "GET" "/api/v1/financial/compensations?page=1&page_size=10" "200" "List Compensações" && ((passed++)) || ((failed++))
((total++))

# Fluxo de Caixa - List
test_endpoint "GET" "/api/v1/financial/cashflow?page=1&page_size=10" "200" "List Fluxo de Caixa" && ((passed++)) || ((failed++))
((total++))

# DRE - List
test_endpoint "GET" "/api/v1/financial/dre?page=1&page_size=10" "200" "List DRE" && ((passed++)) || ((failed++))
((total++))

echo ""

# ==================== SUMMARY ====================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📈 RESUMO DOS TESTES"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Total:  $total testes"
echo -e "${GREEN}Passou: $passed testes${NC}"
echo -e "${RED}Falhou: $failed testes${NC}"

if [ $failed -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Todos os smoke tests passaram!${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}✗ Alguns testes falharam.${NC}"
    exit 1
fi
