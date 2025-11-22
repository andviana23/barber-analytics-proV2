package valueobject

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Percentage representa uma porcentagem (0-100)
type Percentage struct {
	value decimal.Decimal
}

// NewPercentage cria uma nova porcentagem
// Valida que está entre 0 e 100
func NewPercentage(valor decimal.Decimal) (Percentage, error) {
	if valor.LessThan(decimal.Zero) || valor.GreaterThan(decimal.NewFromInt(100)) {
		return Percentage{}, fmt.Errorf("porcentagem deve estar entre 0 e 100, recebido: %s", valor.String())
	}
	return Percentage{value: valor}, nil
}

// NewPercentageFromFloat cria uma porcentagem a partir de float64
func NewPercentageFromFloat(valor float64) (Percentage, error) {
	return NewPercentage(decimal.NewFromFloat(valor))
}

// NewPercentageUnsafe cria uma porcentagem sem validação (use apenas quando já validado)
func NewPercentageUnsafe(valor decimal.Decimal) Percentage {
	return Percentage{value: valor}
}

// Value retorna o valor como decimal
func (p Percentage) Value() decimal.Decimal {
	return p.value
}

// Float retorna o valor como float64
func (p Percentage) Float() float64 {
	f, _ := p.value.Float64()
	return f
}

// String formata como string (ex: "25.50%")
func (p Percentage) String() string {
	return fmt.Sprintf("%s%%", p.value.StringFixed(2))
}

// AsDecimal retorna a porcentagem como decimal (ex: 25% = 0.25)
func (p Percentage) AsDecimal() decimal.Decimal {
	return p.value.Div(decimal.NewFromInt(100))
}

// Equals verifica igualdade
func (p Percentage) Equals(other Percentage) bool {
	return p.value.Equal(other.value)
}

// GreaterThan verifica se é maior
func (p Percentage) GreaterThan(other Percentage) bool {
	return p.value.GreaterThan(other.value)
}

// LessThan verifica se é menor
func (p Percentage) LessThan(other Percentage) bool {
	return p.value.LessThan(other.value)
}

// IsZero verifica se é zero
func (p Percentage) IsZero() bool {
	return p.value.IsZero()
}

// ZeroPercent retorna 0%
func ZeroPercent() Percentage {
	return Percentage{value: decimal.Zero}
}

// HundredPercent retorna 100%
func HundredPercent() Percentage {
	return Percentage{value: decimal.NewFromInt(100)}
}
