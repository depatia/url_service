package delivery

import (
	"bytes"
	"fmt"
	"health_checker/internal/domain/dto"
	"health_checker/internal/domain/interfaces"
	"net/http"

	_ "health_checker/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HTTPHandler struct {
	logger           *logrus.Logger
	healthCheckerSvc interfaces.IHealthCheckerService
	pdfSvc           interfaces.IPdfService
}

func NewHTTPHandler(logger *logrus.Logger, healthCheckerSvc interfaces.IHealthCheckerService, pdfSvc interfaces.IPdfService) *HTTPHandler {
	return &HTTPHandler{
		logger:           logger,
		healthCheckerSvc: healthCheckerSvc,
		pdfSvc:           pdfSvc,
	}
}

func (h HTTPHandler) NewRouters() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Документация сваггер
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	api.POST("/check-sites", h.handleCheckSitesAvailability)
	api.POST("/generate-report", h.handleGenerateReport)

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	return router
}

// handleCheckSitesAvailability проверяет состояние переданных сайтов и возвращает их статус, а также уникальный номер запроса
// @Summary проверяет состояние переданных сайтов и возвращает их статус, а также уникальный номер запроса
// @Accept application/json
// @Produce application/json
// @Param body body dto.CheckSitesAvailabilityRequest true "Ссылки для проверки"
// @Success 200 {object} dto.CheckSitesAvailabilityResponse
// @Failure 400 {object} map[string]string
// @Router /api/check-sites [post]
func (h *HTTPHandler) handleCheckSitesAvailability(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CheckSitesAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Валидация ссылок
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Debugf("Start checking %d sites", len(req.Links))

	// Получение статусов сайтов и linksID
	links, linksID := h.healthCheckerSvc.CheckSitesAvailability(ctx, req.Links)

	// Создаем DTO ответа
	resp := dto.CheckSitesAvailabilityResponse{
		Links:    links,
		LinksNum: linksID,
	}

	h.logger.Debugf("Request %d successfully finished", linksID)

	c.JSON(http.StatusOK, resp)
}

// handleGenerateReport возвращает отчет в виде pdf о проверенных сайтах (из кэша)
// @Summary возвращает отчет в виде pdf о проверенных сайтах (из кэша)
// @Accept application/json
// @Produce application/pdf
// @Param body body dto.GenerateReportRequest true "Уникальный идентификатор запроса"
// @Success 200 {file} binary
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/generate-report [post]
func (h *HTTPHandler) handleGenerateReport(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Валидация ссылок
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pdfText := ""

	for _, linksID := range req.LinksList {
		links := h.healthCheckerSvc.GetSitesAvailabilityReportByLinksID(ctx, linksID)
		if len(links) == 0 {
			h.logger.Warnf("Nothing to generate by linksID: %d", linksID)
		}
		for link, availability := range links {
			pdfText += fmt.Sprintf("%s: %s\n", link, availability)
		}
	}

	if len(pdfText) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nothing to generate"})
		return
	}

	pdfBytes, err := h.pdfSvc.GeneratePdf(ctx, pdfText)
	if err != nil {
		h.logger.Errorf("Failed to generate pdf with text: %s", pdfText)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.DataFromReader(http.StatusOK, int64(len(pdfBytes)), "application/pdf", bytes.NewReader(pdfBytes), nil)
}
