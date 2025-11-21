#!/bin/bash

# Script para iniciar backend em modo E2E (sem autenticação JWT)
# Usado apenas para testes automatizados Playwright

set -e

echo "🧪 Iniciando Backend em Modo E2E..."
echo "   ⚠️  ATENÇÃO: JWT authentication DESABILITADO"
echo ""

# Matar processo existente
pkill -f barber-api 2>/dev/null || true
sleep 1

# Ir para o diretório do backend
cd "$(dirname "$0")/../backend"

# Carregar variáveis de ambiente do arquivo .env
if [ -f ".env" ]; then
    echo "📋 Carregando variáveis de ambiente..."
    set -a
    source .env
    set +a
else
    echo "❌ Erro: arquivo .env não encontrado"
    exit 1
fi

# Forçar modo E2E
export E2E_MODE=true

# Compilar
echo "🔨 Compilando backend..."
go build -o bin/barber-api ./cmd/api

# Iniciar
echo "🚀 Iniciando servidor (porta 8080)..."
./bin/barber-api > /tmp/backend-e2e.log 2>&1 &
BACKEND_PID=$!

echo "✅ Backend iniciado (PID: $BACKEND_PID)"
echo "📋 Logs: tail -f /tmp/backend-e2e.log"
echo ""
echo "Aguardando 5 segundos..."
sleep 5

# Verificar se está rodando
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Backend OK!"
    echo "   Modo: E2E (Auth desabilitado)"
    echo "   Tenant: e2e00000-0000-0000-0000-000000000001"
else
    echo "❌ Backend não respondeu ao health check"
    echo ""
    echo "📋 Últimas linhas do log:"
    tail -20 /tmp/backend-e2e.log
    exit 1
fi
