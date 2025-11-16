package service

import (
	"fmt"
	"ms-venue-go/internal/models"
	"ms-venue-go/internal/repository"

	"github.com/google/uuid"
)

type VenueService interface {
	// Location operations
	CreateLocation(req *models.CreateLocationRequest) (*models.Location, error)
	GetAllLocations() ([]models.Location, error)
	GetLocationByID(id string) (*models.Location, error)
	UpdateLocation(id string, req *models.UpdateLocationRequest) (*models.Location, error)
	DeleteLocation(id string) error

	// Table operations
	CreateTable(req *models.CreateTableRequest) (*models.Table, error)
	GetTablesByLocation(locationID string) ([]models.Table, error)
	GetTableByID(id string) (*models.Table, error)
	UpdateTable(id string, req *models.UpdateTableRequest) (*models.Table, error)
	DeleteTable(id string) error
}

type venueService struct {
	repo repository.VenueRepository
}

func NewVenueService(repo repository.VenueRepository) VenueService {
	return &venueService{
		repo: repo,
	}
}

// Location operations
func (s *venueService) CreateLocation(req *models.CreateLocationRequest) (*models.Location, error) {
	location := &models.Location{
		ID:       generateUUID(),
		Code:     req.Code,
		Name:     req.Name,
		Address:  req.Address,
		IsActive: true,
	}

	err := s.repo.CreateLocation(location)
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return location, nil
}

func (s *venueService) GetAllLocations() ([]models.Location, error) {
	locations, err := s.repo.GetAllLocations()
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}
	return locations, nil
}

func (s *venueService) GetLocationByID(id string) (*models.Location, error) {
	location, err := s.repo.GetLocationByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}
	return location, nil
}

func (s *venueService) UpdateLocation(id string, req *models.UpdateLocationRequest) (*models.Location, error) {
	// Get existing location
	location, err := s.repo.GetLocationByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	// Update fields if provided
	if req.Code != "" {
		location.Code = req.Code
	}
	if req.Name != "" {
		location.Name = req.Name
	}
	if req.Address != "" {
		location.Address = req.Address
	}
	if req.IsActive != nil {
		location.IsActive = *req.IsActive
	}

	err = s.repo.UpdateLocation(id, location)
	if err != nil {
		return nil, fmt.Errorf("failed to update location: %w", err)
	}

	return location, nil
}

func (s *venueService) DeleteLocation(id string) error {
	err := s.repo.DeleteLocation(id)
	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}
	return nil
}

// Table operations
func (s *venueService) CreateTable(req *models.CreateTableRequest) (*models.Table, error) {
	table := &models.Table{
		ID:         generateUUID(),
		LocationID: req.LocationID,
		Code:       req.Code,
		Seats:      req.Seats,
		Status:     "free",
		IsActive:   true,
	}

	err := s.repo.CreateTable(table)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return table, nil
}

func (s *venueService) GetTablesByLocation(locationID string) ([]models.Table, error) {
	tables, err := s.repo.GetTablesByLocation(locationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	return tables, nil
}

func (s *venueService) GetTableByID(id string) (*models.Table, error) {
	table, err := s.repo.GetTableByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}
	return table, nil
}

func (s *venueService) UpdateTable(id string, req *models.UpdateTableRequest) (*models.Table, error) {
	// Get existing table
	table, err := s.repo.GetTableByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	// Update fields if provided
	if req.Code != "" {
		table.Code = req.Code
	}
	if req.Seats > 0 {
		table.Seats = req.Seats
	}
	if req.Status != "" {
		table.Status = req.Status
	}
	if req.IsActive != nil {
		table.IsActive = *req.IsActive
	}

	err = s.repo.UpdateTable(id, table)
	if err != nil {
		return nil, fmt.Errorf("failed to update table: %w", err)
	}

	return table, nil
}

func (s *venueService) DeleteTable(id string) error {
	err := s.repo.DeleteTable(id)
	if err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}
	return nil
}

// Generate UUID v4
func generateUUID() string {
	return uuid.New().String()
}
