package valueobject

import (
	"fmt"
	"regexp"
	"time"
)

// MesAno representa um período no formato YYYY-MM
type MesAno struct {
	ano int
	mes int
}

var mesAnoRegex = regexp.MustCompile(`^(\d{4})-(\d{2})$`)

// NewMesAno cria um novo MesAno a partir de uma string "YYYY-MM"
func NewMesAno(str string) (MesAno, error) {
	matches := mesAnoRegex.FindStringSubmatch(str)
	if matches == nil {
		return MesAno{}, fmt.Errorf("formato inválido para mes_ano, esperado YYYY-MM, recebido: %s", str)
	}

	var ano, mes int
	fmt.Sscanf(matches[1], "%d", &ano)
	fmt.Sscanf(matches[2], "%d", &mes)

	if mes < 1 || mes > 12 {
		return MesAno{}, fmt.Errorf("mês inválido: %d (deve estar entre 1-12)", mes)
	}

	if ano < 2000 || ano > 2100 {
		return MesAno{}, fmt.Errorf("ano inválido: %d (deve estar entre 2000-2100)", ano)
	}

	return MesAno{ano: ano, mes: mes}, nil
}

// NewMesAnoFromTime cria um MesAno a partir de time.Time
func NewMesAnoFromTime(t time.Time) MesAno {
	return MesAno{
		ano: t.Year(),
		mes: int(t.Month()),
	}
}

// NewMesAnoFromInts cria um MesAno a partir de ano e mês
func NewMesAnoFromInts(ano, mes int) (MesAno, error) {
	if mes < 1 || mes > 12 {
		return MesAno{}, fmt.Errorf("mês inválido: %d", mes)
	}
	if ano < 2000 || ano > 2100 {
		return MesAno{}, fmt.Errorf("ano inválido: %d", ano)
	}
	return MesAno{ano: ano, mes: mes}, nil
}

// String retorna o formato "YYYY-MM"
func (m MesAno) String() string {
	return fmt.Sprintf("%04d-%02d", m.ano, m.mes)
}

// Ano retorna o ano
func (m MesAno) Ano() int {
	return m.ano
}

// Mes retorna o mês (1-12)
func (m MesAno) Mes() int {
	return m.mes
}

// PrimeiroDia retorna o primeiro dia do mês
func (m MesAno) PrimeiroDia() time.Time {
	return time.Date(m.ano, time.Month(m.mes), 1, 0, 0, 0, 0, time.UTC)
}

// UltimoDia retorna o último dia do mês
func (m MesAno) UltimoDia() time.Time {
	return m.PrimeiroDia().AddDate(0, 1, -1)
}

// Anterior retorna o mês anterior
func (m MesAno) Anterior() MesAno {
	mes := m.mes - 1
	ano := m.ano
	if mes < 1 {
		mes = 12
		ano--
	}
	return MesAno{ano: ano, mes: mes}
}

// Proximo retorna o próximo mês
func (m MesAno) Proximo() MesAno {
	mes := m.mes + 1
	ano := m.ano
	if mes > 12 {
		mes = 1
		ano++
	}
	return MesAno{ano: ano, mes: mes}
}

// Equals verifica igualdade
func (m MesAno) Equals(other MesAno) bool {
	return m.ano == other.ano && m.mes == other.mes
}

// Before verifica se é anterior a outro MesAno
func (m MesAno) Before(other MesAno) bool {
	if m.ano < other.ano {
		return true
	}
	if m.ano == other.ano && m.mes < other.mes {
		return true
	}
	return false
}

// After verifica se é posterior a outro MesAno
func (m MesAno) After(other MesAno) bool {
	if m.ano > other.ano {
		return true
	}
	if m.ano == other.ano && m.mes > other.mes {
		return true
	}
	return false
}

// MesAtual retorna o MesAno atual
func MesAtual() MesAno {
	return NewMesAnoFromTime(time.Now())
}

// MesAnterior retorna o mês anterior ao atual
func MesAnterior() MesAno {
	return MesAtual().Anterior()
}
