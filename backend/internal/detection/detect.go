package detection

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"geowatch-backend/internal/fetcher"
)

// Result contains the output of a change detection analysis.
type Result struct {
	// The image overlay showing changes.
	ChangeOverlay image.Image
	// A metric for the amount of change, e.g., number of changed pixels.
	ChangeSeverity int
}

// VisualChange detects differences between two images by comparing pixel colors.
// It returns a result struct containing the overlay and the severity.
func VisualChange(imageA, imageB *fetcher.SatelliteImage, threshold float64) (*Result, error) {
	fmt.Println("INFO: Starting visual change detection...")

	if imageA == nil || imageB == nil || imageA.ImageData == nil || imageB.ImageData == nil {
		return nil, fmt.Errorf("cannot compare nil images")
	}

	boundsA := imageA.ImageData.Bounds()
	if boundsA != imageB.ImageData.Bounds() {
		return nil, fmt.Errorf("image dimensions do not match")
	}

	diffImage := image.NewRGBA(boundsA)
	highlightColor := color.RGBA{R: 255, G: 0, B: 0, A: 180}
	changedPixels := 0

	for y := boundsA.Min.Y; y < boundsA.Max.Y; y++ {
		for x := boundsA.Min.X; x < boundsA.Max.X; x++ {
			colorA := imageA.ImageData.At(x, y)
			colorB := imageB.ImageData.At(x, y)
			distance := colorDistance(colorA, colorB)

			if distance > threshold {
				diffImage.Set(x, y, highlightColor)
				changedPixels++
			} else {
				diffImage.Set(x, y, color.Transparent)
			}
		}
	}

	result := &Result{
		ChangeOverlay:  diffImage,
		ChangeSeverity: changedPixels,
	}

	fmt.Printf("INFO: Change detection finished. Found %d changed pixels.\n", changedPixels)
	return result, nil
}

// colorDistance calculates the Euclidean distance between two colors.
func colorDistance(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	r1, g1, b1 = r1/257, g1/257, b1/257
	r2, g2, b2 = r2/257, g2/257, b2/257
	deltaR := float64(r1) - float64(r2)
	deltaG := float64(g1) - float64(g2)
	deltaB := float64(b1) - float64(b2)
	return math.Sqrt(deltaR*deltaR + deltaG*deltaG + deltaB*deltaB)
}


