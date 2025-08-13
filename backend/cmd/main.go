// cmd/main.go (New API Server Version)

package main

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"time"
	"strings"

	"geowatch-backend/internal/fetcher"
	"geowatch-backend/internal/detection"
	"geowatch-backend/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"geowatch-backend/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AppState holds the shared state for our application, like the fetcher instance.
type AppState struct {
	Fetcher *fetcher.Fetcher
	DB      *pgxpool.Pool // Add this line
}


func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}	
	dbPool, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("FATAL: Could not connect to the database: %v", err)
	}
	defer dbPool.Close() // Ensure the connection is closed when main exits
	// --- 1. Initialization ---
	// Create the app state, which includes initializing our fetcher.
	// We do this once when the server starts.
	f, err := fetcher.NewFetcher()
	if err != nil {
		log.Fatalf("FATAL: Failed to create fetcher: %v", err)
	}
	appState := &AppState{
		Fetcher: f,
		DB:      dbPool, // Add this line
	}
	// --- 2. Setup Gin Router ---
	router := gin.Default()

	// --- 3. Define API Routes ---
	// Group routes under /api/v1 for versioning.
	apiV1 := router.Group("/api/v1")
	{
		// Our main endpoint for getting change detection overlays.
		apiV1.GET("/changes", appState.getChangesHandler)

		apiV1.GET("/events", appState.getEventsHandler)
	}

	fmt.Println("INFO: Starting GeoWatch API server on port 8080...")
	// Start the server. It will run forever until stopped.
	router.Run(":8080")
}

// getChangesHandler is the heart of our API. It processes the request and returns the image.
func (app *AppState) getChangesHandler(c *gin.Context) {
	// --- 1. Parse and Validate Query Parameters ---
	// Example Request: /api/v1/changes?lat=40.78&lon=-73.97&days=30
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	daysStr := c.DefaultQuery("days", "30") // Default to 30 days if not provided

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lon, errLon := strconv.ParseFloat(lonStr, 64)
	if errLat != nil || errLon != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing 'lat' and 'lon' query parameters."})
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'days' parameter. Must be a positive integer."})
		return
	}

	// --- 2. Prepare for Fetching ---
	// Define the time range.
	dateAfter := time.Now()
	dateBefore := dateAfter.AddDate(0, 0, -days)

	// Create a bounding box around the requested lat/lon point.
	// This creates a square of about 2.5km x 2.5km. You can adjust the size.
	size := 0.025
	bbox := []float64{lon - size/2, lat - size/2, lon + size/2, lat + size/2}

	// --- 3. Run the Core Logic (Fetch and Compare) ---

	// In cmd/main.go, inside the getChangesHandler function...

// ... (code for parsing parameters and defining bbox is the same) ...

	// --- 3. Run the Core Logic (Fetch, Process, and Compare) ---
	fmt.Println("Handler: Fetching 'before' image...")
	imageBeforeRaw, err := app.Fetcher.FetchImageForLocation(bbox, dateBefore)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch 'before' image", "details": err.Error()})
		return
	}
	// ADD THIS PROCESSING STEP
	imageBefore, err := fetcher.ProcessImage(imageBeforeRaw, 10) // Increase brightness slightly
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process 'before' image", "details": err.Error()})
		return
	}


	fmt.Println("Handler: Fetching 'after' image...")
	imageAfterRaw, err := app.Fetcher.FetchImageForLocation(bbox, dateAfter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch 'after' image", "details": err.Error()})
		return
	}
	// ADD THIS PROCESSING STEP
	imageAfter, err := fetcher.ProcessImage(imageAfterRaw, 10) // Increase brightness slightly
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process 'after' image", "details": err.Error()})
		return
	}

	fmt.Println("Handler: Comparing images...")
    changeThreshold := 60.0
    detectionResult, err := detection.VisualChange(imageBefore, imageAfter, changeThreshold)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run change detection", "details": err.Error()})
        return
    }


	// --- 4. Calculate Severity and Save Event to DB ---
	
	newEvent := storage.ChangeEvent{
        LocationID:  1, // Still using placeholder
        EventType:   "visual_change",
        Description: fmt.Sprintf("Detected %d changed pixels between %s and %s.", detectionResult.ChangeSeverity, dateBefore.Format("2006-01-02"), dateAfter.Format("2006-01-02")),
        DetectedAt:  time.Now(),
        Severity:    detectionResult.ChangeSeverity,
    }

	

	_, err = storage.SaveChangeEvent(app.DB, newEvent, bbox)
	if err != nil {
		// Log the error but don't block the user. The main goal is to return the image.
		log.Printf("ERROR: Failed to save change event to database: %v", err)
	}


 // --- 6. Encode and Return the Result Image ---
    buffer := new(bytes.Buffer)
    // We get the overlay image from the detection result.
    if err := png.Encode(buffer, detectionResult.ChangeOverlay); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode result image."})
        return
    }

    c.Data(http.StatusOK, "image/png", buffer.Bytes())
}

// getEventsHandler queries and returns historical change events from the DB.
func (app *AppState) getEventsHandler(c *gin.Context) {
	// Example Request: /api/v1/events?bbox=-74.0,40.7,-73.9,40.8
	bboxStr := c.Query("bbox")
	if bboxStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'bbox' query parameter. Example format: 'minLon,minLat,maxLon,maxLat'"})
		return
	}

	// Parse the bbox string (e.g., "-74.0,40.7,-73.9,40.8") into a slice of float64
	coords, err := parseBbox(bboxStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'bbox' format.", "details": err.Error()})
		return
	}

	// Call our new loading function
	events, err := storage.LoadChangeEventsInBBox(app.DB, coords)
	if err != nil {
		log.Printf("ERROR: Failed to load events from DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query for events."})
		return
	}

	// Return the found events as a JSON array.
	c.JSON(http.StatusOK, events)
}

// We need a helper to parse the bounding box string.
// Add this helper function also in main.go.
func parseBbox(bboxStr string) ([]float64, error) {
	parts := strings.Split(bboxStr, ",")
	if len(parts) != 4 {
		return nil, fmt.Errorf("bbox must have exactly 4 parts: minLon,minLat,maxLon,maxLat")
	}

	var coords []float64
	for _, part := range parts {
		val, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number in bbox: %s", part)
		}
		coords = append(coords, val)
	}
	return coords, nil
}






// Action Plan

// 1.  **Run `go get github.com/gin-gonic/gin`** if you haven't already.
// 2.  **Replace** the code in `internal/fetcher/fetch.go` with the updated version from Step 2.
// 3.  **Replace** the code in `cmd/main.go` with the new API server code from Step 3.
// 4.  **Run the server:** Open your terminal in the project root and execute:
//     ```bash
//     go run ./cmd/main.go
//     ```
//     You should see the message: `INFO: Starting GeoWatch API server on port 8080...`
// 5.  **Test your live API!** Open your web browser and navigate to this URL:
//     [http://localhost:8080/api/v1/changes?lat=40.78&lon=-73.97](http://localhost:8080/api/v1/changes?lat=40.78&lon=-73.97)

// When you hit that URL, you will trigger the entire process. Your terminal will show all the log messages, and after a few moments, your browser will display the final `change_overlay.png` result. You have now successfully built and deployed your change detection service as a web API!

// You can try different locations by changing the `lat` and `lon` parameters, or the time window by adding a `&days=90` parameter to the URL.