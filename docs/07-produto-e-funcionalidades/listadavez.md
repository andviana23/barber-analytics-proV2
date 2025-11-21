> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

---
title: 'Lista da Vez - Lógica de Funcionamento'
author: 'Andrey Viana'
version: '2.0.0'
last_updated: '18/11/2025'
---

# 📋 Lista da Vez - Lógica de Funcionamento

## 🎯 Conceito Principal

A **Lista da Vez** é um sistema de **fila rotativa baseada em pontos** que distribui clientes de forma equitativa entre barbeiros de uma unidade.

### Princípio Básico

```
Barbeiro atende cliente → Ganha +1 ponto → Vai para o final da fila
```

A fila é **automaticamente reordenada** a cada atendimento, sempre colocando o barbeiro com **menos pontos** no topo.

---

## 🔄 Fluxo de Funcionamento

### 1. Inicialização da Lista

**Quando**: No primeiro acesso ou quando não existe lista para a unidade.

**O que acontece**:

1. Sistema busca todos os barbeiros ativos da unidade
2. Cria registros na tabela `barbers_turn_list`
3. Todos começam com **0 pontos**
4. Ordenação inicial: por data de cadastro (mais antigo primeiro)

```javascript
// Lógica de inicialização
const activeBarbers = await getActiveBarbersByUnit(unitId);
// Cria registro para cada barbeiro:
// { professional_id, unit_id, points: 0, position: auto }
```

### 2. Ordenação da Fila

**Critérios de ordenação (prioridade decrescente)**:

1. **Menor quantidade de pontos** (`current_points ASC`)
2. **Último atendimento mais antigo** (`last_turn_at ASC NULLS FIRST`)
3. **Nome do barbeiro** (desempate alfabético)

```sql
ORDER BY
  current_points ASC,
  last_turn_at ASC NULLS FIRST,
  professional_name ASC
```

**Exemplo prático**:

| Barbeiro | Pontos | Último Atendimento | Posição      |
| -------- | ------ | ------------------ | ------------ |
| João     | 0      | null               | 1º (próximo) |
| Maria    | 0      | 10:30              | 2º           |
| Pedro    | 1      | 09:15              | 3º           |
| Ana      | 2      | 11:00              | 4º           |

### 3. Registro de Atendimento

**Quando**: Barbeiro atende um cliente.

**Processo**:

1. **Identificar barbeiro**: Sistema pega o ID do profissional que atendeu
2. **Incrementar pontos**: `current_points = current_points + 1`
3. **Atualizar timestamp**: `last_turn_at = NOW()`
4. **Reordenar automaticamente**: Lista se reorganiza baseada nos novos pontos

```javascript
// Função fn_record_barber_turn
UPDATE barbers_turn_list
SET
  current_points = current_points + 1,
  last_turn_at = NOW(),
  updated_at = NOW()
WHERE professional_id = p_professional_id
  AND unit_id = p_unit_id;
```

**Resultado**: Barbeiro que atendeu vai automaticamente para uma posição mais abaixo na fila.

### 4. Busca do Próximo Barbeiro

**Quando**: Cliente chega sem preferência de barbeiro.

**Lógica SQL**:

```sql
SELECT
  professional_id,
  professional_name,
  current_points,
  last_turn_at
FROM barbers_turn_list btl
INNER JOIN professionals p ON p.id = btl.professional_id
WHERE
  btl.unit_id = p_unit_id
  AND btl.is_active = true  -- Apenas barbeiros ativos
  AND p.is_active = true    -- Profissional não pausado
ORDER BY
  btl.current_points ASC,  -- Menor pontuação primeiro
  btl.last_turn_at ASC NULLS FIRST,  -- Mais tempo sem atender
  p.name ASC  -- Desempate alfabético
LIMIT 1;  -- Retorna apenas o próximo da fila
```

**Retorno**: Profissional no topo da fila (menor pontuação).

### 5. Reset Mensal Automático

**Quando**: Todo dia 1º de cada mês às 23:00 (via pg_cron).

**Processo**:

1. **Salvar histórico**: Copia estado atual para `barber_turn_history`
   - Salva: pontos finais, posição final, mês/ano de referência

2. **Zerar pontos**: Reseta todos os `current_points` para 0

3. **Limpar timestamps**: `last_turn_at = NULL`

4. **Reiniciar ciclo**: Todos voltam para mesma posição inicial

```sql
-- Executado automaticamente pelo Cron Job
CREATE OR REPLACE FUNCTION fn_reset_barber_turn_list()
RETURNS void AS $$
BEGIN
  -- 1. Salvar estado atual no histórico
  INSERT INTO barber_turn_history (
    professional_id, unit_id, month_year,
    total_turns, final_points
  )
  SELECT
    professional_id,
    unit_id,
    TO_CHAR(CURRENT_DATE - INTERVAL '1 month', 'YYYY-MM'),
    current_points,
    current_points
  FROM barbers_turn_list
  WHERE is_active = true;

  -- 2. Resetar pontos de todos os barbeiros
  UPDATE barbers_turn_list
  SET
    current_points = 0,
    last_turn_at = NULL,
    updated_at = now();
END;
$$ LANGUAGE plpgsql;
```

**Resultado**: No início de cada mês, todos barbeiros começam do zero novamente.

---

## 📊 Estrutura de Dados

### Tabela Principal: `barbers_turn_list`

Armazena o **estado atual** da fila.

```sql
CREATE TABLE barbers_turn_list (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  unit_id UUID NOT NULL REFERENCES units(id),
  professional_id UUID NOT NULL REFERENCES professionals(id),
  current_points INTEGER DEFAULT 0,     -- Pontos acumulados no mês
  last_turn_at TIMESTAMPTZ,             -- Último atendimento
  is_active BOOLEAN DEFAULT true,       -- Se está participando da fila
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),

  UNIQUE(professional_id, unit_id)      -- 1 barbeiro por unidade
);

-- Índices para performance
CREATE INDEX idx_barbers_turn_unit ON barbers_turn_list(unit_id);
CREATE INDEX idx_barbers_turn_points ON barbers_turn_list(current_points);
CREATE INDEX idx_barbers_turn_active ON barbers_turn_list(is_active);
```

**Campos**:

- `professional_id`: UUID do barbeiro
- `unit_id`: UUID da unidade (barbearia)
- `current_points`: Pontos acumulados (0 = início da fila)
- `last_turn_at`: Timestamp do último atendimento
- `is_active`: Se o barbeiro está na fila (true) ou pausado (false)

### Tabela de Histórico: `barber_turn_history`

Armazena **snapshot mensal** para relatórios.

```sql
CREATE TABLE barber_turn_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  professional_id UUID NOT NULL REFERENCES professionals(id),
  unit_id UUID NOT NULL REFERENCES units(id),
  month_year VARCHAR(7) NOT NULL,       -- 'YYYY-MM'
  total_turns INTEGER DEFAULT 0,        -- Total de atendimentos
  final_points INTEGER DEFAULT 0,       -- Pontos no fim do mês
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_turn_history_month ON barber_turn_history(month_year);
CREATE INDEX idx_turn_history_unit ON barber_turn_history(unit_id);
```

---

## ⚙️ Regras de Negócio

### 1. Participação na Fila

**Condições para aparecer na fila**:

- ✅ Profissional com `role = 'barbeiro'`
- ✅ `professionals.is_active = true` (profissional ativo no sistema)
- ✅ `barbers_turn_list.is_active = true` (participando da fila)

**Barbeiro pode ser pausado**: Gerente alterna `is_active` para `false`.

### 2. Cálculo de Pontos

- **+1 ponto** por atendimento registrado
- **Não há pontos negativos** (mínimo sempre é 0)
- **Sem limite máximo** de pontos no mês
- **Reset automático** no dia 1º de cada mês

### 3. Justiça na Distribuição

**Cenário**: 3 barbeiros com pontos diferentes

```
João:  0 pontos → Próximo atendimento
Maria: 0 pontos → Próximo após João
Pedro: 1 ponto  → Próximo após Maria
```

**Após João atender**:

```
Maria: 0 pontos → Agora é o próximo
Pedro: 1 ponto  → Depois de Maria
João:  1 ponto  → Empatou com Pedro, mas atendeu por último
```

**Desempate**: Se dois barbeiros têm mesma pontuação, **atende primeiro quem está há mais tempo sem atender**.

### 4. Casos Especiais

#### Cliente com Preferência

- Cliente pode escolher barbeiro específico
- Barbeiro escolhido ganha +1 ponto normalmente
- Fila se reorganiza após o atendimento

#### Barbeiro em Pausa

- `is_active = false` → **não aparece** na lista de próximos
- Pontos acumulados são **mantidos**
- Quando retorna (`is_active = true`), volta com mesmos pontos

#### Novo Barbeiro no Meio do Mês

- Entra com **0 pontos**
- Vai para o topo da fila automaticamente
- Tem vantagem inicial (começou zerado quando outros já tinham pontos)

---

## 🔐 Segurança e Permissões

### Row Level Security (RLS)

**Políticas aplicadas**:

```sql
-- SELECT: Usuário vê apenas sua unidade
CREATE POLICY "view_own_unit_turn_list"
ON barbers_turn_list FOR SELECT
USING (
  unit_id IN (
    SELECT unit_id FROM professionals
    WHERE user_id = auth.uid()
  )
);

-- INSERT/UPDATE: Apenas admin/gerente
CREATE POLICY "manage_turn_list"
ON barbers_turn_list FOR ALL
USING (
  EXISTS (
    SELECT 1 FROM professionals
    WHERE user_id = auth.uid()
    AND role IN ('admin', 'gerente', 'manager')
  )
);
```

### Controle de Acesso

| Ação                   | Barbeiro | Gerente | Admin |
| ---------------------- | -------- | ------- | ----- |
| Ver lista              | ✅       | ✅      | ✅    |
| Registrar atendimento  | ❌       | ✅      | ✅    |
| Pausar/ativar barbeiro | ❌       | ✅      | ✅    |
| Executar reset manual  | ❌       | ✅      | ✅    |
| Ver histórico          | ✅       | ✅      | ✅    |

---

## 📈 Casos de Uso

### Caso 1: Cliente Chega na Barbearia

**Fluxo**:

1. Recepcionista abre a página "Lista da Vez"
2. Sistema mostra **próximo barbeiro** em destaque (topo da fila)
3. Recepcionista clica em "Registrar Atendimento"
4. Sistema:
   - Adiciona +1 ponto ao barbeiro
   - Atualiza `last_turn_at`
   - Reordena fila automaticamente
5. **Novo próximo barbeiro** aparece no topo

### Caso 2: Barbeiro Sai para Almoço

**Fluxo**:

1. Gerente clica no **switch de ativação** do barbeiro
2. Sistema atualiza `is_active = false`
3. Barbeiro **desaparece** da lista de próximos
4. Pontos são **mantidos**
5. Quando retorna, gerente ativa novamente
6. Barbeiro volta com **mesmos pontos** acumulados

### Caso 3: Final do Mês

**Fluxo automático (pg_cron)**:

1. **Dia 1º às 23:00**: Cron job dispara `fn_reset_barber_turn_list()`
2. Estado atual é **salvo** em `barber_turn_history`:
   - Pontos finais de cada barbeiro
   - Posição final no ranking
   - Mês/ano de referência
3. Todos os pontos são **zerados**
4. Timestamps `last_turn_at` são **limpos**
5. **Novo ciclo começa** com todos em pé de igualdade

### Caso 4: Consulta de Relatório Mensal

**Fluxo**:

1. Gerente acessa aba "Histórico"
2. Seleciona mês/ano desejado
3. Sistema busca em `barber_turn_history`
4. Mostra ranking final:
   - Quantos pontos cada barbeiro fez
   - Posição final no ranking
   - Comparação com meses anteriores

---

## 🧮 Algoritmo de Reordenação

### Pseudocódigo

```javascript
function reorderTurnList(unitId) {
  // 1. Buscar todos os barbeiros ativos
  const barbers = SELECT * FROM barbers_turn_list
                  WHERE unit_id = unitId
                  AND is_active = true;

  // 2. Ordenar por critérios
  const sorted = barbers.sort((a, b) => {
    // Critério 1: Menor pontuação
    if (a.current_points !== b.current_points) {
      return a.current_points - b.current_points;
    }

    // Critério 2: Último atendimento mais antigo
    if (a.last_turn_at !== b.last_turn_at) {
      if (a.last_turn_at === null) return -1;  // null = nunca atendeu = prioridade
      if (b.last_turn_at === null) return 1;
      return a.last_turn_at - b.last_turn_at;
    }

    // Critério 3: Nome alfabético
    return a.professional_name.localeCompare(b.professional_name);
  });

  // 3. Retornar lista ordenada
  return sorted;
}
```

### Complexidade

- **Tempo**: O(n log n) - Ordenação de array
- **Espaço**: O(n) - Lista temporária
- **Frequência**: A cada atendimento registrado

---

## 📊 Métricas e KPIs

### Estatísticas em Tempo Real

```javascript
const stats = {
  totalBarbers: 8, // Barbeiros na unidade
  totalPoints: 15, // Soma de todos os pontos
  averagePoints: 1.88, // Média de pontos por barbeiro
  barbersWithPoints: 5, // Quantos já atenderam no mês
  lastUpdated: '2025-11-18T14:30:00', // Última atualização
};
```

### Relatório Mensal

```javascript
const monthlyReport = {
  month: 11,
  year: 2025,
  unitName: 'Unidade Centro',
  barbers: [
    { name: 'João Silva', totalPoints: 45, finalPosition: 1 },
    { name: 'Maria Santos', totalPoints: 42, finalPosition: 2 },
    { name: 'Pedro Costa', totalPoints: 38, finalPosition: 3 },
  ],
  totalPoints: 125,
  averagePoints: 41.67,
};
```

---

## 🎯 Vantagens do Sistema

### 1. Equidade

- Distribuição justa de clientes
- Todos barbeiros têm mesmas oportunidades
- Sem favoritismo ou disputa

### 2. Transparência

- Fila visível para todos
- Histórico completo armazenado
- Critérios claros e objetivos

### 3. Automação

- Reordenação automática após cada atendimento
- Reset mensal sem intervenção manual
- Sem necessidade de controle manual

### 4. Flexibilidade

- Barbeiros podem ser pausados/ativados
- Cliente pode escolher barbeiro preferido
- Gerente pode executar reset manual se necessário

### 5. Rastreabilidade

- Histórico mensal completo
- Auditoria de todos os atendimentos
- Relatórios e estatísticas

---

## 🔄 Integração com Outros Módulos

### Receitas (Revenues)

**Conexão**: Quando receita é registrada com `professional_id`:

- Sistema pode automaticamente adicionar +1 ponto na Lista da Vez
- Sincronização: receita → atendimento → ponto

### Dashboard

**Widgets**:

- Card "Próximo Barbeiro da Vez"
- Gráfico de distribuição de pontos
- Ranking mensal em tempo real

### Notificações

**Alertas automáticos**:

- Reset mensal executado com sucesso
- Barbeiro atingiu X pontos no mês
- Discrepância na distribuição (alguém muito abaixo ou acima da média)

---

## 🚀 Fluxo Técnico Completo

### 1. Frontend → Backend

```javascript
// 1. Hook (useListaDaVez.js)
const { mutate: recordTurn } = useMutation({
  mutationFn: (professionalId) =>
    listaDaVezService.recordTurn(professionalId, unitId, user)
});

// 2. Service (listaDaVezService.js)
async recordTurn(professionalId, unitId, user) {
  // Validar permissão
  if (!this.canManage(user)) {
    return { error: 'Apenas gerentes' };
  }

  // Chamar repository
  return await listaDaVezRepository.recordTurn(professionalId, unitId);
}

// 3. Repository (listaDaVezRepository.js)
async recordTurn(professionalId, unitId) {
  const { error } = await supabase.rpc('fn_record_barber_turn', {
    p_professional_id: professionalId,
    p_unit_id: unitId
  });
  return { error };
}
```

### 2. Banco de Dados

```sql
-- 4. Função SQL (fn_record_barber_turn)
CREATE OR REPLACE FUNCTION fn_record_barber_turn(
  p_professional_id UUID,
  p_unit_id UUID
)
RETURNS void AS $$
BEGIN
  UPDATE barbers_turn_list
  SET
    current_points = current_points + 1,
    last_turn_at = now(),
    updated_at = now()
  WHERE
    professional_id = p_professional_id
    AND unit_id = p_unit_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Barbeiro não encontrado na lista da vez';
  END IF;
END;
$$ LANGUAGE plpgsql;
```

### 3. Retorno e Atualização

```javascript
// 5. Hook atualiza cache
onSuccess: () => {
  queryClient.invalidateQueries(['lista-da-vez', unitId]);
  queryClient.invalidateQueries(['next-barber', unitId]);
  toast.success('Atendimento registrado!');
};

// 6. UI atualiza automaticamente
// TanStack Query refetch as queries
// Lista se reordena com novos pontos
```

---

## 📝 Exemplo Prático Completo

### Estado Inicial (09:00)

```
┌──────────┬────────┬──────────────────┐
│ Barbeiro │ Pontos │ Último Atendim.  │
├──────────┼────────┼──────────────────┤
│ João     │ 0      │ null             │ ← PRÓXIMO
│ Maria    │ 0      │ null             │
│ Pedro    │ 0      │ null             │
└──────────┴────────┴──────────────────┘
```

### Após 1º Cliente (09:15)

João atende → Ganha +1 ponto

```
┌──────────┬────────┬──────────────────┐
│ Barbeiro │ Pontos │ Último Atendim.  │
├──────────┼────────┼──────────────────┤
│ Maria    │ 0      │ null             │ ← PRÓXIMO
│ Pedro    │ 0      │ null             │
│ João     │ 1      │ 09:15            │
└──────────┴────────┴──────────────────┘
```

### Após 2º Cliente (09:30)

Maria atende → Ganha +1 ponto

```
┌──────────┬────────┬──────────────────┐
│ Barbeiro │ Pontos │ Último Atendim.  │
├──────────┼────────┼──────────────────┤
│ Pedro    │ 0      │ null             │ ← PRÓXIMO
│ João     │ 1      │ 09:15            │
│ Maria    │ 1      │ 09:30            │
└──────────┴────────┴──────────────────┘
```

### Após 3º Cliente (09:45)

Pedro atende → Todos empatados

```
┌──────────┬────────┬──────────────────┐
│ Barbeiro │ Pontos │ Último Atendim.  │
├──────────┼────────┼──────────────────┤
│ João     │ 1      │ 09:15            │ ← PRÓXIMO (mais antigo)
│ Maria    │ 1      │ 09:30            │
│ Pedro    │ 1      │ 09:45            │
└──────────┴────────┴──────────────────┘
```

### Final do Dia (18:00)

```
┌──────────┬────────┬──────────────────┐
│ Barbeiro │ Pontos │ Último Atendim.  │
├──────────┼────────┼──────────────────┤
│ Pedro    │ 8      │ 17:45            │ ← Próximo amanhã
│ Maria    │ 9      │ 17:30            │
│ João     │ 9      │ 17:50            │
└──────────┴────────┴──────────────────┘
```

### Início do Próximo Mês (01/12 - 00:00)

Reset automático executado

```
┌──────────┬────────┬──────────────────┐
│ Barbeiro │ Pontos │ Último Atendim.  │
├──────────┼────────┼──────────────────┤
│ João     │ 0      │ null             │ ← PRÓXIMO
│ Maria    │ 0      │ null             │
│ Pedro    │ 0      │ null             │
└──────────┴────────┴──────────────────┘

Histórico salvo em barber_turn_history:
- João: 45 pontos (novembro)
- Maria: 42 pontos (novembro)
- Pedro: 38 pontos (novembro)
```

---

## 🎓 Resumo da Lógica

### Fluxo Simplificado

1. **Cliente chega** → Sistema busca próximo barbeiro (menor pontuação)
2. **Barbeiro atende** → Sistema adiciona +1 ponto
3. **Fila reordena** → Barbeiro desce na lista automaticamente
4. **Próximo cliente** → Novo barbeiro no topo atende
5. **Final do mês** → Reset automático, todos voltam a 0

### Regras Essenciais

- 🎯 **Menor pontuação = Prioridade** (0 pontos = topo da fila)
- ⏰ **Último atendimento** = Critério de desempate
- 🔄 **Automático** = Sem intervenção manual necessária
- 📊 **Histórico** = Todos os meses salvos para relatórios
- 🔐 **Seguro** = RLS impede manipulação indevida

---

**Última atualização**: 18 de novembro de 2025
**Versão**: 2.0.0
**Autor**: Andrey Viana
