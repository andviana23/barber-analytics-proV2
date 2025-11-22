package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/metas"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// MetasHandler agrupa os handlers de metas
type MetasHandler struct {
	// Metas Mensais
	setMetaMensalUC    *metas.SetMetaMensalUseCase
	getMetaMensalUC    *metas.GetMetaMensalUseCase
	listMetasMensaisUC *metas.ListMetasMensaisUseCase
	updateMetaMensalUC *metas.UpdateMetaMensalUseCase
	deleteMetaMensalUC *metas.DeleteMetaMensalUseCase

	// Metas Barbeiro
	setMetaBarbeiroUC    *metas.SetMetaBarbeiroUseCase
	getMetaBarbeiroUC    *metas.GetMetaBarbeiroUseCase
	listMetasBarbeiroUC  *metas.ListMetasBarbeiroUseCase
	updateMetaBarbeiroUC *metas.UpdateMetaBarbeiroUseCase
	deleteMetaBarbeiroUC *metas.DeleteMetaBarbeiroUseCase

	// Metas Ticket Médio
	setMetaTicketUC         *metas.SetMetaTicketUseCase
	getMetaTicketMedioUC    *metas.GetMetaTicketMedioUseCase
	listMetasTicketMedioUC  *metas.ListMetasTicketMedioUseCase
	updateMetaTicketMedioUC *metas.UpdateMetaTicketMedioUseCase
	deleteMetaTicketMedioUC *metas.DeleteMetaTicketMedioUseCase

	logger *zap.Logger
}

// NewMetasHandler cria um novo handler de metas
func NewMetasHandler(
	setMetaMensalUC *metas.SetMetaMensalUseCase,
	getMetaMensalUC *metas.GetMetaMensalUseCase,
	listMetasMensaisUC *metas.ListMetasMensaisUseCase,
	updateMetaMensalUC *metas.UpdateMetaMensalUseCase,
	deleteMetaMensalUC *metas.DeleteMetaMensalUseCase,
	setMetaBarbeiroUC *metas.SetMetaBarbeiroUseCase,
	getMetaBarbeiroUC *metas.GetMetaBarbeiroUseCase,
	listMetasBarbeiroUC *metas.ListMetasBarbeiroUseCase,
	updateMetaBarbeiroUC *metas.UpdateMetaBarbeiroUseCase,
	deleteMetaBarbeiroUC *metas.DeleteMetaBarbeiroUseCase,
	setMetaTicketUC *metas.SetMetaTicketUseCase,
	getMetaTicketMedioUC *metas.GetMetaTicketMedioUseCase,
	listMetasTicketMedioUC *metas.ListMetasTicketMedioUseCase,
	updateMetaTicketMedioUC *metas.UpdateMetaTicketMedioUseCase,
	deleteMetaTicketMedioUC *metas.DeleteMetaTicketMedioUseCase,
	logger *zap.Logger,
) *MetasHandler {
	return &MetasHandler{
		setMetaMensalUC:         setMetaMensalUC,
		getMetaMensalUC:         getMetaMensalUC,
		listMetasMensaisUC:      listMetasMensaisUC,
		updateMetaMensalUC:      updateMetaMensalUC,
		deleteMetaMensalUC:      deleteMetaMensalUC,
		setMetaBarbeiroUC:       setMetaBarbeiroUC,
		getMetaBarbeiroUC:       getMetaBarbeiroUC,
		listMetasBarbeiroUC:     listMetasBarbeiroUC,
		updateMetaBarbeiroUC:    updateMetaBarbeiroUC,
		deleteMetaBarbeiroUC:    deleteMetaBarbeiroUC,
		setMetaTicketUC:         setMetaTicketUC,
		getMetaTicketMedioUC:    getMetaTicketMedioUC,
		listMetasTicketMedioUC:  listMetasTicketMedioUC,
		updateMetaTicketMedioUC: updateMetaTicketMedioUC,
		deleteMetaTicketMedioUC: deleteMetaTicketMedioUC,
		logger:                  logger,
	}
}

// SetMetaMensal godoc
// @Summary Definir/atualizar meta mensal
// @Description Define ou atualiza a meta de faturamento mensal
// @Tags Metas
// @Accept json
// @Produce json
// @Param request body dto.SetMetaMensalRequest true "Dados da meta mensal"
// @Success 200 {object} dto.MetaMensalResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/metas/monthly [post]
// @Security BearerAuth
func (h *MetasHandler) SetMetaMensal(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.SetMetaMensalRequest
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

	mesAno, metaFaturamento, origem, err := mapper.FromSetMetaMensalRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "conversion_error",
			Message: err.Error(),
		})
	}

	meta, err := h.setMetaMensalUC.Execute(ctx, metas.SetMetaMensalInput{
		TenantID:        tenantID,
		MesAno:          mesAno,
		MetaFaturamento: metaFaturamento,
		Origem:          origem,
	})
	if err != nil {
		h.logger.Error("Erro ao definir meta mensal", zap.Error(err), zap.String("tenant_id", tenantID))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao definir meta mensal",
		})
	}

	response := mapper.ToMetaMensalResponse(meta)
	return c.JSON(http.StatusOK, response)
}

// SetMetaBarbeiro godoc
// @Summary Definir/atualizar meta de barbeiro
// @Description Define ou atualiza as metas individuais de um barbeiro
// @Tags Metas
// @Accept json
// @Produce json
// @Param request body dto.SetMetaBarbeiroRequest true "Dados da meta do barbeiro"
// @Success 200 {object} dto.MetaBarbeiroResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/metas/barbers [post]
// @Security BearerAuth
func (h *MetasHandler) SetMetaBarbeiro(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.SetMetaBarbeiroRequest
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

	mesAno, metaGerais, metaExtras, metaProdutos, err := mapper.FromSetMetaBarbeiroRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "conversion_error",
			Message: err.Error(),
		})
	}

	meta, err := h.setMetaBarbeiroUC.Execute(ctx, metas.SetMetaBarbeiroInput{
		TenantID:           tenantID,
		BarbeiroID:         req.BarbeiroID,
		MesAno:             mesAno,
		MetaServicosGerais: metaGerais,
		MetaServicosExtras: metaExtras,
		MetaProdutos:       metaProdutos,
	})
	if err != nil {
		h.logger.Error("Erro ao definir meta de barbeiro", zap.Error(err), zap.String("barbeiro_id", req.BarbeiroID))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao definir meta de barbeiro",
		})
	}

	response := mapper.ToMetaBarbeiroResponse(meta)
	return c.JSON(http.StatusOK, response)
}

// SetMetaTicket godoc
// @Summary Definir/atualizar meta de ticket médio
// @Description Define ou atualiza a meta de ticket médio (geral ou por barbeiro)
// @Tags Metas
// @Accept json
// @Produce json
// @Param request body dto.SetMetaTicketRequest true "Dados da meta de ticket médio"
// @Success 200 {object} dto.MetaTicketResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/metas/ticket [post]
// @Security BearerAuth
func (h *MetasHandler) SetMetaTicket(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.SetMetaTicketRequest
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

	// Parse adicional para meta ticket
	mesAno, tipo, barbeiroID, metaValor, err := mapper.FromSetMetaTicketRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	meta, err := h.setMetaTicketUC.Execute(ctx, metas.SetMetaTicketInput{
		TenantID:   tenantID,
		MesAno:     mesAno,
		Tipo:       tipo,
		BarbeiroID: barbeiroID,
		MetaValor:  metaValor,
	})
	if err != nil {
		h.logger.Error("Erro ao definir meta de ticket", zap.Error(err), zap.String("tipo", req.Tipo))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao definir meta de ticket médio",
		})
	}

	response := mapper.ToMetaTicketResponse(meta)
	return c.JSON(http.StatusOK, response)
}

// GetMetaMensal godoc
// @Summary Buscar meta mensal
// @Description Retorna uma meta mensal específica por ID
// @Tags Metas
// @Produce json
// @Param id path string true "ID da meta mensal"
// @Success 200 {object} dto.MetaMensalResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/metas/monthly/{id} [get]
// @Security BearerAuth
func (h *MetasHandler) GetMetaMensal(c echo.Context) error {
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
			Message: "ID obrigatório",
		})
	}

	meta, err := h.getMetaMensalUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar meta mensal", zap.Error(err), zap.String("id", id))
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Meta não encontrada",
		})
	}

	response := mapper.ToMetaMensalResponse(meta)
	return c.JSON(http.StatusOK, response)
}

// ListMetasMensais godoc
// @Summary Listar metas mensais
// @Description Lista todas as metas mensais do tenant
// @Tags Metas
// @Produce json
// @Success 200 {array} dto.MetaMensalResponse
// @Router /api/v1/metas/monthly [get]
// @Security BearerAuth
func (h *MetasHandler) ListMetasMensais(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// TODO: Implementar parsing de filtros de período
	// Por ora, usa defaults amplos
	inicio, _ := valueobject.NewMesAno("2020-01")
	fim, _ := valueobject.NewMesAno("2099-12")

	metasList, err := h.listMetasMensaisUC.Execute(ctx, metas.ListMetasMensaisInput{
		TenantID: tenantID,
		Inicio:   inicio,
		Fim:      fim,
	})
	if err != nil {
		h.logger.Error("Erro ao listar metas mensais", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar metas",
		})
	}

	responses := make([]dto.MetaMensalResponse, len(metasList))
	for i, meta := range metasList {
		responses[i] = mapper.ToMetaMensalResponse(meta)
	}

	return c.JSON(http.StatusOK, responses)
}

// UpdateMetaMensal godoc
// @Summary Atualizar meta mensal
// @Description Atualiza uma meta mensal existente
// @Tags Metas
// @Accept json
// @Produce json
// @Param id path string true "ID da meta"
// @Param request body dto.SetMetaMensalRequest true "Dados da meta"
// @Success 200 {object} dto.MetaMensalResponse
// @Router /api/v1/metas/monthly/{id} [put]
// @Security BearerAuth
func (h *MetasHandler) UpdateMetaMensal(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID inválido",
		})
	}

	var req dto.SetMetaMensalRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Requisição inválida",
		})
	}

	_, metaFaturamento, _, err := mapper.FromSetMetaMensalRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_data",
			Message: err.Error(),
		})
	}

	meta, err := h.updateMetaMensalUC.Execute(ctx, metas.UpdateMetaMensalInput{
		TenantID:        tenantID,
		ID:              id,
		MetaFaturamento: metaFaturamento,
	})
	if err != nil {
		h.logger.Error("Erro ao atualizar meta mensal",
			zap.String("tenant_id", tenantID),
			zap.String("id", id),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "update_failed",
			Message: "Erro ao atualizar meta",
		})
	}

	return c.JSON(http.StatusOK, mapper.ToMetaMensalResponse(meta))
}

// DeleteMetaMensal godoc
// @Summary Deletar meta mensal
// @Description Remove uma meta mensal
// @Tags Metas
// @Param id path string true "ID da meta"
// @Success 204
// @Router /api/v1/metas/monthly/{id} [delete]
// @Security BearerAuth
func (h *MetasHandler) DeleteMetaMensal(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID inválido",
		})
	}

	err := h.deleteMetaMensalUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao deletar meta mensal",
			zap.String("tenant_id", tenantID),
			zap.String("id", id),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "delete_failed",
			Message: "Erro ao deletar meta",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetMetaBarbeiro godoc
// @Summary Buscar meta de barbeiro
// @Tags Metas
// @Produce json
// @Param id path string true "ID da meta"
// @Success 200 {object} dto.MetaBarbeiroResponse
// @Router /api/v1/metas/barbers/{id} [get]
// @Security BearerAuth
func (h *MetasHandler) GetMetaBarbeiro(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID inválido",
		})
	}

	meta, err := h.getMetaBarbeiroUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar meta barbeiro",
			zap.String("tenant_id", tenantID),
			zap.String("id", id),
			zap.Error(err),
		)
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Meta não encontrada",
		})
	}

	return c.JSON(http.StatusOK, mapper.ToMetaBarbeiroResponse(meta))
}

// ListMetasBarbeiro godoc
// @Summary Listar metas de barbeiros
// @Tags Metas
// @Produce json
// @Param barbeiro_id query string false "ID do barbeiro (opcional)"
// @Success 200 {array} dto.MetaBarbeiroResponse
// @Router /api/v1/metas/barbers [get]
// @Security BearerAuth
func (h *MetasHandler) ListMetasBarbeiro(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Query params opcionais
	barbeiroID := c.QueryParam("barbeiro_id")
	var barbeiroIDPtr *string
	if barbeiroID != "" {
		barbeiroIDPtr = &barbeiroID
	}

	// Período padrão (amplo para pegar todas)
	mesAnoInicio, _ := valueobject.NewMesAno("2020-01")
	mesAnoFim, _ := valueobject.NewMesAno("2099-12")

	metas, err := h.listMetasBarbeiroUC.Execute(ctx, metas.ListMetasBarbeiroInput{
		TenantID:   tenantID,
		BarbeiroID: barbeiroIDPtr,
		Inicio:     mesAnoInicio,
		Fim:        mesAnoFim,
	})
	if err != nil {
		h.logger.Error("Erro ao listar metas barbeiro",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "list_failed",
			Message: "Erro ao listar metas",
		})
	}

	// Converter para DTOs
	responses := make([]dto.MetaBarbeiroResponse, len(metas))
	for i, meta := range metas {
		responses[i] = mapper.ToMetaBarbeiroResponse(meta)
	}

	return c.JSON(http.StatusOK, responses)
}

// UpdateMetaBarbeiro godoc
// @Summary Atualizar meta de barbeiro
// @Tags Metas
// @Accept json
// @Produce json
// @Param id path string true "ID da meta"
// @Param request body dto.SetMetaBarbeiroRequest true "Dados atualizados"
// @Success 200 {object} dto.MetaBarbeiroResponse
// @Router /api/v1/metas/barbers/{id} [put]
// @Security BearerAuth
func (h *MetasHandler) UpdateMetaBarbeiro(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID inválido",
		})
	}

	var req dto.SetMetaBarbeiroRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Requisição inválida",
		})
	}

	_, metaServicosGerais, metaServicosExtras, metaProdutos, err := mapper.FromSetMetaBarbeiroRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_data",
			Message: err.Error(),
		})
	}

	// Calcular meta total de faturamento
	metaTotal := metaServicosGerais.Add(metaServicosExtras).Add(metaProdutos)

	meta, err := h.updateMetaBarbeiroUC.Execute(ctx, metas.UpdateMetaBarbeiroInput{
		TenantID:        tenantID,
		ID:              id,
		MetaFaturamento: metaTotal,
		MetaServicos:    0, // TODO: extrair do request se necessário
	})
	if err != nil {
		h.logger.Error("Erro ao atualizar meta barbeiro",
			zap.String("tenant_id", tenantID),
			zap.String("id", id),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "update_failed",
			Message: "Erro ao atualizar meta",
		})
	}

	return c.JSON(http.StatusOK, mapper.ToMetaBarbeiroResponse(meta))
}

// DeleteMetaBarbeiro godoc
// @Summary Deletar meta de barbeiro
// @Tags Metas
// @Param id path string true "ID da meta"
// @Success 204
// @Router /api/v1/metas/barbers/{id} [delete]
// @Security BearerAuth
func (h *MetasHandler) DeleteMetaBarbeiro(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID inválido",
		})
	}

	err := h.deleteMetaBarbeiroUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao deletar meta barbeiro",
			zap.String("tenant_id", tenantID),
			zap.String("id", id),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "delete_failed",
			Message: "Erro ao deletar meta",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetMetaTicket godoc
// @Summary Buscar meta de ticket médio
// @Tags Metas
// @Produce json
// @Param id path string true "ID da meta"
// @Success 200 {object} dto.MetaTicketResponse
// @Router /api/v1/metas/ticket/{id} [get]
// @Security BearerAuth
func (h *MetasHandler) GetMetaTicket(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID inválido",
		})
	}

	meta, err := h.getMetaTicketMedioUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar meta ticket médio",
			zap.String("tenant_id", tenantID),
			zap.String("id", id),
			zap.Error(err),
		)
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Meta não encontrada",
		})
	}

	return c.JSON(http.StatusOK, mapper.ToMetaTicketResponse(meta))
}

// ListMetasTicket godoc
// @Summary Listar metas de ticket médio
// @Tags Metas
// @Produce json
// @Success 200 {array} dto.MetaTicketResponse
// @Router /api/v1/metas/ticket [get]
// @Security BearerAuth
func (h *MetasHandler) ListMetasTicket(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Período padrão
	mesAnoInicio, _ := valueobject.NewMesAno("2020-01")
	mesAnoFim, _ := valueobject.NewMesAno("2099-12")

	metas, err := h.listMetasTicketMedioUC.Execute(ctx, metas.ListMetasTicketMedioInput{
		TenantID: tenantID,
		Inicio:   mesAnoInicio,
		Fim:      mesAnoFim,
	})
	if err != nil {
		h.logger.Error("Erro ao listar metas ticket médio",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "list_failed",
			Message: "Erro ao listar metas",
		})
	}

	// Converter para DTOs
	responses := make([]dto.MetaTicketResponse, len(metas))
	for i, meta := range metas {
		responses[i] = mapper.ToMetaTicketResponse(meta)
	}

	return c.JSON(http.StatusOK, responses)
}

// UpdateMetaTicket godoc
// @Summary Atualizar meta de ticket médio
// @Tags Metas
// @Accept json
// @Produce json
// @Param id path string true "ID da meta"
// @Param request body dto.SetMetaTicketRequest true "Dados atualizados"
// @Success 200 {object} dto.MetaTicketResponse
// @Router /api/v1/metas/ticket/{id} [put]
// @Security BearerAuth
func (h *MetasHandler) UpdateMetaTicket(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID inválido",
		})
	}

	var req dto.SetMetaTicketRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Requisição inválida",
		})
	}

	_, _, _, metaValor, err := mapper.FromSetMetaTicketRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_data",
			Message: err.Error(),
		})
	}

	meta, err := h.updateMetaTicketMedioUC.Execute(ctx, metas.UpdateMetaTicketMedioInput{
		TenantID:        tenantID,
		ID:              id,
		MetaTicketMedio: metaValor,
	})
	if err != nil {
		h.logger.Error("Erro ao atualizar meta ticket médio",
			zap.String("tenant_id", tenantID),
			zap.String("id", id),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "update_failed",
			Message: "Erro ao atualizar meta",
		})
	}

	return c.JSON(http.StatusOK, mapper.ToMetaTicketResponse(meta))
}

// DeleteMetaTicket godoc
// @Summary Deletar meta de ticket médio
// @Tags Metas
// @Param id path string true "ID da meta"
// @Success 204
// @Router /api/v1/metas/ticket/{id} [delete]
// @Security BearerAuth
func (h *MetasHandler) DeleteMetaTicket(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID inválido",
		})
	}

	err := h.deleteMetaTicketMedioUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao deletar meta ticket médio",
			zap.String("tenant_id", tenantID),
			zap.String("id", id),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "delete_failed",
			Message: "Erro ao deletar meta",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registra todas as rotas de metas
func (h *MetasHandler) RegisterRoutes(g *echo.Group) {
	// Metas mensais
	g.POST("/monthly", h.SetMetaMensal)
	g.GET("/monthly/:id", h.GetMetaMensal)
	g.GET("/monthly", h.ListMetasMensais)
	g.PUT("/monthly/:id", h.UpdateMetaMensal)
	g.DELETE("/monthly/:id", h.DeleteMetaMensal)

	// Metas de barbeiro
	g.POST("/barbers", h.SetMetaBarbeiro)
	g.GET("/barbers/:id", h.GetMetaBarbeiro)
	g.GET("/barbers", h.ListMetasBarbeiro)
	g.PUT("/barbers/:id", h.UpdateMetaBarbeiro)
	g.DELETE("/barbers/:id", h.DeleteMetaBarbeiro)

	// Metas de ticket médio
	g.POST("/ticket", h.SetMetaTicket)
	g.GET("/ticket/:id", h.GetMetaTicket)
	g.GET("/ticket", h.ListMetasTicket)
	g.PUT("/ticket/:id", h.UpdateMetaTicket)
	g.DELETE("/ticket/:id", h.DeleteMetaTicket)
}
