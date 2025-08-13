package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// We can reuse the ChangeEvent struct from save.go, but we need to add fields
// that we want to read back from the database, like the ID and the geometry.
type ChangeEventWithGeom struct {
	ID          int       `json:"id"`
	LocationID  int       `json:"location_id"`
	EventType   string    `json:"event_type"`
	Description string    `json:"description"`
	DetectedAt  time.Time `json:"detected_at"`
	Severity    int       `json:"severity"`
	// We'll read the geometry as GeoJSON, which is very frontend-friendly.
	GeoJSON string `json:"geom_geojson"`
}

// LoadChangeEventsInBBox finds all change events whose saved geometry intersects
// with the provided bounding box.
func LoadChangeEventsInBBox(pool *pgxpool.Pool, bbox []float64) ([]ChangeEventWithGeom, error) {
	// This query finds events where the event's geometry (`geom`) intersects with
	// a new polygon (`query_geom`) that we create from the user's request bbox.
	// ST_AsGeoJSON converts the geometry into a JSON string, perfect for APIs.
	query := `
		SELECT id, location_id, event_type, description, detected_at, severity, ST_AsGeoJSON(geom)
		FROM change_events
		WHERE ST_Intersects(geom, ST_MakeEnvelope($1, $2, $3, $4, 4326))
		ORDER BY detected_at DESC;
	`

	rows, err := pool.Query(context.Background(), query, bbox[0], bbox[1], bbox[2], bbox[3])
	if err != nil {
		return nil, fmt.Errorf("failed to execute query to load change events: %w", err)
	}
	defer rows.Close()

	var events []ChangeEventWithGeom

	// We iterate through all the rows returned by the query.
	for rows.Next() {
		var event ChangeEventWithGeom
		// For each row, we scan the columns into the fields of our struct.
		if err := rows.Scan(
			&event.ID,
			&event.LocationID,
			&event.EventType,
			&event.Description,
			&event.DetectedAt,
			&event.Severity,
			&event.GeoJSON,
		); err != nil {
			// If one row fails, we log it and continue, so the user still gets partial results.
			fmt.Printf("WARNING: Failed to scan one change event row: %v\n", err)
			continue
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over event rows: %w", err)
	}
	
	fmt.Printf("INFO: Loaded %d change events from the database.\n", len(events))
	return events, nil
}	