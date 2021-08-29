package analysis

import (
//	"fmt"
	"sync"
	"image"
	"runtime"
)

// Return true if the board has been erased
func BoardIsErased(img *image.RGBA, erasedBoard *image.RGBA) bool {
	return AreImgsEqual(img, erasedBoard)
}

// Return true if the two images are equal
func AreImgsEqual(img1, img2 *image.RGBA) bool {
	rectNum := squaresAmount()

	// Images divided into rectangles
	img1Div := make([]image.Image, rectNum)
	img2Div := make([]image.Image, rectNum)

	// Height and width of each rectangle
	width := img1.Rect.Max.X / int(rectNum)
	height := img1.Rect.Max.Y

	// Divide the images in rectNum rectangles, each of which is as high as the
	// image and as wide as imgWidth
	for i := 0; i < int(rectNum); i++ {
		// Area of the sub-image
		subArea := image.Rect((width*i)+1, 0, width*(i+1), height)

		if i == rectNum - 1 {
			// The last one has also the remaining part (division remainder)
			subArea.Max.X = img1.Rect.Max.X
		}
		
		// Get sub-images and add them
		subImg1 := img1.SubImage(subArea)
		subImg2 := img1.SubImage(subArea)
		img1Div[i] = subImg1
		img2Div[i] = subImg2
	}

	areEqual := true

	// Signal-only channel to stop goroutines whenever one of them finds a
	// difference
	stopSig := make(chan struct{})
	defer close(stopSig)

	var wg sync.WaitGroup

	// Launch a number of goroutines equal to the number of rectangles
	for i := 0; i < int(rectNum); i++ {
		wg.Add(1)
		go checkRect(stopSig, &wg, img1Div[i], img2Div[i])
	}
	
	// Get the result
	// (some code here...)

	// Wait for goroutines to stop
	wg.Wait()

	return areEqual
}

// Check rectangles (in the same position) of two images
func checkRect(stopCh chan struct{}, wg *sync.WaitGroup, rect1, rect2 image.Image) {
	defer wg.Done()
	
	r1Bounds := rect1.Bounds()

	for x := r1Bounds.Min.X; x < r1Bounds.Max.X; x++ {
		for y := r1Bounds.Min.Y; y < r1Bounds.Max.Y; y++ {
			// Stop if a stop signal is received
			select {
			case <-stopCh:
				return
			default:
			}

			if rect1.At(x, y) != rect2.At(x, y) {
				stopCh <- struct{}{}
				return
			}
		}
	}
}

// Return the number of squares that the image will be divided into
func squaresAmount() int {
	procs := runtime.GOMAXPROCS(0)

	return procs*2
}
