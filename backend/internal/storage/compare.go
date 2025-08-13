// // internal/storage/compare.go
package storage
// package storage

// import (
// 	"fmt"
// 	"image"
// 	"image/color"
// 	"math"

// 	// You'll need to replace 'geowatch-backend' with your actual Go module name
// 	// defined in your go.mod file.
// 	"geowatch-backend/internal/fetcher"
// )

// // CompareImages takes two satellite images and a sensitivity threshold, then returns
// // a new image highlighting the differences.
// // The threshold determines how sensitive the change detection is. A lower value means
// // more minor changes will be detected. A good starting point is often around 50-75.
// func CompareImages(imageA, imageB *fetcher.SatelliteImage, threshold float64) (image.Image, error) {
// 	fmt.Println("INFO: Starting image comparison...")

// 	// --- 1. Basic Validation ---
// 	if imageA == nil || imageB == nil || imageA.ImageData == nil || imageB.ImageData == nil {
// 		return nil, fmt.Errorf("cannot compare nil images")
// 	}

// 	boundsA := imageA.ImageData.Bounds()
// 	boundsB := imageB.ImageData.Bounds()

// 	if boundsA != boundsB {
// 		return nil, fmt.Errorf("image dimensions do not match: A=%v, B=%v", boundsA, boundsB)
// 	}

// 	// --- 2. Prepare the Output Image ---
// 	// We create a new RGBA image to store our results. It's the same size as the inputs.
// 	diffImage := image.NewRGBA(boundsA)

// 	// Define the color we will use to highlight changes.
// 	// A semi-transparent red is great for overlays.
// 	highlightColor := color.RGBA{R: 255, G: 0, B: 0, A: 180} // Bright red, 70% opaque

// 	// --- 3. Pixel-by-Pixel Comparison ---
// 	// We loop over every pixel in the image.
// 	for y := boundsA.Min.Y; y < boundsA.Max.Y; y++ {
// 		for x := boundsA.Min.X; x < boundsA.Max.X; x++ {
// 			// Get the color of the pixel from both images at the same (x, y) coordinate.
// 			colorA := imageA.ImageData.At(x, y)
// 			colorB := imageB.ImageData.At(x, y)

// 			// Calculate the "distance" between the two colors.
// 			distance := colorDistance(colorA, colorB)

// 			// If the distance is greater than our threshold, a significant change has occurred.
// 			if distance > threshold {
// 				// We color this pixel in our output image with the highlight color.
// 				diffImage.Set(x, y, highlightColor)
// 			} else {
// 				// Otherwise, there was no significant change, so we make this pixel
// 				// completely transparent. This is KEY for creating an overlay.
// 				diffImage.Set(x, y, color.Transparent)
// 			}
// 		}
// 	}

// 	fmt.Println("INFO: Image comparison finished successfully.")
// 	return diffImage, nil
// }

// // colorDistance calculates the Euclidean distance between two colors in RGB space.
// // This gives us a single number representing how "different" the two colors are.
// func colorDistance(c1, c2 color.Color) float64 {
// 	// The `At()` method returns a generic `color.Color`. We convert it to RGBA
// 	// so we can access the individual Red, Green, Blue, and Alpha channels.
// 	r1, g1, b1, _ := c1.RGBA()
// 	r2, g2, b2, _ := c2.RGBA()

// 	// The values returned by RGBA() are 16-bit (0-65535). We scale them down to
// 	// the more common 8-bit range (0-255) for our calculation.
// 	r1, g1, b1 = r1/257, g1/257, b1/257
// 	r2, g2, b2 = r2/257, g2/257, b2/257

// 	// The Euclidean distance formula: sqrt((r2-r1)^2 + (g2-g1)^2 + (b2-b1)^2)
// 	deltaR := float64(r1) - float64(r2)
// 	deltaG := float64(g1) - float64(g2)
// 	deltaB := float64(b1) - float64(b2)

// 	return math.Sqrt(deltaR*deltaR + deltaG*deltaG + deltaB*deltaB)
// }