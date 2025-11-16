#!/bin/bash

# ============================================================================
# 🚀 Barber Analytics Pro V2 — Start All Services
# Liga Backend + Frontend + Prometheus em um único comando
# ============================================================================

PROJECT_ROOT="/home/andrey/projetos/barber-Analytic-proV2"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# ============================================================================
# HEADER
# ============================================================================

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║${NC}  🚀 Barber Analytics Pro V2 — Iniciando Serviços"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo

# ============================================================================
# SETUP ENVIRONMENT
# ============================================================================

# Credenciais Neon corretas (Projeto: BarberAnalicV2, ID: old-queen-78246613)
export DATABASE_URL="postgresql://neondb_owner:npg_83COkAjHMotv@ep-winter-leaf-adhqz08p-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
export NEXT_PUBLIC_API_URL="http://localhost:8080/api/v1"
export ENV="production"
export LOG_LEVEL="info"

echo -e "${YELLOW}⚙️  Configurando variáveis de ambiente...${NC}"
echo -e "   DATABASE_URL: ${BLUE}Neon PostgreSQL (neondb_owner)${NC}"
echo -e "   BANCO: ${BLUE}neondb${NC} (BarberAnalicV2)"
echo -e "   API_URL: ${BLUE}http://localhost:8080/api/v1${NC}"
echo

# ============================================================================
# KILL PORTAS ANTIGAS
# ============================================================================

echo -e "${YELLOW}🔍 Limpando portas antigas...${NC}"

for port in 8080 3000 9090; do
    if lsof -ti :$port 2>/dev/null | xargs kill -9 2>/dev/null; then
        echo -e "   ${GREEN}✅${NC} Porta $port liberada"
    fi
done

sleep 1
echo

# ============================================================================
# BACKEND
# ============================================================================

echo -e "${YELLOW}🔧 Iniciando Backend (Go)...${NC}"

cd "$BACKEND_DIR"

# Build
echo -e "   ${BLUE}→${NC} Compilando..."
if go build -o bin/barber-api cmd/api/main.go 2>&1 | grep -v "no Go files" | head -5; then
    :
fi

# Run
echo -e "   ${BLUE}→${NC} Iniciando servidor na porta 8080..."
./bin/barber-api > /tmp/backend.log 2>&1 &
BACKEND_PID=$!
echo $BACKEND_PID > /tmp/barber-backend.pid

sleep 2

if lsof -ti :8080 >/dev/null 2>&1; then
    echo -e "   ${GREEN}✅ Backend iniciado${NC} (PID: $BACKEND_PID)"
else
    echo -e "   ${RED}❌ Backend falhou!${NC}"
    cat /tmp/backend.log | tail -10
    exit 1
fi

echo

# ============================================================================
# FRONTEND
# ============================================================================

echo -e "${YELLOW}⚛️  Iniciando Frontend (Next.js)...${NC}"

cd "$FRONTEND_DIR"

echo -e "   ${BLUE}→${NC} Instalando dependências..."
pnpm install --frozen-lockfile > /tmp/pnpm.log 2>&1 || true

echo -e "   ${BLUE}→${NC} Iniciando servidor na porta 3000..."
pnpm dev > /tmp/frontend.log 2>&1 &
FRONTEND_PID=$!
echo $FRONTEND_PID > /tmp/barber-frontend.pid

sleep 3

if lsof -ti :3000 >/dev/null 2>&1; then
    echo -e "   ${GREEN}✅ Frontend iniciado${NC} (PID: $FRONTEND_PID)"
else
    echo -e "   ${RED}❌ Frontend falhou!${NC}"
    cat /tmp/frontend.log | tail -10
    exit 1
fi

echo

# ============================================================================
# PROMETHEUS
# ============================================================================

if command -v prometheus >/dev/null 2>&1; then
    echo -e "${YELLOW}📊 Iniciando Prometheus...${NC}"

    prometheus --config.file="$PROJECT_ROOT/prometheus.yml" > /tmp/prometheus.log 2>&1 &
    PROM_PID=$!
    echo $PROM_PID > /tmp/barber-prometheus.pid

    sleep 2

    if lsof -ti :9090 >/dev/null 2>&1; then
        echo -e "   ${GREEN}✅ Prometheus iniciado${NC} (PID: $PROM_PID)"
    else
        echo -e "   ${YELLOW}⚠️  Prometheus falhou (continuando assim mesmo)${NC}"
    fi

    echo
fi

# ============================================================================
# SUMMARY
# ============================================================================

echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ BARBER ANALYTICS PRO V2 — PRONTO PARA USAR!${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
echo

echo -e "🌐 ${BLUE}Frontend:${NC}     http://localhost:3000"
echo -e "📡 ${BLUE}API:${NC}          http://localhost:8080"
echo -e "📊 ${BLUE}Prometheus:${NC}   http://localhost:9090"
echo

echo -e "${YELLOW}🔐 LOGIN (Dev Mode):${NC}"
echo -e "   Email: ${BLUE}test@barber.com${NC}"
echo -e "   Senha: ${BLUE}123456${NC} (qualquer coisa em dev mode)"
echo

echo -e "${YELLOW}📋 COMANDOS ÚTEIS:${NC}"
echo -e "   Ver logs backend:   ${BLUE}tail -f /tmp/backend.log${NC}"
echo -e "   Ver logs frontend:  ${BLUE}tail -f /tmp/frontend.log${NC}"
echo -e "   Ver logs Prometheus:${BLUE}tail -f /tmp/prometheus.log${NC}"
echo -e "   Parar tudo:         ${BLUE}./scripts/stop-all.sh${NC}"
echo -e "   Testar API:         ${BLUE}curl http://localhost:8080/api/v1/ping${NC}"
echo

echo -e "${YELLOW}⏸️  Pressione ${RED}CTRL+C${NC}${YELLOW} para parar tudo...${NC}"
echo

# ============================================================================
# TRAP PARA CLEANUP
# ============================================================================

trap_handler() {
    echo
    echo -e "${YELLOW}Parando serviços...${NC}"
    ./scripts/stop-all.sh
    exit 0
}

trap trap_handler INT TERM

# Aguardar indefinidamente
while true; do
    sleep 1
done
