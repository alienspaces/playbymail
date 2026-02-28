package scanner

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // Register PNG decoder

	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

const (
	maxImageDimension = 2000            // Max width or height in pixels
	jpegQuality       = 85              // JPEG compression quality (0-100)
	maxImageSizeBytes = 5 * 1024 * 1024 // 5MB max after optimization
)

// optimizeImageForOCR resizes and compresses images to reduce payload size
// for OpenAI API calls while maintaining OCR readability.
func optimizeImageForOCR(l logger.Logger, imageData []byte, originalMIME string) ([]byte, string, error) {
	l = l.WithFunctionContext("ImageOptimizer/optimizeImageForOCR")

	optStart := len(imageData)

	l.Info("optimizing image for OCR original_size=%d mime=%s", optStart, originalMIME)

	// Decode image
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	l.Info("decoded image dimensions=%dx%d format=%s", width, height, format)

	// Resize if needed
	if width > maxImageDimension || height > maxImageDimension {
		scale := float64(maxImageDimension) / float64(max(width, height))
		newWidth := int(float64(width) * scale)
		newHeight := int(float64(height) * scale)

		l.Info("resizing image from %dx%d to %dx%d scale=%.2f", width, height, newWidth, newHeight, scale)

		img = resizeImage(img, newWidth, newHeight)
		bounds = img.Bounds()
		width = bounds.Dx()
		height = bounds.Dy()
	}

	// Encode as JPEG (smaller than PNG for photos/scans)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, "", fmt.Errorf("failed to encode JPEG: %w", err)
	}

	optimized := buf.Bytes()
	optEnd := len(optimized)

	// Validate JPEG signature (starts with FF D8 FF)
	if len(optimized) < 3 || optimized[0] != 0xFF || optimized[1] != 0xD8 || optimized[2] != 0xFF {
		l.Warn("optimized JPEG data appears invalid (missing JPEG signature), using original")
		return imageData, originalMIME, nil
	}

	// Verify the optimized JPEG can be decoded (validates it's a real image)
	if _, _, err := image.Decode(bytes.NewReader(optimized)); err != nil {
		l.Warn("optimized JPEG data cannot be decoded, using original error=%v", err)
		return imageData, originalMIME, nil
	}

	reduction := float64(optStart-optEnd) / float64(optStart) * 100

	// If optimization made the image larger or it's already small, return original
	if optEnd >= optStart || optStart < 100*1024 { // Skip if larger or already < 100KB
		l.Debug("skipping optimization original_size=%d optimized_size=%d (would be larger or already small)", optStart, optEnd)
		return imageData, originalMIME, nil
	}

	l.Info("optimized image original_size=%d optimized_size=%d reduction=%.1f%% dimensions=%dx%d",
		optStart, optEnd, reduction, width, height)

	if optEnd > maxImageSizeBytes {
		l.Warn("optimized image still large size=%d max=%d", optEnd, maxImageSizeBytes)
	}

	return optimized, "image/jpeg", nil
}

// resizeImage resizes an image using nearest-neighbor interpolation
func resizeImage(src image.Image, width, height int) image.Image {
	srcBounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	srcW := float64(srcBounds.Dx())
	srcH := float64(srcBounds.Dy())
	dstW := float64(width)
	dstH := float64(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := int((float64(x) + 0.5) * srcW / dstW)
			srcY := int((float64(y) + 0.5) * srcH / dstH)

			if srcX >= srcBounds.Dx() {
				srcX = srcBounds.Dx() - 1
			}
			if srcY >= srcBounds.Dy() {
				srcY = srcBounds.Dy() - 1
			}

			dst.Set(x, y, src.At(srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY))
		}
	}

	return dst
}

