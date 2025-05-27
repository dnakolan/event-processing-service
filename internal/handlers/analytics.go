package handlers

import (
	"net/http"
	"time"

	"github.com/dnakolan/event-processing-service/internal/models"
	"github.com/dnakolan/event-processing-service/internal/services"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service services.AnalyticsService
}

func NewAnalyticsHandler(service services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

func (h *AnalyticsHandler) GetAnalyticsHandler(c *gin.Context) {
	window := c.Query("window")

	filter, err := buildFilterFromWindow(&window)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	analytics, err := h.service.GetAnalytics(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, analytics)
}

func buildFilterFromWindow(window *string) (*models.EventFilter, error) {
	timeWindow, err := time.ParseDuration(*window)
	if err != nil {
		return nil, err
	}

	start := time.Now().Add(-timeWindow)
	end := time.Now()
	return &models.EventFilter{
		StartTimestamp: &start,
		EndTimestamp:   &end,
	}, nil
}
