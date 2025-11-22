package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/pricing"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// PricingHandler agrupa os handlers de precificação
type PricingHandler struct {
	// Config use cases
	saveConfigUC   *pricing.SaveConfigPrecificacaoUseCase
	getConfigUC    *pricing.GetPrecificacaoConfigUseCase
	updateConfigUC *pricing.UpdatePrecificacaoConfigUseCase
	deleteConfigUC *pricing.DeletePrecificacaoConfigUseCase
	// Simulação use cases
	simularPrecoUC    *pricing.SimularPrecoUseCase
	saveSimulacaoUC   *pricing.SaveSimulacaoUseCase
	getSimulacaoUC    *pricing.GetSimulacaoUseCase
	listSimulacoesUC  *pricing.ListSimulacoesUseCase
	deleteSimulacaoUC *pricing.DeleteSimulacaoUseCase
	logger            *zap.Logger
}

// NewPricingHandler cria um novo handler de precificação
func NewPricingHandler(
	// Config
	saveConfigUC *pricing.SaveConfigPrecificacaoUseCase,
	getConfigUC *pricing.GetPrecificacaoConfigUseCase,
	updateConfigUC *pricing.UpdatePrecificacaoConfigUseCase,
	deleteConfigUC *pricing.DeletePrecificacaoConfigUseCase,
	// Simulação
	simularPrecoUC *pricing.SimularPrecoUseCase,
	saveSimulacaoUC *pricing.SaveSimulacaoUseCase,
	getSimulacaoUC *pricing.GetSimulacaoUseCase,
	listSimulacoesUC *pricing.ListSimulacoesUseCase,
	deleteSimulacaoUC *pricing.DeleteSimulacaoUseCase,
	logger *zap.Logger,
) *PricingHandler {
	return &PricingHandler{
		saveConfigUC:      saveConfigUC,
		getConfigUC:       getConfigUC,
		updateConfigUC:    updateConfigUC,
		deleteConfigUC:    deleteConfigUC,
		simularPrecoUC:    simularPrecoUC,
		saveSimulacaoUC:   saveSimulacaoUC,
		getSimulacaoUC:    getSimulacaoUC,
		listSimulacoesUC:  listSimulacoesUC,
		deleteSimulacaoUC: deleteSimulacaoUC,
		logger:            logger,
	}
}

// SaveConfig godoc
// @Summary Salvar configuração de precificação
// @Description Salva ou atualiza a configuração de precificação do tenant
// @Tags Pricing
// @Accept json
// @Produce json
// @Param request body dto.SaveConfigPrecificacaoRequest true "Configuração de precificação"
// @Success 200 {object} dto.PrecificacaoConfigResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/pricing/config [post]
// @Security BearerAuth
func (h *PricingHandler) SaveConfig(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.SaveConfigPrecificacaoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	margemDesejada, impostoPercentual, comissaoDefault, markupAlvo, err := mapper.FromSaveConfigPrecificacaoRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "conversion_error",
			Message: err.Error(),
		})
	}

	config, err := h.saveConfigUC.Execute(ctx, pricing.SaveConfigPrecificacaoInput{
		TenantID:          tenantID,
		MargemDesejada:    margemDesejada,
		MarkupAlvo:        markupAlvo,
		ImpostoPercentual: impostoPercentual,
		ComissaoDefault:   comissaoDefault,
	})
	if err != nil {
		h.logger.Error("Erro ao salvar configuração de precificação", zap.Error(err), zap.String("tenant_id", tenantID))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao salvar configuração",
		})
	}

	response := mapper.ToPrecificacaoConfigResponse(config)
	return c.JSON(http.StatusOK, response)
}

// SimularPreco godoc
// @Summary Simular preço
// @Description Simula o preço sugerido com base nos custos e configuração
// @Tags Pricing
// @Accept json
// @Produce json
// @Param request body dto.SimularPrecoRequest true "Dados para simulação"
// @Success 200 {object} dto.PrecificacaoSimulacaoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/pricing/simulate [post]
// @Security BearerAuth
func (h *PricingHandler) SimularPreco(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.SimularPrecoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	custoMateriais, custoMaoDeObra, precoAtual, parametros, err := mapper.FromSimularPrecoRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "conversion_error",
			Message: err.Error(),
		})
	}

	simulacao, err := h.simularPrecoUC.Execute(ctx, pricing.SimularPrecoInput{
		TenantID:       tenantID,
		ItemID:         req.ItemID,
		TipoItem:       req.TipoItem,
		CustoMateriais: custoMateriais,
		CustoMaoDeObra: custoMaoDeObra,
		PrecoAtual:     precoAtual,
		Parametros:     parametros,
	})
	if err != nil {
		h.logger.Error("Erro ao simular preço", zap.Error(err), zap.String("item_id", req.ItemID))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao simular preço",
		})
	}

	response := mapper.ToPrecificacaoSimulacaoResponse(simulacao)
	return c.JSON(http.StatusOK, response)
}

// GetConfig godoc
// @Summary Buscar configuração de precificação
// @Tags Pricing
// @Produce json
// @Success 200 {object} dto.PrecificacaoConfigResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/pricing/config [get]
// @Security BearerAuth
func (h *PricingHandler) GetConfig(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	config, err := h.getConfigUC.Execute(ctx, tenantID)
	if err != nil {
		h.logger.Error("Erro ao buscar configuração", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar configuração",
		})
	}

	response := mapper.ToPrecificacaoConfigResponse(config)
	return c.JSON(http.StatusOK, response)
}

// UpdateConfig godoc
// @Summary Atualizar configuração de precificação
// @Tags Pricing
// @Accept json
// @Produce json
// @Param request body dto.SaveConfigPrecificacaoRequest true "Configuração atualizada"
// @Success 200 {object} dto.PrecificacaoConfigResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/pricing/config [put]
// @Security BearerAuth
func (h *PricingHandler) UpdateConfig(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.SaveConfigPrecificacaoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	// Buscar config atual para atualizar
	configAtual, err := h.getConfigUC.Execute(ctx, tenantID)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Configuração não encontrada",
		})
	}

	margemDesejada, impostoPercentual, comissaoDefault, markupAlvo, err := mapper.FromSaveConfigPrecificacaoRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "conversion_error",
			Message: err.Error(),
		})
	}

	// Atualizar campos
	configAtual.MargemDesejada = margemDesejada
	configAtual.MarkupAlvo = markupAlvo
	configAtual.ImpostoPercentual = impostoPercentual
	configAtual.ComissaoPercentualDefault = comissaoDefault

	configAtualizada, err := h.updateConfigUC.Execute(ctx, pricing.UpdatePrecificacaoConfigInput{
		TenantID: tenantID,
		Config:   configAtual,
	})
	if err != nil {
		h.logger.Error("Erro ao atualizar configuração", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao atualizar configuração",
		})
	}

	response := mapper.ToPrecificacaoConfigResponse(configAtualizada)
	return c.JSON(http.StatusOK, response)
}

// DeleteConfig godoc
// @Summary Deletar configuração de precificação
// @Tags Pricing
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/pricing/config [delete]
// @Security BearerAuth
func (h *PricingHandler) DeleteConfig(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	if err := h.deleteConfigUC.Execute(ctx, tenantID); err != nil {
		h.logger.Error("Erro ao deletar configuração", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao deletar configuração",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetSimulacao godoc
// @Summary Buscar simulação de precificação
// @Tags Pricing
// @Produce json
// @Param id path string true "ID da simulação"
// @Success 200 {object} dto.PrecificacaoSimulacaoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/pricing/simulations/{id} [get]
// @Security BearerAuth
func (h *PricingHandler) GetSimulacao(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID é obrigatório",
		})
	}

	sim, err := h.getSimulacaoUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar simulação", zap.Error(err), zap.String("id", id))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar simulação",
		})
	}

	response := mapper.ToPrecificacaoSimulacaoResponse(sim)
	return c.JSON(http.StatusOK, response)
}

// ListSimulacoes godoc
// @Summary Listar simulações de precificação
// @Tags Pricing
// @Produce json
// @Param item_id query string false "Filtrar por item"
// @Param tipo_item query string false "Filtrar por tipo de item"
// @Param page query int false "Número da página" default(1)
// @Param page_size query int false "Tamanho da página" default(10)
// @Success 200 {array} dto.PrecificacaoSimulacaoResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/pricing/simulations [get]
// @Security BearerAuth
func (h *PricingHandler) ListSimulacoes(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Query params opcionais
	itemID := c.QueryParam("item_id")
	tipoItem := c.QueryParam("tipo_item")

	input := pricing.ListSimulacoesInput{
		TenantID: tenantID,
	}
	if itemID != "" {
		input.ItemID = &itemID
	}
	if tipoItem != "" {
		input.TipoItem = &tipoItem
	}

	// Aceita paginação mas não precisa ser usada (lista vazia é válida)
	// page := c.QueryParam("page")
	// pageSize := c.QueryParam("page_size")

	sims, err := h.listSimulacoesUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao listar simulações", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar simulações",
		})
	}

	// Se não houver simulações, retorna array vazio (válido)
	responses := make([]dto.PrecificacaoSimulacaoResponse, len(sims))
	for i, sim := range sims {
		responses[i] = mapper.ToPrecificacaoSimulacaoResponse(sim)
	}

	return c.JSON(http.StatusOK, responses)
}

// DeleteSimulacao godoc
// @Summary Deletar simulação de precificação
// @Tags Pricing
// @Param id path string true "ID da simulação"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/pricing/simulations/{id} [delete]
// @Security BearerAuth
func (h *PricingHandler) DeleteSimulacao(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID é obrigatório",
		})
	}

	if err := h.deleteSimulacaoUC.Execute(ctx, tenantID, id); err != nil {
		h.logger.Error("Erro ao deletar simulação", zap.Error(err), zap.String("id", id))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao deletar simulação",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registra todas as rotas de precificação
func (h *PricingHandler) RegisterRoutes(g *echo.Group) {
	// Configuração
	g.POST("/config", h.SaveConfig)
	g.GET("/config", h.GetConfig)
	g.PUT("/config", h.UpdateConfig)
	g.DELETE("/config", h.DeleteConfig)

	// Simulações
	g.POST("/simulate", h.SimularPreco)
	g.GET("/simulations/:id", h.GetSimulacao)
	g.GET("/simulations", h.ListSimulacoes)
	g.DELETE("/simulations/:id", h.DeleteSimulacao)
}
