#!/bin/bash

# Script de teste E2E para MetaMensal CRUD
# Testa vertical slice completo: POST → GET → LIST → PUT → DELETE

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
TENANT_ID="${TENANT_ID:-00000000-0000-0000-0000-000000000001}"

echo "🧪 Teste E2E - Meta Mensal CRUD"
echo "================================"
echo "Base URL: $BASE_URL"
echo "Tenant ID: $TENANT_ID"
echo ""

# 1. CREATE - POST /api/v1/metas/monthly
echo "📝 1. Criando meta mensal (POST)..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/metas/monthly" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "mes_ano": "2024-12",
    "meta_faturamento": "50000.00",
    "origem": "MANUAL"
  }')

echo "$CREATE_RESPONSE" | jq .

META_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id')
if [ "$META_ID" = "null" ] || [ -z "$META_ID" ]; then
  echo "❌ Erro: ID não retornado na criação"
  exit 1
fi
echo "✅ Meta criada com ID: $META_ID"
echo ""

# 2. GET - GET /api/v1/metas/monthly/:id
echo "🔍 2. Buscando meta mensal (GET)..."
GET_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/metas/monthly/$META_ID" \
  -H "X-Tenant-ID: $TENANT_ID")

echo "$GET_RESPONSE" | jq .

GET_ID=$(echo "$GET_RESPONSE" | jq -r '.id')
if [ "$GET_ID" != "$META_ID" ]; then
  echo "❌ Erro: ID retornado diferente do esperado"
  exit 1
fi
echo "✅ Meta encontrada corretamente"
echo ""

# 3. LIST - GET /api/v1/metas/monthly
echo "📋 3. Listando metas mensais (LIST)..."
LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/metas/monthly" \
  -H "X-Tenant-ID: $TENANT_ID")

echo "$LIST_RESPONSE" | jq .

LIST_COUNT=$(echo "$LIST_RESPONSE" | jq 'length')
if [ "$LIST_COUNT" -lt 1 ]; then
  echo "❌ Erro: Nenhuma meta retornada na listagem"
  exit 1
fi
echo "✅ Listagem retornou $LIST_COUNT meta(s)"
echo ""

# 4. UPDATE - PUT /api/v1/metas/monthly/:id
echo "✏️  4. Atualizando meta mensal (PUT)..."
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/v1/metas/monthly/$META_ID" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "mes_ano": "2024-12",
    "meta_faturamento": "75000.00",
    "origem": "MANUAL"
  }')

echo "$UPDATE_RESPONSE" | jq .

UPDATED_VALOR=$(echo "$UPDATE_RESPONSE" | jq -r '.meta_faturamento')
if [ "$UPDATED_VALOR" != "75000.00" ]; then
  echo "❌ Erro: Valor não foi atualizado corretamente"
  echo "   Esperado: 75000.00, Recebido: $UPDATED_VALOR"
  exit 1
fi
echo "✅ Meta atualizada com sucesso"
echo ""

# 5. DELETE - DELETE /api/v1/metas/monthly/:id
echo "🗑️  5. Deletando meta mensal (DELETE)..."
DELETE_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE \
  "$BASE_URL/api/v1/metas/monthly/$META_ID" \
  -H "X-Tenant-ID: $TENANT_ID")

if [ "$DELETE_STATUS" != "204" ]; then
  echo "❌ Erro: Status de deleção incorreto"
  echo "   Esperado: 204, Recebido: $DELETE_STATUS"
  exit 1
fi
echo "✅ Meta deletada com sucesso (Status: $DELETE_STATUS)"
echo ""

# 6. Verificar deleção - GET deve retornar 404 ou erro
echo "🔍 6. Verificando deleção (GET após DELETE)..."
VERIFY_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X GET \
  "$BASE_URL/api/v1/metas/monthly/$META_ID" \
  -H "X-Tenant-ID: $TENANT_ID")

if [ "$VERIFY_STATUS" = "200" ]; then
  echo "❌ Erro: Meta ainda existe após deleção"
  exit 1
fi
echo "✅ Meta não encontrada após deleção (Status: $VERIFY_STATUS)"
echo ""

echo "✅ ========================================="
echo "✅ TODOS OS TESTES PASSARAM! 🎉"
echo "✅ ========================================="
echo ""
echo "Resumo:"
echo "  ✅ CREATE (POST)   - Meta criada"
echo "  ✅ GET             - Meta encontrada"
echo "  ✅ LIST            - Listagem funcionando"
echo "  ✅ UPDATE (PUT)    - Meta atualizada"
echo "  ✅ DELETE          - Meta deletada"
echo "  ✅ VERIFY          - Deleção confirmada"
echo ""
echo "🚀 Vertical Slice MetaMensal 100% funcional!"
