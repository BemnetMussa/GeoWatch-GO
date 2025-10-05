// cmd/main.go (New API Server Version)

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"encoding/json"
	"geowatch-backend/internal/storage"
	"geowatch-backend/pkg/db"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// AppState holds the shared state for our application, like the fetcher instance.
type AppState struct {
	DB *pgxpool.Pool
}

type AnalysisRequest struct {
	AOI       [][][]float64 `json:"aoi"`
	StartDate string        `json:"startDate"`
	EndDate   string        `json:"endDate"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	dbPool, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("FATAL: Could not connect to the database: %v", err)
	}
	defer dbPool.Close()

	appState := &AppState{
		DB: dbPool,
	}

	router := gin.Default()

	// Your CORS setup is perfect. It allows our frontend on port 5173 to talk to this backend.
	// CORS from env, default to frontend port 8082
	allowedOrigin := os.Getenv("FRONTEND_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:8082"
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// --- 3. Define API Routes ---
	apiV1 := router.Group("/api/v1")
	{
		// Changed to POST to accept a JSON body
		apiV1.POST("/changes", appState.getChangesHandler)
		apiV1.GET("/events", appState.getEventsHandler)
	}

	// --- KEY CHANGE: RUN ON A DIFFERENT PORT ---
	goServerPort := ":8000"
	if p := os.Getenv("GO_SERVER_PORT"); p != "" {
		if !strings.HasPrefix(p, ":") {
			p = ":" + p
		}
		goServerPort = p
	}
	fmt.Printf("INFO: Starting GeoWatch Go API server on port %s\n", goServerPort)
	router.Run(goServerPort)
}

// getChangesHandler is the heart of our API. It processes the request and returns the image.
// backend/cmd/main.go

// In backend/cmd/main.go

func (app *AppState) getChangesHandler(c *gin.Context) {
	fmt.Println("Go Backend: Received request from frontend.")

	// We still need to read the JSON from the frontend
	var requestData AnalysisRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Marshal the received data to be sent to the Python service
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request for Python service"})
		return
	}

	// The Python service is running on port 8080
	pythonServiceURL := os.Getenv("PYTHON_SERVICE_URL") + "/analyze-changes"

	fmt.Println("Go Backend: Forwarding request to Python GEE service...")
	resp, err := http.Post(pythonServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("ERROR: Could not connect to Python service: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to communicate with the GEE analysis service"})
		return
	}
	defer resp.Body.Close()

	// Check if the Python service itself returned an error
	if resp.StatusCode != http.StatusOK {
		// Try to read the error body from Python to give a better message
		errorBody, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR: Python service returned status %d with body: %s", resp.StatusCode, string(errorBody))
		c.JSON(http.StatusBadGateway, gin.H{"error": "The GEE analysis service returned an error"})
		return
	}

	// --- MORE EFFICIENT PROXYING ---
	// We get the content type from Python's response
	contentType := resp.Header.Get("Content-Type")
	// And the content length
	contentLength := resp.ContentLength

	// We use Gin's Stream method to copy the response body directly
	// without buffering the whole image in memory.
	fmt.Println("Go Backend: Streaming PNG response back to client...")
	c.Stream(func(w io.Writer) bool {
		_, err := io.Copy(w, resp.Body)
		return err != nil
	})

	// Set the headers on our response to match Python's
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
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
//     [http://localhost:8081/api/v1/changes?lat=40.78&lon=-73.97](http://localhost:8080/api/v1/changes?lat=40.78&lon=-73.97)

// When you hit that URL, you will trigger the entire process. Your terminal will show all the log messages, and after a few moments, your browser will display the final `change_overlay.png` result. You have now successfully built and deployed your change detection service as a web API!

// You can try different locations by changing the `lat` and `lon` parameters, or the time window by adding a `&days=90` parameter to the URL.
