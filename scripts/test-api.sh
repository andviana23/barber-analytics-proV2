#!/bin/bash

# ============================================================================
# 🧪 Barber Analytics Pro V2 — Test API
# Testa endpoints críticos da API
# ============================================================================

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_URL="http://localhost:8080/api/v1"
TENANT_ID="e2e00000-0000-0000-0000-000000000001"

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║${NC}  🧪 Testando API — Barber Analytics Pro V2"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo

# ============================================================================
# TEST PING
# ============================================================================

echo -e "${YELLOW}📡 Test 1: PING${NC}"
response=$(curl -s "$BASE_URL/ping")
if [[ $response == *"pong"* ]]; then
    echo -e "   ${GREEN}✅ Backend respondendo${NC}"
else
    echo -e "   ${RED}❌ Backend não respondeu${NC}"
    exit 1
fi

echo

# ============================================================================
# TEST HEALTH
# ============================================================================

echo -e "${YELLOW}💚 Test 2: HEALTH CHECK${NC}"
http_code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health")
if [ "$http_code" -eq 200 ]; then
    echo -e "   ${GREEN}✅ Health check OK${NC}"
else
    echo -e "   ${RED}❌ Health check falhou (HTTP $http_code)${NC}"
fi

echo

# ============================================================================
# TEST FEATURE FLAGS
# ============================================================================

echo -e "${YELLOW}🚩 Test 3: FEATURE FLAGS${NC}"
response=$(curl -s -H "X-Tenant-ID: $TENANT_ID" "$BASE_URL/feature-flags")
if [[ $response == *"use_v2_financial"* ]]; then
    echo -e "   ${GREEN}✅ Feature flags acessíveis${NC}"
    echo -e "   ${BLUE}Resposta:${NC} $response" | head -1
else
    echo -e "   ${RED}❌ Feature flags não funcionou${NC}"
fi

echo

# ============================================================================
# TEST RECEITAS
# ============================================================================

echo -e "${YELLOW}💰 Test 4: LIST RECEITAS${NC}"
http_code=$(curl -s -o /dev/null -w "%{http_code}" -H "X-Tenant-ID: $TENANT_ID" "$BASE_URL/receitas")
if [ "$http_code" -eq 200 ]; then
    echo -e "   ${GREEN}✅ Receitas endpoint OK${NC}"
else
    echo -e "   ${RED}❌ Receitas falhou (HTTP $http_code)${NC}"
fi

echo

# ============================================================================
# TEST METRICS
# ============================================================================

echo -e "${YELLOW}📊 Test 5: PROMETHEUS METRICS${NC}"
response=$(curl -s "http://localhost:8080/metrics")
if [[ $response == *"http_requests_total"* ]]; then
    echo -e "   ${GREEN}✅ Prometheus metrics ativo${NC}"
    echo -e "   ${BLUE}Métricas encontradas:${NC}"
    echo "$response" | grep "^http_" | head -3 | sed 's/^/      /'
else
    echo -e "   ${RED}❌ Prometheus metrics não funcionou${NC}"
fi

echo

# ============================================================================
# TEST FRONTEND
# ============================================================================

echo -e "${YELLOW}🌐 Test 6: FRONTEND${NC}"
http_code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:3000")
if [ "$http_code" -eq 200 ]; then
    echo -e "   ${GREEN}✅ Frontend respondendo${NC}"
else
    echo -e "   ${RED}❌ Frontend não respondeu (HTTP $http_code)${NC}"
fi

echo

# ============================================================================
# SUMMARY
# ============================================================================

echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ TESTES COMPLETOS!${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
echo

echo -e "🌐 Acesse o frontend: ${BLUE}http://localhost:3000${NC}"
echo -e "📡 Teste a API: ${BLUE}http://localhost:8080/api/v1/ping${NC}"
echo -e "📊 Metricas: ${BLUE}http://localhost:8080/metrics${NC}"
echo
