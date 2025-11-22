package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/google/uuid"
)

// UserPreferences representa as preferências de privacidade do usuário (LGPD)
type UserPreferences struct {
	ID     string
	UserID string

	DataSharingConsent     bool
	MarketingConsent       bool
	AnalyticsConsent       bool
	ThirdPartyConsent      bool
	PersonalizedAdsConsent bool

	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewUserPreferences cria novas preferências de usuário
func NewUserPreferences(userID string) (*UserPreferences, error) {
	if userID == "" {
		return nil, domain.ErrInvalidID
	}

	now := time.Now()
	return &UserPreferences{
		ID:                     uuid.NewString(),
		UserID:                 userID,
		DataSharingConsent:     false,
		MarketingConsent:       false,
		AnalyticsConsent:       false,
		ThirdPartyConsent:      false,
		PersonalizedAdsConsent: false,
		CriadoEm:               now,
		AtualizadoEm:           now,
	}, nil
}

// AtualizarConsentimentos atualiza os consentimentos
func (u *UserPreferences) AtualizarConsentimentos(
	dataSharing, marketing, analytics, thirdParty, personalizedAds bool,
) {
	u.DataSharingConsent = dataSharing
	u.MarketingConsent = marketing
	u.AnalyticsConsent = analytics
	u.ThirdPartyConsent = thirdParty
	u.PersonalizedAdsConsent = personalizedAds
	u.AtualizadoEm = time.Now()
}

// RevogarTodosConsentimentos revoga todos os consentimentos
func (u *UserPreferences) RevogarTodosConsentimentos() {
	u.DataSharingConsent = false
	u.MarketingConsent = false
	u.AnalyticsConsent = false
	u.ThirdPartyConsent = false
	u.PersonalizedAdsConsent = false
	u.AtualizadoEm = time.Now()
}

// Validate valida as regras de negócio
func (u *UserPreferences) Validate() error {
	if u.UserID == "" {
		return domain.ErrInvalidID
	}
	return nil
}
