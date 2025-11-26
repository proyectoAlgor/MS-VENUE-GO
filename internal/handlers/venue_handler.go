package handlers

import (
	"net/http"
	"strings"

	"ms-venue-go/internal/models"
	"ms-venue-go/internal/service"

	"github.com/gin-gonic/gin"
)

type VenueHandler struct {
	venueService service.VenueService
}

func NewVenueHandler(venueService service.VenueService) *VenueHandler {
	return &VenueHandler{
		venueService: venueService,
	}
}

// Location handlers
func (h *VenueHandler) CreateLocation(c *gin.Context) {
	var req models.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := h.venueService.CreateLocation(&req)
	if err != nil {
		// Check if it's a duplicate error
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, location)
}

func (h *VenueHandler) GetAllLocations(c *gin.Context) {
	// Get user ID and token from context (set by OptionalAuth middleware)
	userID, exists := c.Get("user_id")

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	token := ""
	if authHeader != "" {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// If no user ID in context but there's a token, try to use it anyway
	// This handles cases where OptionalAuth might not have set user_id
	userIDStr := ""
	if exists {
		userIDStr = userID.(string)
	}

	locations, err := h.venueService.GetAllLocations(userIDStr, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, locations)
}

func (h *VenueHandler) GetLocationByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location ID is required"})
		return
	}

	location, err := h.venueService.GetLocationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	c.JSON(http.StatusOK, location)
}

func (h *VenueHandler) UpdateLocation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location ID is required"})
		return
	}

	var req models.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := h.venueService.UpdateLocation(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, location)
}

func (h *VenueHandler) DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location ID is required"})
		return
	}

	err := h.venueService.DeleteLocation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location deleted successfully"})
}

// Table handlers
func (h *VenueHandler) CreateTable(c *gin.Context) {
	var req models.CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table, err := h.venueService.CreateTable(&req)
	if err != nil {
		// Check if it's a duplicate error
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, table)
}

func (h *VenueHandler) GetTablesByLocation(c *gin.Context) {
	locationID := c.Param("locationId")
	if locationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location ID is required"})
		return
	}

	// Get user ID and token from context (set by auth middleware if authenticated)
	userID, exists := c.Get("user_id")
	userIDStr := ""
	if exists {
		userIDStr = userID.(string)
	}

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	token := ""
	if authHeader != "" {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	}

	tables, err := h.venueService.GetTablesByLocation(locationID, userIDStr, token)
	if err != nil {
		// Check if it's an access denied error
		if strings.Contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}

func (h *VenueHandler) GetTableByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
		return
	}

	table, err := h.venueService.GetTableByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	c.JSON(http.StatusOK, table)
}

func (h *VenueHandler) UpdateTable(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
		return
	}

	var req models.UpdateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table, err := h.venueService.UpdateTable(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

func (h *VenueHandler) UpdateTableStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=available occupied"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateReq := &models.UpdateTableRequest{
		Status: req.Status,
	}

	table, err := h.venueService.UpdateTable(id, updateReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

func (h *VenueHandler) DeleteTable(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
		return
	}

	err := h.venueService.DeleteTable(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Table deleted successfully"})
}

// Health check
func (h *VenueHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "ms-venue-go",
	})
}
