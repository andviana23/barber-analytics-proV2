package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/pricing"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/handler"
	"github.com/andviana23/barber-analytics-backend/internal/infra/repository/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPricingHandler() (*handler.PricingHandler, *echo.Echo) {
	queries := db.New(testDBPool)

	// Repositories
	precificacaoConfigRepo := postgres.NewPrecificacaoConfigRepository(queries)
	precificacaoSimulacaoRepo := postgres.NewPrecificacaoSimulacaoRepository(queries)

	// Use cases
	saveConfigUC := pricing.NewSaveConfigPrecificacaoUseCase(precificacaoConfigRepo, testLogger)
	getConfigUC := pricing.NewGetPrecificacaoConfigUseCase(precificacaoConfigRepo, testLogger)
	updateConfigUC := pricing.NewUpdatePrecificacaoConfigUseCase(precificacaoConfigRepo, testLogger)
	deleteConfigUC := pricing.NewDeletePrecificacaoConfigUseCase(precificacaoConfigRepo, testLogger)
	simularPrecoUC := pricing.NewSimularPrecoUseCase(precificacaoConfigRepo, precificacaoSimulacaoRepo, testLogger)
	saveSimulacaoUC := pricing.NewSaveSimulacaoUseCase(precificacaoSimulacaoRepo, testLogger)
	getSimulacaoUC := pricing.NewGetSimulacaoUseCase(precificacaoSimulacaoRepo, testLogger)
	listSimulacoesUC := pricing.NewListSimulacoesUseCase(precificacaoSimulacaoRepo, testLogger)
	deleteSimulacaoUC := pricing.NewDeleteSimulacaoUseCase(precificacaoSimulacaoRepo, testLogger)

	// Handler
	pricingHandler := handler.NewPricingHandler(
		saveConfigUC,
		getConfigUC,
		updateConfigUC,
		deleteConfigUC,
		simularPrecoUC,
		saveSimulacaoUC,
		getSimulacaoUC,
		listSimulacoesUC,
		deleteSimulacaoUC,
		testLogger,
	)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	return pricingHandler, e
}

func TestPricingHandler_SaveConfig_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de integração")
	}

	pricingHandler, e := setupPricingHandler()

	payload := map[string]interface{}{
		"margem_desejada":    "30.00",
		"markup_alvo":        "1.43",
		"imposto_percentual": "8.50",
		"comissao_default":   "40.00",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pricing/config", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("X-Tenant-ID", "e2e00000-0000-0000-0000-000000000001")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("tenant_id", "e2e00000-0000-0000-0000-000000000001")

	err := pricingHandler.SaveConfig(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response["id"])
	assert.Equal(t, "30.00", response["margem_desejada"])
}

func TestPricingHandler_SimularPreco_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de integração")
	}

	pricingHandler, e := setupPricingHandler()

	payload := map[string]interface{}{
		"item_id":           "00000000-0000-0000-0000-000000000003",
		"tipo_item":         "SERVICO",
		"custo_materiais":   "5000",
		"custo_mao_de_obra": "10000",
		"preco_atual":       "30000",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pricing/simulate", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("X-Tenant-ID", "e2e00000-0000-0000-0000-000000000001")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("tenant_id", "e2e00000-0000-0000-0000-000000000001")

	err := pricingHandler.SimularPreco(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response["preco_sugerido"])
	assert.NotEmpty(t, response["margem_calculada"])
}
