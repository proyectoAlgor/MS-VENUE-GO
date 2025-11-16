package handlers

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, location)
}

func (h *VenueHandler) GetAllLocations(c *gin.Context) {
	locations, err := h.venueService.GetAllLocations()
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

	tables, err := h.venueService.GetTablesByLocation(locationID)
	if err != nil {
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
