#!/bin/bash

# Teste do fluxo de onboarding (Step 2)

echo "🧪 Testando fluxo de onboarding..."
echo ""

# 1. Fazer login (QA User para pegar token dummy válido em DEV)
echo "1️⃣ Fazendo login (QA User)..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "qa@barberpro.dev", "password": "password"}')

echo "$LOGIN_RESPONSE" | jq .

# Extrair token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Falha ao obter token!"
  exit 1
fi

echo ""
echo "✅ Token obtido: ${TOKEN:0:50}..."
echo ""

# 2. Configurar Tenant (Step 2)
echo "2️⃣ Configurando Tenant (Step 2)..."
CONFIGURE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/onboarding/configure \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d @/tmp/test-configure-tenant.json)

echo "$CONFIGURE_RESPONSE" | jq .

# Verificar resposta
CODE=$(echo "$CONFIGURE_RESPONSE" | jq -r '.code')

echo ""
if [ "$CODE" == "OK" ]; then
  echo "✅ Sucesso! Tenant configurado."
else
  echo "❌ Falha na configuração!"
fi
