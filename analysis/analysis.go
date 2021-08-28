package analysis

import (
	"image"
	"errors"
)

// Return true if the two images are equal
func AreImgsEqual(img1 *image.RGBA, img2 *image.RGBA) bool {
	// Images divided into rectangles
	img1Div := make([]image.Image, rectNum)
	img2Div := make([]image.Image, rectNum)
	
	rectNum := squaresAmount()

	// Height and width of each rectangle
	width := img1.Rect.Max.X / int(rectNum)
	height := img1.Rect.Max.Y

	// Divide the images in rectNum rectangles, each of which is as high as the
	// image and as wide as imgWidth
	for i := 0; i < int(rectNum); i++ {
		// Area of the sub-image
		subArea := image.Rect(width*i, height*i, width*(i+1)-1, height*(i+1)-1)
		
		// Get sub-images and add them
		subImg1 := img1.SubImage(subArea)
		subImg2 := img1.SubImage(subArea)
		img1Div[i] = subImg1
		img2Div[i] = subImg2
	}
	
	areImgsEqual := true
	
	// Signal-only channel to stop goroutines whenever one of them finds a
	// difference
	stopSig := make(chan struct{})
	defer close(stopSig)
	
	// Error signal to comunicate that an error happened in a goroutine
	errSig := make(chan error)
	defer close(errSig)
	
	// Launch a number of goroutines equal to the number of rectangles
	for i := 0; i < int(rectNum); i++ {
		go checkRect(stopSig, errSig, img1Div[i], img2Div[i])
	}
	
	return areImgsEqual
}

// Return true if the board has been erased
func BoardIsErased(img *image.RGBA, erasedBoard *image.RGBA) bool {
	return AreImgsEqual(img, erasedBoard)
}

// Check rectangles (in the same position) of two images
func checkRect(stopCh chan struct{}, errCh chan error, rect1 image.Image, rect2 image.Image) {
	r1Bounds := rect1.Bounds()
	r2Bounds := rect2.Bounds()

	// Rectangles must have the same size and the same position
	if r1Bounds != r2Bounds {
		errCh <- errors.New("Rectangles differ by size and/or position")
	}

	isStop := false

	for x := r1Bounds.Min.X; x < r1Bounds.Max.X; x++ {
		for y := r1Bounds.Min.Y; y < r1Bounds.Max.Y; y++ {
			// Stop if a stop signal is received
			select {
			case <-stopCh:
				isStop = true
			case <-errCh:
				isStop = true
			default:
				if rect1.At(x, y) != rect2.At(x, y) {
					stopCh <- struct{}{}
				}
			}

			if isStop {
				// Break inner loop
				break
			}
		}

		if isStop {
			// Break outer loop
			break
		}
	}
}

// Return the number of squares that the image will be divided into
func squaresAmount() int {
	procs := runtime.GOMAXPROCS(0)
	
	return uint(procs*2)
}
