// internal/storage/storage.go
package storage

import (
	"context"
	"database/sql"
	// The blank import is used for the side-effect of registering the pgx driver.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Store will hold the database connection pool.
type Store struct {
	db *sql.DB
}

// NewStore creates and returns a new Store with a database connection.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Location represents the structure of a location in the database.
type Location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetLocations retrieves a list of locations from the database.
func (s *Store) GetLocations(ctx context.Context) ([]Location, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name FROM locations LIMIT 10")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var loc Location
		if err := rows.Scan(&loc.ID, &loc.Name); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}

	return locations, nil
}

// CreateLocation inserts a new location into the database using its name and WKT geometry.
func (s *Store) CreateLocation(ctx context.Context, name, wkt string) (int, error) {
	query := `INSERT INTO locations (name, geom) VALUES ($1, ST_GeomFromText($2)) RETURNING id`
	var id int
	err := s.db.QueryRowContext(ctx, query, name, wkt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}