#!/bin/bash

# Teste do endpoint /auth/me

echo "🧪 Testando autenticação..."
echo ""

# 1. Fazer login
echo "1️⃣ Fazendo login..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "andrey@tratodebarbados.com", "password": "123456"}')

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

# 2. Testar /auth/me
echo "2️⃣ Testando GET /auth/me..."
ME_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "$ME_RESPONSE" | jq .

# Verificar resposta
STATUS=$(echo "$ME_RESPONSE" | jq -r '.status')
NOME=$(echo "$ME_RESPONSE" | jq -r '.data.nome // .data.name')
EMAIL=$(echo "$ME_RESPONSE" | jq -r '.data.email')

echo ""
if [ "$STATUS" == "success" ]; then
  echo "✅ Sucesso!"
  echo "   Nome: $NOME"
  echo "   Email: $EMAIL"
else
  echo "❌ Falha na requisição!"
fi
