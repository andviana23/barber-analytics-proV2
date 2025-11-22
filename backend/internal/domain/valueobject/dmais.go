package valueobject

import (
	"fmt"
	"time"
)

// DMais representa dias de compensação bancária (D+0, D+1, D+30, etc)
type DMais struct {
	dias int
}

// NewDMais cria um novo DMais
func NewDMais(dias int) (DMais, error) {
	if dias < 0 {
		return DMais{}, fmt.Errorf("d+ não pode ser negativo: %d", dias)
	}
	if dias > 365 {
		return DMais{}, fmt.Errorf("d+ não pode exceder 365 dias: %d", dias)
	}
	return DMais{dias: dias}, nil
}

// NewDMaisUnsafe cria um DMais sem validação
func NewDMaisUnsafe(dias int) DMais {
	return DMais{dias: dias}
}

// Dias retorna o número de dias
func (d DMais) Dias() int {
	return d.dias
}

// String retorna a representação como string (ex: "D+1")
func (d DMais) String() string {
	return fmt.Sprintf("D+%d", d.dias)
}

// CalcularDataCompensacao calcula a data de compensação a partir de uma data base
func (d DMais) CalcularDataCompensacao(dataBase time.Time) time.Time {
	return dataBase.AddDate(0, 0, d.dias)
}

// Equals verifica igualdade
func (d DMais) Equals(other DMais) bool {
	return d.dias == other.dias
}

// GreaterThan verifica se tem mais dias
func (d DMais) GreaterThan(other DMais) bool {
	return d.dias > other.dias
}

// LessThan verifica se tem menos dias
func (d DMais) LessThan(other DMais) bool {
	return d.dias < other.dias
}

// Constantes comuns
var (
	DPlus0  = DMais{dias: 0}  // Imediato (PIX/Dinheiro)
	DPlus1  = DMais{dias: 1}  // Débito/Transferência
	DPlus30 = DMais{dias: 30} // Crédito
)
