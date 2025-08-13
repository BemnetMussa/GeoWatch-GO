// internal/api/handlers.go
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// Make sure your module name in go.mod is correct. Assuming 'geowatch'.
	"geowatch-backend/internal/storage"
)

// API holds the dependencies for the API handlers, such as the data store.
type API struct {
	Store *storage.Store
}

// NewAPI creates a new API struct with the given store.
func NewAPI(store *storage.Store) *API {
	return &API{Store: store}
}

// HealthCheckHandler provides a simple health check endpoint.
func (a *API) HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetLocationsHandler handles requests to retrieve locations.
func (a *API) GetLocationsHandler(c *gin.Context) {
	locations, err := a.Store.GetLocations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locations)
}

// CreateLocationInput defines the expected JSON input for creating a location.
type CreateLocationInput struct {
	Name string `json:"name" binding:"required"`
	WKT  string `json:"wkt"  binding:"required"`
}

// CreateLocationHandler handles requests to create a new location.
func (a *API) CreateLocationHandler(c *gin.Context) {
	var input CreateLocationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := a.Store.CreateLocation(c.Request.Context(), input.Name, input.WKT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}