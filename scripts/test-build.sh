#!/bin/bash

# Script de teste de build completo do sistema
# Verifica TypeScript, build de produção e estrutura de arquivos

set -e  # Exit on error

echo "🧪 Teste Completo de Build - Barber Analytics Pro v2.0"
echo "======================================================"
echo ""

# Cores
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Funções auxiliares
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_step() {
    echo ""
    echo "📍 $1"
    echo "---------------------------------------------------"
}

# Navegar para frontend
cd "$(dirname "$0")/../frontend"

# 1. Verificar TypeScript
print_step "1. Verificando TypeScript"
if pnpm tsc --noEmit 2>&1 | grep -q "error TS"; then
    print_error "Erros TypeScript encontrados"
    pnpm tsc --noEmit 2>&1 | grep "error TS" | head -10
    exit 1
else
    print_success "TypeScript: 0 erros"
fi

# 2. Verificar estrutura de componentes
print_step "2. Verificando estrutura de componentes"

COMPONENTS=(
    "app/components/ui/FormInputField.tsx"
    "app/components/ui/CheckboxField.tsx"
    "app/components/ui/DatePickerField.tsx"
    "app/components/ui/TimePickerField.tsx"
    "app/components/ui/SelectField.tsx"
    "app/components/ui/Modal.tsx"
    "app/components/ui/Button.tsx"
    "app/components/ui/DataTable.tsx"
)

for component in "${COMPONENTS[@]}"; do
    if [ -f "$component" ]; then
        print_success "Componente existe: $component"
    else
        print_error "Componente faltando: $component"
        exit 1
    fi
done

# 3. Verificar "use client" directives
print_step "3. Verificando 'use client' directives"

CLIENT_COMPONENTS=(
    "app/components/ui/CheckboxField.tsx"
    "app/components/ui/DatePickerField.tsx"
    "app/components/ui/TimePickerField.tsx"
    "app/components/providers/AppThemeProvider.tsx"
    "app/providers.tsx"
)

for component in "${CLIENT_COMPONENTS[@]}"; do
    if grep -q '"use client"' "$component"; then
        print_success "'use client' presente em: $component"
    else
        print_warning "'use client' ausente em: $component"
    fi
done

# 4. Build de produção
print_step "4. Build de produção (Next.js)"
if pnpm run build > /tmp/build-output.log 2>&1; then
    print_success "Build de produção completado"

    # Verificar rotas geradas
    ROUTES_COUNT=$(grep -c "○ /" /tmp/build-output.log || echo "0")
    print_success "Rotas estáticas geradas: $ROUTES_COUNT"

    # Verificar se build gerou arquivos
    if [ -d ".next" ]; then
        print_success "Diretório .next criado"

        BUILD_SIZE=$(du -sh .next | cut -f1)
        print_success "Tamanho do build: $BUILD_SIZE"
    else
        print_error "Diretório .next não encontrado"
        exit 1
    fi
else
    print_error "Build de produção falhou"
    cat /tmp/build-output.log | tail -30
    exit 1
fi

# 5. Verificar arquivos críticos
print_step "5. Verificando arquivos críticos"

CRITICAL_FILES=(
    ".next/build-manifest.json"
    ".next/app-path-routes-manifest.json"
    ".next/export-marker.json"
)

for file in "${CRITICAL_FILES[@]}"; do
    if [ -f "$file" ]; then
        print_success "Arquivo gerado: $file"
    else
        print_error "Arquivo faltando: $file"
        exit 1
    fi
done

# 6. Resumo final
print_step "6. Resumo Final"
echo ""
print_success "✅ TypeScript: PASSOU"
print_success "✅ Componentes: TODOS PRESENTES"
print_success "✅ Build Produção: SUCESSO"
print_success "✅ Arquivos Gerados: OK"
echo ""
echo "======================================================"
echo -e "${GREEN}🎉 TODOS OS TESTES PASSARAM!${NC}"
echo "======================================================"
echo ""
echo "Sistema pronto para deploy/desenvolvimento"
echo ""

# Cleanup
rm -f /tmp/build-output.log
