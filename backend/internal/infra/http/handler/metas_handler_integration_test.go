package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/metas"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/handler"
	"github.com/andviana23/barber-analytics-backend/internal/infra/repository/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	testDBPool *pgxpool.Pool
	testLogger *zap.Logger
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if cv.validator == nil {
		cv.validator = validator.New()
	}
	return cv.validator.Struct(i)
}

func TestMain(m *testing.M) {
	// Setup
	var err error
	testLogger, _ = zap.NewDevelopment()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		testLogger.Warn("DATABASE_URL não configurada, pulando testes de integração")
		os.Exit(0)
	}

	ctx := context.Background()
	testDBPool, err = pgxpool.New(ctx, databaseURL)
	if err != nil {
		testLogger.Fatal("Erro ao conectar ao banco de testes", zap.Error(err))
	}

	if err := testDBPool.Ping(ctx); err != nil {
		testLogger.Fatal("Erro ao fazer ping no banco", zap.Error(err))
	}

	// Run tests
	code := m.Run()

	// Teardown
	testDBPool.Close()
	os.Exit(code)
}

func setupTestHandler() (*handler.MetasHandler, *echo.Echo) {
	queries := db.New(testDBPool)

	// Repositories
	metaMensalRepo := postgres.NewMetaMensalRepository(queries)
	metaBarbeiroRepo := postgres.NewMetaBarbeiroRepository(queries)
	metaTicketMedioRepo := postgres.NewMetasTicketMedioRepository(queries)

	// Use cases
	setMetaMensalUC := metas.NewSetMetaMensalUseCase(metaMensalRepo, testLogger)
	getMetaMensalUC := metas.NewGetMetaMensalUseCase(metaMensalRepo, testLogger)
	listMetasMensaisUC := metas.NewListMetasMensaisUseCase(metaMensalRepo, testLogger)
	updateMetaMensalUC := metas.NewUpdateMetaMensalUseCase(metaMensalRepo, testLogger)
	deleteMetaMensalUC := metas.NewDeleteMetaMensalUseCase(metaMensalRepo, testLogger)

	setMetaBarbeiroUC := metas.NewSetMetaBarbeiroUseCase(metaBarbeiroRepo, testLogger)
	getMetaBarbeiroUC := metas.NewGetMetaBarbeiroUseCase(metaBarbeiroRepo, testLogger)
	listMetasBarbeiroUC := metas.NewListMetasBarbeiroUseCase(metaBarbeiroRepo, testLogger)
	updateMetaBarbeiroUC := metas.NewUpdateMetaBarbeiroUseCase(metaBarbeiroRepo, testLogger)
	deleteMetaBarbeiroUC := metas.NewDeleteMetaBarbeiroUseCase(metaBarbeiroRepo, testLogger)

	setMetaTicketMedioUC := metas.NewSetMetaTicketUseCase(metaTicketMedioRepo, testLogger)
	getMetaTicketMedioUC := metas.NewGetMetaTicketMedioUseCase(metaTicketMedioRepo, testLogger)
	listMetasTicketMedioUC := metas.NewListMetasTicketMedioUseCase(metaTicketMedioRepo, testLogger)
	updateMetaTicketMedioUC := metas.NewUpdateMetaTicketMedioUseCase(metaTicketMedioRepo, testLogger)
	deleteMetaTicketMedioUC := metas.NewDeleteMetaTicketMedioUseCase(metaTicketMedioRepo, testLogger)

	// Handler
	metasHandler := handler.NewMetasHandler(
		setMetaMensalUC,
		getMetaMensalUC,
		listMetasMensaisUC,
		updateMetaMensalUC,
		deleteMetaMensalUC,
		setMetaBarbeiroUC,
		getMetaBarbeiroUC,
		listMetasBarbeiroUC,
		updateMetaBarbeiroUC,
		deleteMetaBarbeiroUC,
		setMetaTicketMedioUC,
		getMetaTicketMedioUC,
		listMetasTicketMedioUC,
		updateMetaTicketMedioUC,
		deleteMetaTicketMedioUC,
		testLogger,
	)

	e := echo.New()
	e.Validator = &CustomValidator{}
	return metasHandler, e
}

func TestMetasHandler_SetMetaMensal_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de integração")
	}

	metasHandler, e := setupTestHandler()

	// Preparar payload
	payload := map[string]interface{}{
		"mes_ano":          "2025-01",
		"meta_faturamento": "150000",
		"origem":           "MANUAL",
	}
	body, _ := json.Marshal(payload)

	// Criar request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/metas/monthly", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("X-Tenant-ID", "e2e00000-0000-0000-0000-000000000001")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("tenant_id", "e2e00000-0000-0000-0000-000000000001")

	// Executar
	err := metasHandler.SetMetaMensal(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response["id"])
	assert.Equal(t, "2025-01", response["mes_ano"])
	assert.Equal(t, "150000", response["meta_faturamento"])
}

func TestMetasHandler_ListMetasMensais_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de integração")
	}

	metasHandler, e := setupTestHandler()

	// Criar request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/metas/monthly?page=1&page_size=10", nil)
	req.Header.Set("X-Tenant-ID", "e2e00000-0000-0000-0000-000000000001")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("tenant_id", "e2e00000-0000-0000-0000-000000000001")

	// Executar
	err := metasHandler.ListMetasMensais(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotNil(t, response["data"])
	assert.NotNil(t, response["total"])
	assert.NotNil(t, response["page"])
	assert.NotNil(t, response["page_size"])
}
