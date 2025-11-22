// Package postgres contém funções auxiliares para conversão entre tipos.
package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
)

// uuidToPgtype converte uuid.UUID para pgtype.UUID.
func uuidToPgtype(id uuid.UUID) (pgtype.UUID, error) {
	var pgUUID pgtype.UUID
	err := pgUUID.Scan(id.String())
	return pgUUID, err
}

// uuidStringToPgtype converte string UUID para pgtype.UUID.
func uuidStringToPgtype(id string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return uuidToPgtype(parsed)
}

// pgUUIDToUUID converte pgtype.UUID para uuid.UUID.
func pgUUIDToUUID(pgUUID pgtype.UUID) uuid.UUID {
	if !pgUUID.Valid {
		return uuid.Nil
	}
	id, _ := uuid.FromBytes(pgUUID.Bytes[:])
	return id
}

// pgUUIDToString converte pgtype.UUID para string.
func pgUUIDToString(pgUUID pgtype.UUID) string {
	if !pgUUID.Valid {
		return ""
	}
	id, _ := uuid.FromBytes(pgUUID.Bytes[:])
	return id.String()
}

// decimalToNumeric converte decimal.Decimal para pgtype.Numeric.
func decimalToNumeric(d decimal.Decimal) pgtype.Numeric {
	var num pgtype.Numeric
	_ = num.Scan(d.String())
	return num
}

// numericToDecimal converte pgtype.Numeric para decimal.Decimal.
func numericToDecimal(num pgtype.Numeric) decimal.Decimal {
	if !num.Valid || num.Int == nil {
		return decimal.Zero
	}

	return decimal.NewFromBigInt(num.Int, num.Exp)
}

// timestampToTimestamptz converte time.Time para pgtype.Timestamptz.
func timestampToTimestamptz(t time.Time) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	_ = ts.Scan(t)
	return ts
}

// timestamptzToTime converte pgtype.Timestamptz para time.Time.
func timestamptzToTime(ts pgtype.Timestamptz) time.Time {
	if !ts.Valid {
		return time.Time{}
	}
	return ts.Time
}

// dateToDate converte time.Time para pgtype.Date.
func dateToDate(t time.Time) pgtype.Date {
	var d pgtype.Date
	_ = d.Scan(t)
	return d
}

// dateToTime converte pgtype.Date para time.Time.
func dateToTime(d pgtype.Date) time.Time {
	if !d.Valid {
		return time.Time{}
	}
	return d.Time
}

// stringPtr retorna um ponteiro para string (auxiliar para campos opcionais).
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// boolPtr retorna um ponteiro para bool (auxiliar para campos opcionais).
func boolPtr(b bool) *bool {
	return &b
}

// int32Ptr retorna um ponteiro para int32 (auxiliar para campos opcionais).
func int32Ptr(i int32) *int32 {
	return &i
}

// moneyToNumeric converte Money para pgtype.Numeric.
func moneyToNumeric(m valueobject.Money) pgtype.Numeric {
	return decimalToNumeric(m.Value())
}

// numericToMoney converte pgtype.Numeric para Money.
func numericToMoney(num pgtype.Numeric) valueobject.Money {
	return valueobject.NewMoneyFromDecimal(numericToDecimal(num))
}

// percentageToNumeric converte Percentage para pgtype.Numeric.
func percentageToNumeric(p valueobject.Percentage) pgtype.Numeric {
	return decimalToNumeric(p.Value())
}

// numericToPercentage converte pgtype.Numeric para Percentage.
func numericToPercentage(num pgtype.Numeric) (valueobject.Percentage, error) {
	return valueobject.NewPercentage(numericToDecimal(num))
}

// dateToTimePtr converte pgtype.Date para *time.Time (opcional).
func dateToTimePtr(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	t := d.Time
	return &t
}

// timePtrToDate converte *time.Time para pgtype.Date.
func timePtrToDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return dateToDate(*t)
}

// uuidPtrToPgtype converte *string UUID para pgtype.UUID.
func uuidPtrToPgtype(s *string) (pgtype.UUID, error) {
	if s == nil || *s == "" {
		return pgtype.UUID{Valid: false}, nil
	}
	return uuidStringToPgtype(*s)
}

// pgtypeToUuidPtr converte pgtype.UUID para *string.
func pgtypeToUuidPtr(pgUUID pgtype.UUID) *string {
	if !pgUUID.Valid {
		return nil
	}
	id, _ := uuid.FromBytes(pgUUID.Bytes[:])
	s := id.String()
	return &s
}

// stringPtrToPgText converte *string para pgtype.Text.
func stringPtrToPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// pgTextToStringPtr converte pgtype.Text para *string.
func pgTextToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

// int32ToDMais converte int32 para DMais.
func int32ToDMais(dias int32) valueobject.DMais {
	return valueobject.NewDMaisUnsafe(int(dias))
}

// dmaisToInt32 converte DMais para int32.
func dmaisToInt32(d valueobject.DMais) int32 {
	return int32(d.Dias())
}

// decimalToMoney converte decimal.Decimal para Money.
func decimalToMoney(d decimal.Decimal) valueobject.Money {
	return valueobject.NewMoneyFromDecimal(d)
}

// moneyToDecimal converte Money para decimal.Decimal.
func moneyToDecimal(m valueobject.Money) decimal.Decimal {
	return m.Value()
}

// percentageToDecimal converte Percentage para decimal.Decimal.
func percentageToDecimal(p valueobject.Percentage) decimal.Decimal {
	return p.Value()
}

// decimalToPercentage converte decimal.Decimal para Percentage.
func decimalToPercentage(d decimal.Decimal) (valueobject.Percentage, error) {
	return valueobject.NewPercentage(d)
}
