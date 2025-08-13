// The purpose of internal/fetcher/process.go is to handle any pre-processing or cleaning of the images after they are downloaded but before they are compared. This is a common step in real-world image analysis pipelines.
// Why is this useful?
// Normalization: Satellite images can have different brightness or contrast due to the time of day or atmospheric conditions. A processing step can help normalize this.
// Filtering: We could apply filters to reduce "noise" in the image, which might lead to a cleaner comparison and fewer false positives.
// Cloud Masking (Advanced): A very common problem is clouds obscuring the view. An advanced processing step could try to detect clouds and exclude those areas from the comparison.


package fetcher

import (
	"fmt"
	"image"
	"image/color"
)

// ProcessImage applies adjustments to a raw satellite image to prepare it for comparison.
// For this implementation, we'll add a simple brightness adjustment.
//
// 'brightnessChange' can be positive (to brighten) or negative (to darken). A value
// between -100 and 100 is reasonable.
func ProcessImage(satImage *SatelliteImage, brightnessChange int) (*SatelliteImage, error) {
	if satImage == nil || satImage.ImageData == nil {
		return nil, fmt.Errorf("cannot process a nil image")
	}

	fmt.Printf("INFO: Processing image %s with brightness adjustment: %d\n", satImage.ID, brightnessChange)

	// Get the original image's properties.
	bounds := satImage.ImageData.Bounds()
	processedImg := image.NewRGBA(bounds) // Create a new image to hold the result.

	// Loop over every pixel of the original image.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := satImage.ImageData.At(x, y)
			r, g, b, a := originalColor.RGBA()

			// Convert 16-bit color channels to 8-bit (0-255).
			// We work with 8-bit because it's simpler for brightness adjustments.
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// Apply the brightness adjustment to each color channel.
			newR := adjustBrightness(r8, brightnessChange)
			newG := adjustBrightness(g8, brightnessChange)
			newB := adjustBrightness(b8, brightnessChange)

			// Set the new, adjusted pixel color in our output image.
			processedImg.Set(x, y, color.RGBA{R: newR, G: newG, B: newB, A: a8})
		}
	}

	// Return a new SatelliteImage struct with the processed image data.
	return &SatelliteImage{
		ID:         satImage.ID + "_processed",
		AcquiredAt: satImage.AcquiredAt,
		ImageData:  processedImg,
	}, nil
}

// adjustBrightness adds the change value to a single color channel,
// clamping the result to the valid 0-255 range.
func adjustBrightness(channelVal uint8, change int) uint8 {
	newVal := int(channelVal) + change
	if newVal < 0 {
		return 0
	}
	if newVal > 255 {
		return 255
	}
	return uint8(newVal)
}