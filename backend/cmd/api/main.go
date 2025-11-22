package main

import (
	"context"
	"log"
	"os"

	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/financial"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/metas"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/pricing"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/handler"
	"github.com/andviana23/barber-analytics-backend/internal/infra/repository/postgres"
	"github.com/andviana23/barber-analytics-backend/internal/infra/scheduler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Erro ao inicializar logger: %v", err)
	}
	defer logger.Sync()

	// Load environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logger.Fatal("DATABASE_URL não configurada")
	}

	// Initialize database connection
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		logger.Fatal("Erro ao conectar ao banco de dados", zap.Error(err))
	}
	defer dbPool.Close()

	// Test connection
	if err := dbPool.Ping(ctx); err != nil {
		logger.Fatal("Erro ao fazer ping no banco", zap.Error(err))
	}
	logger.Info("Conectado ao banco de dados PostgreSQL")

	// Initialize sqlc queries
	queries := db.New(dbPool)

	// Initialize repositories
	metaMensalRepo := postgres.NewMetaMensalRepository(queries)
	metaBarbeiroRepo := postgres.NewMetaBarbeiroRepository(queries)
	metaTicketMedioRepo := postgres.NewMetasTicketMedioRepository(queries)

	// Pricing repositories
	precificacaoConfigRepo := postgres.NewPrecificacaoConfigRepository(queries)
	precificacaoSimulacaoRepo := postgres.NewPrecificacaoSimulacaoRepository(queries)

	// Financial repositories
	contaPagarRepo := postgres.NewContaPagarRepository(queries)
	contaReceberRepo := postgres.NewContaReceberRepository(queries)
	compensacaoRepo := postgres.NewCompensacaoBancariaRepository(queries)
	fluxoCaixaRepo := postgres.NewFluxoCaixaDiarioRepository(queries)
	dreRepo := postgres.NewDREMensalRepository(queries)

	// Initialize use cases - Meta Mensal
	setMetaMensalUC := metas.NewSetMetaMensalUseCase(metaMensalRepo, logger)
	getMetaMensalUC := metas.NewGetMetaMensalUseCase(metaMensalRepo, logger)
	listMetasMensaisUC := metas.NewListMetasMensaisUseCase(metaMensalRepo, logger)
	updateMetaMensalUC := metas.NewUpdateMetaMensalUseCase(metaMensalRepo, logger)
	deleteMetaMensalUC := metas.NewDeleteMetaMensalUseCase(metaMensalRepo, logger)

	// Initialize use cases - Meta Barbeiro
	setMetaBarbeiroUC := metas.NewSetMetaBarbeiroUseCase(metaBarbeiroRepo, logger)
	getMetaBarbeiroUC := metas.NewGetMetaBarbeiroUseCase(metaBarbeiroRepo, logger)
	listMetasBarbeiroUC := metas.NewListMetasBarbeiroUseCase(metaBarbeiroRepo, logger)
	updateMetaBarbeiroUC := metas.NewUpdateMetaBarbeiroUseCase(metaBarbeiroRepo, logger)
	deleteMetaBarbeiroUC := metas.NewDeleteMetaBarbeiroUseCase(metaBarbeiroRepo, logger)

	// Initialize use cases - Meta Ticket Médio
	setMetaTicketMedioUC := metas.NewSetMetaTicketUseCase(metaTicketMedioRepo, logger)
	getMetaTicketMedioUC := metas.NewGetMetaTicketMedioUseCase(metaTicketMedioRepo, logger)
	listMetasTicketMedioUC := metas.NewListMetasTicketMedioUseCase(metaTicketMedioRepo, logger)
	updateMetaTicketMedioUC := metas.NewUpdateMetaTicketMedioUseCase(metaTicketMedioRepo, logger)
	deleteMetaTicketMedioUC := metas.NewDeleteMetaTicketMedioUseCase(metaTicketMedioRepo, logger)

	// Initialize use cases - Pricing (9 use cases)
	saveConfigUC := pricing.NewSaveConfigPrecificacaoUseCase(precificacaoConfigRepo, logger)
	getConfigUC := pricing.NewGetPrecificacaoConfigUseCase(precificacaoConfigRepo, logger)
	updateConfigUC := pricing.NewUpdatePrecificacaoConfigUseCase(precificacaoConfigRepo, logger)
	deleteConfigUC := pricing.NewDeletePrecificacaoConfigUseCase(precificacaoConfigRepo, logger)
	simularPrecoUC := pricing.NewSimularPrecoUseCase(precificacaoConfigRepo, precificacaoSimulacaoRepo, logger)
	saveSimulacaoUC := pricing.NewSaveSimulacaoUseCase(precificacaoSimulacaoRepo, logger)
	getSimulacaoUC := pricing.NewGetSimulacaoUseCase(precificacaoSimulacaoRepo, logger)
	listSimulacoesUC := pricing.NewListSimulacoesUseCase(precificacaoSimulacaoRepo, logger)
	deleteSimulacaoUC := pricing.NewDeleteSimulacaoUseCase(precificacaoSimulacaoRepo, logger)

	// Initialize use cases - Financial (23 use cases)
	// ContaPagar
	createContaPagarUC := financial.NewCreateContaPagarUseCase(contaPagarRepo, logger)
	getContaPagarUC := financial.NewGetContaPagarUseCase(contaPagarRepo, logger)
	listContasPagarUC := financial.NewListContasPagarUseCase(contaPagarRepo, logger)
	updateContaPagarUC := financial.NewUpdateContaPagarUseCase(contaPagarRepo, logger)
	deleteContaPagarUC := financial.NewDeleteContaPagarUseCase(contaPagarRepo, logger)
	marcarPagamentoUC := financial.NewMarcarPagamentoUseCase(contaPagarRepo, logger)
	// ContaReceber
	createContaReceberUC := financial.NewCreateContaReceberUseCase(contaReceberRepo, logger)
	getContaReceberUC := financial.NewGetContaReceberUseCase(contaReceberRepo, logger)
	listContasReceberUC := financial.NewListContasReceberUseCase(contaReceberRepo, logger)
	updateContaReceberUC := financial.NewUpdateContaReceberUseCase(contaReceberRepo, logger)
	deleteContaReceberUC := financial.NewDeleteContaReceberUseCase(contaReceberRepo, logger)
	marcarRecebimentoUC := financial.NewMarcarRecebimentoUseCase(contaReceberRepo, logger)
	// Compensação
	createCompensacaoUC := financial.NewCreateCompensacaoUseCase(compensacaoRepo, logger)
	getCompensacaoUC := financial.NewGetCompensacaoUseCase(compensacaoRepo, logger)
	listCompensacoesUC := financial.NewListCompensacoesUseCase(compensacaoRepo, logger)
	deleteCompensacaoUC := financial.NewDeleteCompensacaoUseCase(compensacaoRepo, logger)
	marcarCompensacaoUC := financial.NewMarcarCompensacaoUseCase(compensacaoRepo, logger)
	// FluxoCaixa (com dependências de ContaPagar e ContaReceber)
	generateFluxoDiarioUC := financial.NewGenerateFluxoDiarioUseCase(fluxoCaixaRepo, contaPagarRepo, contaReceberRepo, logger)
	getFluxoCaixaUC := financial.NewGetFluxoCaixaUseCase(fluxoCaixaRepo, logger)
	listFluxoCaixaUC := financial.NewListFluxoCaixaUseCase(fluxoCaixaRepo, logger)
	// DRE (com dependências de ContaPagar e ContaReceber)
	generateDREUC := financial.NewGenerateDREUseCase(dreRepo, contaPagarRepo, contaReceberRepo, logger)
	getDREUC := financial.NewGetDREUseCase(dreRepo, logger)
	listDREUC := financial.NewListDREUseCase(dreRepo, logger)

	// Initialize scheduler for cron jobs
	sched := scheduler.New(logger)

	// Register financial cron jobs
	financialDeps := scheduler.FinancialJobDeps{
		GenerateDRE:         generateDREUC,
		GenerateFluxoDiario: generateFluxoDiarioUC,
		MarcarCompensacoes:  marcarCompensacaoUC,
	}

	// Parse tenant list from ENV (SCHEDULER_TENANTS="tenant1,tenant2,...")
	tenants := scheduler.ParseTenantEnv("SCHEDULER_TENANTS")

	scheduler.RegisterFinancialJobs(sched, logger, financialDeps, tenants)

	// Start scheduler in background
	sched.Start()
	defer sched.Stop(ctx)

	// Initialize handlers - Metas completo (15 use cases)
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
		logger,
	)

	// Initialize handlers - Pricing (9 use cases)
	pricingHandler := handler.NewPricingHandler(
		// Config
		saveConfigUC,
		getConfigUC,
		updateConfigUC,
		deleteConfigUC,
		// Simulação
		simularPrecoUC,
		saveSimulacaoUC,
		getSimulacaoUC,
		listSimulacoesUC,
		deleteSimulacaoUC,
		logger,
	)

	// Initialize handlers - Financial (23 use cases)
	financialHandler := handler.NewFinancialHandler(
		// ContaPagar
		createContaPagarUC,
		getContaPagarUC,
		listContasPagarUC,
		updateContaPagarUC,
		deleteContaPagarUC,
		marcarPagamentoUC,
		// ContaReceber
		createContaReceberUC,
		getContaReceberUC,
		listContasReceberUC,
		updateContaReceberUC,
		deleteContaReceberUC,
		marcarRecebimentoUC,
		// Compensação
		createCompensacaoUC,
		getCompensacaoUC,
		listCompensacoesUC,
		deleteCompensacaoUC,
		marcarCompensacaoUC,
		// FluxoCaixa
		generateFluxoDiarioUC,
		getFluxoCaixaUC,
		listFluxoCaixaUC,
		// DRE
		generateDREUC,
		getDREUC,
		listDREUC,
		logger,
	)

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:8000"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Health Check Endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status": "ok",
			"app":    "Barber Analytics Pro v2.0",
		})
	})

	// API Routes
	api := e.Group("/api/v1")

	// Middleware para injetar tenant_id (mock por enquanto - TODO: implementar JWT)
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: Extrair de JWT
			// Por enquanto, aceita via header para testes
			tenantID := c.Request().Header.Get("X-Tenant-ID")
			if tenantID == "" {
				tenantID = "00000000-0000-0000-0000-000000000001" // tenant mock
			}
			c.Set("tenant_id", tenantID)
			return next(c)
		}
	})

	// Metas routes - 15 endpoints completos
	metasGroup := api.Group("/metas")

	// Meta Mensal (5 endpoints)
	metasGroup.POST("/monthly", metasHandler.SetMetaMensal)
	metasGroup.GET("/monthly/:id", metasHandler.GetMetaMensal)
	metasGroup.GET("/monthly", metasHandler.ListMetasMensais)
	metasGroup.PUT("/monthly/:id", metasHandler.UpdateMetaMensal)
	metasGroup.DELETE("/monthly/:id", metasHandler.DeleteMetaMensal)

	// Meta Barbeiro (5 endpoints)
	metasGroup.POST("/barbers", metasHandler.SetMetaBarbeiro)
	metasGroup.GET("/barbers/:id", metasHandler.GetMetaBarbeiro)
	metasGroup.GET("/barbers", metasHandler.ListMetasBarbeiro)
	metasGroup.PUT("/barbers/:id", metasHandler.UpdateMetaBarbeiro)
	metasGroup.DELETE("/barbers/:id", metasHandler.DeleteMetaBarbeiro)

	// Meta Ticket Médio (5 endpoints)
	metasGroup.POST("/ticket", metasHandler.SetMetaTicket)
	metasGroup.GET("/ticket/:id", metasHandler.GetMetaTicket)
	metasGroup.GET("/ticket", metasHandler.ListMetasTicket)
	metasGroup.PUT("/ticket/:id", metasHandler.UpdateMetaTicket)
	metasGroup.DELETE("/ticket/:id", metasHandler.DeleteMetaTicket)

	// Pricing routes - 9 endpoints
	pricingGroup := api.Group("/pricing")
	pricingHandler.RegisterRoutes(pricingGroup)

	// Financial routes - 19 endpoints (20 total, mas 1 é cronjob)
	financialGroup := api.Group("/financial")

	// ContaPagar (6 endpoints: 5 CRUD + 1 marcarPagamento)
	financialGroup.POST("/payables", financialHandler.CreateContaPagar)
	financialGroup.GET("/payables/:id", financialHandler.GetContaPagar)
	financialGroup.GET("/payables", financialHandler.ListContasPagar)
	financialGroup.PUT("/payables/:id", financialHandler.UpdateContaPagar)
	financialGroup.DELETE("/payables/:id", financialHandler.DeleteContaPagar)
	financialGroup.POST("/payables/:id/payment", financialHandler.MarcarPagamento)

	// ContaReceber (6 endpoints: 5 CRUD + 1 marcarRecebimento)
	financialGroup.POST("/receivables", financialHandler.CreateContaReceber)
	financialGroup.GET("/receivables/:id", financialHandler.GetContaReceber)
	financialGroup.GET("/receivables", financialHandler.ListContasReceber)
	financialGroup.PUT("/receivables/:id", financialHandler.UpdateContaReceber)
	financialGroup.DELETE("/receivables/:id", financialHandler.DeleteContaReceber)
	financialGroup.POST("/receivables/:id/receipt", financialHandler.MarcarRecebimento)

	// Compensação (3 endpoints: Get, List, Delete)
	financialGroup.GET("/compensations/:id", financialHandler.GetCompensacao)
	financialGroup.GET("/compensations", financialHandler.ListCompensacoes)
	financialGroup.DELETE("/compensations/:id", financialHandler.DeleteCompensacao)

	// FluxoCaixa (2 endpoints: Get, List)
	financialGroup.GET("/cashflow/:id", financialHandler.GetFluxoCaixa)
	financialGroup.GET("/cashflow", financialHandler.ListFluxoCaixa)

	// DRE (2 endpoints: Get, List)
	financialGroup.GET("/dre/:month", financialHandler.GetDRE)
	financialGroup.GET("/dre", financialHandler.ListDRE)

	// Placeholder endpoint
	api.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "pong",
		})
	})

	// Start server
	logger.Info("🚀 Servidor iniciado", zap.String("port", port))
	if err := e.Start(":" + port); err != nil {
		logger.Fatal("Erro ao iniciar servidor", zap.Error(err))
	}
}
