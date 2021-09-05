package analysis

import (
	"fmt"
	"sync"
	"image"
	"context"
	"runtime"
)

// Return true if the board has been erased
func BoardIsErased(img *image.RGBA, erasedBoard *image.RGBA) bool {
	return AreImgsEqual(img, erasedBoard)
}

// Return true if the two images are equal
func AreImgsEqual(img1, img2 *image.RGBA) bool {
	rectNum := int(squaresAmount())

	img1Div, img2Div := divideImgs(img1, img2, uint(rectNum))

	areEqual := true

	// Buffered channel for the result of each goroutine
	resCh := make(chan bool, rectNum)
	defer close(resCh)

	var wg sync.WaitGroup

	// Context made to stop all goroutines at a certain condition
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Launch a number of goroutines equal to the number of rectangles
	for i := 0; i < rectNum; i++ {
		wg.Add(1)
		go func(img1, img2 image.Image) {
			defer wg.Done()
			checkRect(ctx, cancel, resCh, img1, img2)
		}(img1Div[i], img2Div[i])

	}

	// Read the results
	for i := 0; i < rectNum; i++ {
		r := <-resCh
		if !r {
			areEqual = false
			break
		}
	}

	// Wait for goroutines to stop
	wg.Wait()

	return areEqual
}

// Check rectangles (in the same position) of two images asynchronously
func checkRect(ctx context.Context, cancel func(), resCh chan bool, rect1, rect2 image.Image) {
	r1Bounds := rect1.Bounds()

	for x := r1Bounds.Min.X; x < r1Bounds.Max.X; x++ {
		for y := r1Bounds.Min.Y; y < r1Bounds.Max.Y; y++ {
			if ctx.Err() != nil {
				// Another goroutine said to stop working
				return
			}

			if rect1.At(x, y) != rect2.At(x, y) {
				// The two images are not equal
				cancel()
				resCh <- false
				return
			}
		}
	}

	resCh <- true
}

// Divide two images in a certain amount of parts
func divideImgs(img1, img2 *image.RGBA, parts uint) ([]image.Image, []image.Image) {
	// Images divided into rectangles
	img1Div := make([]image.Image, parts)
	img2Div := make([]image.Image, parts)

	// Height and width of each rectangle
	// Rectangulars have the same height as the image and their width is divided
	// by the number of parts
	width := img1.Rect.Max.X / int(parts)
	height := img1.Rect.Max.Y

	// Divide the images in rectNum rectangles, each of which is as high as the
	// image and as wide as imgWidth
	for i := 0; i < int(parts); i++ {
		// Area of the sub-image
		subArea := image.Rect(width*i, 0, width*(i+1), height)

		if i == int(parts) - 1 {
			// The last one has also the remaining part (division remainder)
			subArea.Max.X = img1.Rect.Max.X
		}

		// Get sub-images and add them
		subImg1 := img1.SubImage(subArea)
		subImg2 := img2.SubImage(subArea)
		img1Div[i] = subImg1
		img2Div[i] = subImg2
	}

	return img1Div, img2Div
}

// Return the number of squares that the image will be divided into
func squaresAmount() uint {
	procs := runtime.GOMAXPROCS(0)

	return uint(procs*2)
}
