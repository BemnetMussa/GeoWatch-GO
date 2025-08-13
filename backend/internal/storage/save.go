// internal/storage/save.go

package storage

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ChangeEvent represents the data we want to save to the database.
type ChangeEvent struct {
	LocationID  int       // Foreign key to the locations table
	EventType   string    // e.g., 'visual_change_detected'
	Description string    // A summary of the event
	DetectedAt  time.Time // When the analysis was run
	Severity    int       // A measure of the change, e.g., number of changed pixels
	// The GEOMETRY will be handled by PostGIS functions in the SQL query
}

// SaveChangeEvent saves a detected change event to the database.
// It takes the database connection pool and the event details.
// bbox is the bounding box used for the analysis, which we'll save as the event's geometry.
func SaveChangeEvent(pool *pgxpool.Pool, event ChangeEvent, bbox []float64) (int, error) {
	// The SQL query to insert a new record into the change_events table.
	// We use ST_MakeEnvelope to create a PostGIS polygon geometry from the bounding box.
	// ST_SetSRID sets the spatial reference system (4326 is standard WGS84 lat/lon).
	// `RETURNING id` gives us back the ID of the newly created row.
	query := `
		INSERT INTO change_events (location_id, event_type, description, detected_at, severity, geom)
		VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakeEnvelope($6, $7, $8, $9), 4326))
		RETURNING id;
	`

	var eventID int
	// We use QueryRow because we expect exactly one row to be returned (the new ID).
	err := pool.QueryRow(
		context.Background(),
		query,
		event.LocationID,
		event.EventType,
		event.Description,
		event.DetectedAt,
		event.Severity,
		bbox[0], // min Longitude
		bbox[1], // min Latitude
		bbox[2], // max Longitude
		bbox[3], // max Latitude
	).Scan(&eventID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert change event into database: %w", err)
	}

	fmt.Printf("INFO: Successfully saved change event with ID %d to the database.\n", eventID)
	return eventID, nil
}

// CountChangedPixels is a helper function to calculate a simple severity metric.
// It counts the number of non-transparent pixels in the difference image.
func CountChangedPixels(diffImage image.Image) int {
	bounds := diffImage.Bounds()
	changedPixels := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the alpha channel of the pixel.
			// Our comparison function makes unchanged pixels fully transparent (alpha=0).
			_, _, _, a := diffImage.At(x, y).RGBA()
			if a > 0 {
				changedPixels++
			}
		}
	}
	return changedPixels
}