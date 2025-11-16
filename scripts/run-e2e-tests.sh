#!/bin/bash

# ============================================================================
# Script: Executar Testes E2E com Backend + Neon
# Versão: 1.0.0
# Autor: Andrey Viana
# Data: 15/11/2025
# ============================================================================
#
# Este script automatiza a execução de testes E2E do frontend,
# garantindo que o backend esteja rodando e os dados de teste carregados.
#
# Uso:
#   ./scripts/run-e2e-tests.sh [opções]
#
# Opções:
#   --headed        Executar com navegador visível
#   --ui            Abrir interface interativa do Playwright
#   --debug         Executar em modo debug
#   --skip-seed     Pular execução do seed (usar dados existentes)
#   --help          Exibir esta mensagem
#
# ============================================================================

set -e  # Exit on error

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variáveis
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"
BACKEND_PID_FILE="$BACKEND_DIR/.backend-e2e.pid"

# Flags
HEADED=false
UI=false
DEBUG=false
SKIP_SEED=false

# ============================================================================
# Funções Helper
# ============================================================================

print_header() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

print_step() {
    echo -e "${GREEN}▶${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✖${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 não encontrado. Por favor, instale $1."
        exit 1
    fi
}

# ============================================================================
# Parse de Argumentos
# ============================================================================

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --headed)
                HEADED=true
                shift
                ;;
            --ui)
                UI=true
                shift
                ;;
            --debug)
                DEBUG=true
                shift
                ;;
            --skip-seed)
                SKIP_SEED=true
                shift
                ;;
            --help)
                cat << EOF
Uso: ./scripts/run-e2e-tests.sh [opções]

Opções:
  --headed        Executar com navegador visível
  --ui            Abrir interface interativa do Playwright
  --debug         Executar em modo debug
  --skip-seed     Pular execução do seed (usar dados existentes)
  --help          Exibir esta mensagem

Exemplos:
  ./scripts/run-e2e-tests.sh
  ./scripts/run-e2e-tests.sh --headed
  ./scripts/run-e2e-tests.sh --ui --skip-seed

EOF
                exit 0
                ;;
            *)
                print_error "Opção desconhecida: $1"
                echo "Use --help para ver opções disponíveis."
                exit 1
                ;;
        esac
    done
}

# ============================================================================
# Verificar Pré-requisitos
# ============================================================================

check_prerequisites() {
    print_header "Verificando Pré-requisitos"

    print_step "Verificando ferramentas instaladas..."
    check_command "go"
    check_command "pnpm"
    check_command "curl"
    print_success "Todas as ferramentas necessárias estão instaladas"

    print_step "Verificando DATABASE_URL..."
    if [ -z "$DATABASE_URL" ]; then
        if [ -f "$BACKEND_DIR/.env" ]; then
            export $(grep DATABASE_URL "$BACKEND_DIR/.env" | xargs)
        fi

        if [ -z "$DATABASE_URL" ]; then
            print_error "DATABASE_URL não configurada"
            echo ""
            echo "Configure a variável de ambiente:"
            echo "  export DATABASE_URL=\"postgresql://...\""
            echo ""
            echo "Ou adicione no arquivo backend/.env"
            exit 1
        fi
    fi
    print_success "DATABASE_URL configurada"

    print_step "Verificando JWT keys..."
    if [ ! -f "$BACKEND_DIR/keys/private.pem" ] || [ ! -f "$BACKEND_DIR/keys/public.pem" ]; then
        print_warning "JWT keys não encontradas. Gerando..."
        mkdir -p "$BACKEND_DIR/keys"
        openssl genrsa -out "$BACKEND_DIR/keys/private.pem" 2048 2>/dev/null
        openssl rsa -in "$BACKEND_DIR/keys/private.pem" -pubout -out "$BACKEND_DIR/keys/public.pem" 2>/dev/null
        print_success "JWT keys geradas"
    else
        print_success "JWT keys encontradas"
    fi
}

# ============================================================================
# Preparar Banco de Dados
# ============================================================================

prepare_database() {
    print_header "Preparando Banco de Dados"

    print_step "Aplicando migrations..."
    cd "$BACKEND_DIR"
    make migrate-up || print_warning "Migrations já aplicadas ou erro (continuando...)"
    print_success "Migrations aplicadas"

    if [ "$SKIP_SEED" = false ]; then
        print_step "Executando seed de dados de teste..."
        go run scripts/seed_test_data.go
        print_success "Dados de teste carregados"
    else
        print_warning "Pulando seed (--skip-seed especificado)"
    fi
}

# ============================================================================
# Iniciar Backend
# ============================================================================

start_backend() {
    print_header "Iniciando Backend"

    # Verificar se já está rodando
    if [ -f "$BACKEND_PID_FILE" ]; then
        OLD_PID=$(cat "$BACKEND_PID_FILE")
        if kill -0 $OLD_PID 2>/dev/null; then
            print_warning "Backend já está rodando (PID: $OLD_PID)"
            return
        else
            rm "$BACKEND_PID_FILE"
        fi
    fi

    print_step "Iniciando servidor Go..."
    cd "$BACKEND_DIR"

    # Iniciar backend em background
    nohup go run cmd/api/main.go > /tmp/backend-e2e.log 2>&1 &
    BACKEND_PID=$!
    echo $BACKEND_PID > "$BACKEND_PID_FILE"

    print_success "Backend iniciado (PID: $BACKEND_PID)"

    print_step "Aguardando backend estar pronto..."
    local MAX_ATTEMPTS=30
    local ATTEMPT=0

    while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
        if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
            print_success "Backend está pronto!"
            curl http://localhost:8080/health 2>/dev/null | jq '.' || echo ""
            return
        fi

        sleep 1
        ATTEMPT=$((ATTEMPT + 1))

        if [ $((ATTEMPT % 5)) -eq 0 ]; then
            echo -n "."
        fi
    done

    echo ""
    print_error "Backend não respondeu após 30 segundos"
    print_warning "Verificar logs em: /tmp/backend-e2e.log"
    cat /tmp/backend-e2e.log
    cleanup_backend
    exit 1
}

# ============================================================================
# Executar Testes
# ============================================================================

run_tests() {
    print_header "Executando Testes E2E"

    cd "$FRONTEND_DIR"

    # Determinar comando baseado em flags
    local TEST_CMD="pnpm test:e2e"

    if [ "$UI" = true ]; then
        TEST_CMD="pnpm test:e2e:ui"
    elif [ "$DEBUG" = true ]; then
        TEST_CMD="pnpm test:e2e:debug"
    elif [ "$HEADED" = true ]; then
        TEST_CMD="pnpm test:e2e:headed"
    fi

    print_step "Executando: $TEST_CMD"
    echo ""

    set +e  # Não parar em erro (queremos fazer cleanup)
    $TEST_CMD
    TEST_EXIT_CODE=$?
    set -e

    echo ""

    if [ $TEST_EXIT_CODE -eq 0 ]; then
        print_success "Todos os testes passaram! ✨"
    else
        print_error "Alguns testes falharam (exit code: $TEST_EXIT_CODE)"
    fi

    return $TEST_EXIT_CODE
}

# ============================================================================
# Cleanup
# ============================================================================

cleanup_backend() {
    if [ -f "$BACKEND_PID_FILE" ]; then
        BACKEND_PID=$(cat "$BACKEND_PID_FILE")
        if kill -0 $BACKEND_PID 2>/dev/null; then
            print_step "Parando backend (PID: $BACKEND_PID)..."
            kill $BACKEND_PID
            rm "$BACKEND_PID_FILE"
            print_success "Backend parado"
        fi
    fi
}

cleanup() {
    print_header "Cleanup"
    cleanup_backend

    # Remover logs temporários
    if [ -f "/tmp/backend-e2e.log" ]; then
        rm /tmp/backend-e2e.log
    fi
}

# Trap para garantir cleanup mesmo em caso de erro
trap cleanup EXIT INT TERM

# ============================================================================
# Main
# ============================================================================

main() {
    parse_args "$@"

    print_header "🧪 Barber Analytics Pro - E2E Tests"

    check_prerequisites
    prepare_database
    start_backend

    # Executar testes
    run_tests
    TEST_RESULT=$?

    # Cleanup será feito pelo trap
    exit $TEST_RESULT
}

main "$@"
