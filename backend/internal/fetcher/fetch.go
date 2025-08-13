// internal/fetcher/fetch.go (Updated Version)

package fetcher

import (
	// ... keep all your existing imports: bytes, encoding/json, fmt, etc.
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// authTokenResponse, SatelliteImage, and Fetcher struct remain exactly the same.
type authTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type SatelliteImage struct {
	ID         string
	AcquiredAt time.Time
	ImageData  image.Image
}
type Fetcher struct {
	client       *http.Client
	clientID     string
	clientSecret string
	token        string
	tokenExpiry  time.Time
	mu           sync.Mutex
}

// NewFetcher, getAccessToken, and ensureValidToken remain exactly the same.
func NewFetcher() (*Fetcher, error) {
	clientID := os.Getenv("SENTINELHUB_CLIENT_ID")
	clientSecret := os.Getenv("SENTINELHUB_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("environment variables SENTINELHUB_CLIENT_ID and SENTINELHUB_CLIENT_SECRET must be set")
	}
	return &Fetcher{
		client:       &http.Client{Timeout: time.Minute * 2},
		clientID:     clientID,
		clientSecret: clientSecret,
	}, nil
}
func (f *Fetcher) getAccessToken() error {
	// ... this function's code does not change
	tokenURL := "https://services.sentinel-hub.com/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", f.clientID)
	data.Set("client_secret", f.clientSecret)
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil { return fmt.Errorf("failed to create token request: %w", err) }
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := f.client.Do(req)
	if err != nil { return fmt.Errorf("failed to execute token request: %w", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get token, status: %s, body: %s", resp.Status, string(bodyBytes))
	}
	var tokenResp authTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}
	f.token = tokenResp.AccessToken
	f.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)
	fmt.Println("INFO: Successfully fetched new Sentinel Hub access token.")
	return nil
}
func (f *Fetcher) ensureValidToken() error {
	// ... this function's code does not change
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.token == "" || time.Now().After(f.tokenExpiry) {
		fmt.Println("INFO: Access token is expired or missing. Fetching a new one.")
		return f.getAccessToken()
	}
	return nil
}

// --- THIS IS THE MODIFIED FUNCTION ---
// It now accepts a `bbox` (bounding box) slice of float64.
func (f *Fetcher) FetchImageForLocation(bbox []float64, date time.Time) (*SatelliteImage, error) {
	if err := f.ensureValidToken(); err != nil {
		return nil, err
	}

	fmt.Printf("INFO: Fetching image from Sentinel Hub for bbox %v on %s\n", bbox, date.Format("2006-01-02"))

	requestURL := "https://services.sentinel-hub.com/api/v1/process"
	evalscript := `
		//VERSION=3
		function setup() {
			return { input: ["B04", "B03", "B02"], output: { bands: 3 } };
		}
		function evaluatePixel(sample) {
			return [2.5 * sample.B04, 2.5 * sample.B03, 2.5 * sample.B02];
		}`

	// The request body now uses the 'bbox' passed into the function.
	requestBody, err := json.Marshal(map[string]interface{}{
		"input": map[string]interface{}{
			"bounds": map[string]interface{}{
				"bbox": bbox, // Use the dynamic bbox here
				"properties": map[string]interface{}{
					"crs": "http://www.opengis.net/def/crs/OGC/1.3/CRS84",
				},
			},
			"data": []map[string]interface{}{
				{
					"type": "sentinel-2-l2a",
					"dataFilter": map[string]interface{}{
						"timeRange": map[string]string{
							"from": date.UTC().Format(time.RFC3339),
							"to":   date.UTC().Add(24 * time.Hour).Format(time.RFC3339),
						},
					},
				},
			},
		},
		"output": map[string]interface{}{
			"width":  512,
			"height": 512,
			"format": map[string]string{ "type": "image/png" },
		},
		"evalscript": evalscript,
	})
	// ...The rest of the function (creating and sending the request) is identical...

	if err != nil { return nil, fmt.Errorf("failed to marshal request body: %w", err) }
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(requestBody))
	if err != nil { return nil, fmt.Errorf("failed to create process request: %w", err) }
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+f.token)
	resp, err := f.client.Do(req)
	if err != nil { return nil, fmt.Errorf("failed to execute process request: %w", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch image, status: %s, body: %s", resp.Status, string(bodyBytes))
	}
	img, err := png.Decode(resp.Body)
	if err != nil { return nil, fmt.Errorf("failed to decode image: %w", err) }
	fmt.Printf("INFO: Successfully fetched and decoded image from Sentinel Hub.\n")

	result := &SatelliteImage{
		ID:         fmt.Sprintf("SH_IMG_BBOX%v_%d", bbox, date.Unix()),
		AcquiredAt: date,
		ImageData:  img,
	}
	return result, nil
}