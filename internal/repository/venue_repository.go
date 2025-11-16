package repository

import (
	"database/sql"
	"fmt"
	"ms-venue-go/internal/models"

	_ "github.com/lib/pq"
)

type VenueRepository interface {
	// Location operations
	CreateLocation(location *models.Location) error
	GetAllLocations() ([]models.Location, error)
	GetLocationByID(id string) (*models.Location, error)
	UpdateLocation(id string, location *models.Location) error
	DeleteLocation(id string) error

	// Table operations
	CreateTable(table *models.Table) error
	GetTablesByLocation(locationID string) ([]models.Table, error)
	GetTableByID(id string) (*models.Table, error)
	UpdateTable(id string, table *models.Table) error
	DeleteTable(id string) error
}

type venueRepository struct {
	db *sql.DB
}

func NewVenueRepository(dbURL string) (VenueRepository, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &venueRepository{
		db: db,
	}

	return repo, nil
}

// Location operations
func (r *venueRepository) CreateLocation(location *models.Location) error {
	query := `
		INSERT INTO bar_system.locations (id, code, name, address, is_active)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (code) DO NOTHING
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(query, location.ID, location.Code, location.Name,
		location.Address, location.IsActive).Scan(&location.CreatedAt, &location.UpdatedAt)

	return err
}

func (r *venueRepository) GetAllLocations() ([]models.Location, error) {
	query := `
		SELECT id, code, name, address, is_active, created_at, updated_at
		FROM bar_system.locations 
		WHERE is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := make([]models.Location, 0)
	for rows.Next() {
		var location models.Location
		err := rows.Scan(&location.ID, &location.Code, &location.Name,
			&location.Address, &location.IsActive, &location.CreatedAt, &location.UpdatedAt)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (r *venueRepository) GetLocationByID(id string) (*models.Location, error) {
	query := `
		SELECT id, code, name, address, is_active, created_at, updated_at
		FROM bar_system.locations 
		WHERE id = $1
	`

	location := &models.Location{}
	err := r.db.QueryRow(query, id).Scan(&location.ID, &location.Code, &location.Name,
		&location.Address, &location.IsActive, &location.CreatedAt, &location.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return location, nil
}

func (r *venueRepository) UpdateLocation(id string, location *models.Location) error {
	query := `
		UPDATE bar_system.locations 
		SET code = $1, name = $2, address = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	_, err := r.db.Exec(query, location.Code, location.Name, location.Address, location.IsActive, id)
	return err
}

func (r *venueRepository) DeleteLocation(id string) error {
	query := `
		UPDATE bar_system.locations 
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	return err
}

// Table operations
func (r *venueRepository) CreateTable(table *models.Table) error {
	query := `
		INSERT INTO bar_system.tables (id, location_id, code, seats, status, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (location_id, code) DO NOTHING
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(query, table.ID, table.LocationID, table.Code,
		table.Seats, table.Status, table.IsActive).Scan(&table.CreatedAt, &table.UpdatedAt)

	return err
}

func (r *venueRepository) GetTablesByLocation(locationID string) ([]models.Table, error) {
	query := `
		SELECT id, location_id, code, seats, status, is_active, created_at, updated_at
		FROM bar_system.tables 
		WHERE location_id = $1 AND is_active = true
		ORDER BY code
	`

	rows, err := r.db.Query(query, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make([]models.Table, 0)
	for rows.Next() {
		var table models.Table
		err := rows.Scan(&table.ID, &table.LocationID, &table.Code,
			&table.Seats, &table.Status, &table.IsActive, &table.CreatedAt, &table.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (r *venueRepository) GetTableByID(id string) (*models.Table, error) {
	query := `
		SELECT id, location_id, code, seats, status, is_active, created_at, updated_at
		FROM bar_system.tables 
		WHERE id = $1
	`

	table := &models.Table{}
	err := r.db.QueryRow(query, id).Scan(&table.ID, &table.LocationID, &table.Code,
		&table.Seats, &table.Status, &table.IsActive, &table.CreatedAt, &table.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return table, nil
}

func (r *venueRepository) UpdateTable(id string, table *models.Table) error {
	query := `
		UPDATE bar_system.tables 
		SET code = $1, seats = $2, status = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	_, err := r.db.Exec(query, table.Code, table.Seats, table.Status, table.IsActive, id)
	return err
}

func (r *venueRepository) DeleteTable(id string) error {
	query := `
		UPDATE bar_system.tables 
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	return err
}
