#!/bin/bash

# ============================================================================
# 🛑 Barber Analytics Pro V2 — Stop All Services
# Para Backend + Frontend + Prometheus
# ============================================================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║${NC}  🛑 Barber Analytics Pro V2 — Parando Serviços"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo

# ============================================================================
# KILL SERVICES
# ============================================================================

echo -e "${YELLOW}Parando serviços...${NC}"
echo

# Kill Backend
if [ -f /tmp/barber-backend.pid ]; then
    PID=$(cat /tmp/barber-backend.pid)
    if kill $PID 2>/dev/null; then
        echo -e "   ${GREEN}✅${NC} Backend parado (PID: $PID)"
    fi
    rm -f /tmp/barber-backend.pid
fi

# Kill Frontend
if [ -f /tmp/barber-frontend.pid ]; then
    PID=$(cat /tmp/barber-frontend.pid)
    if kill $PID 2>/dev/null; then
        echo -e "   ${GREEN}✅${NC} Frontend parado (PID: $PID)"
    fi
    rm -f /tmp/barber-frontend.pid
fi

# Kill Prometheus
if [ -f /tmp/barber-prometheus.pid ]; then
    PID=$(cat /tmp/barber-prometheus.pid)
    if kill $PID 2>/dev/null; then
        echo -e "   ${GREEN}✅${NC} Prometheus parado (PID: $PID)"
    fi
    rm -f /tmp/barber-prometheus.pid
fi

# Kill by port (fallback)
for port in 8080 3000 9090; do
    if lsof -ti :$port 2>/dev/null | xargs kill -9 2>/dev/null; then
        echo -e "   ${GREEN}✅${NC} Porta $port liberada"
    fi
done

echo
echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ Tudo parado!${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
