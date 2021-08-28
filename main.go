package main

import (
	"fmt"
	"log"
	"image"
	"image/jpeg"
	"bytes"
	"net/http"
	"github.com/kbinani/screenshot"
)

func main() {
	var screenArea image.Rectangle
	var prevImg, erasedImg *image.RGBA
	remoteUrl := "" // Fill

	fmt.Print("Enter the coordinates of the rectangle to observe as\n" +
		" <x1> <y1> <x2> <y2>\n" +
		"where x1 and y1 are the coordinates of the top-left corner and x2" +
		"and y2 are the coordinates of the bottom-right corner: ")
	fmt.Scanf("%d %d %d %d", &screenArea.Max.X, &screenArea.Max.Y,
		&screenArea.Min.X, &screenArea.Min.Y)

	fmt.Println("\nAssuming that currently the board is empty...")

	firstImg := true

	for {
		currImg, err := screenshot.CaptureRect(screenArea)

		if err != nil {
			log.Println("Unable to capture the screenshot:", err)
			continue
		}

		log.Println("Screenshot captured")

		// Skip analysis if the current image is the first image
		if firstImg {
			prevImg = currImg
			erasedImg = currImg
			firstImg = false
			continue
		}

		if !imgsAreEqual(prevImg, currImg) {
			if boardIsErased(currImg, erasedImg) {
				err := sendImg(remoteUrl, prevImg)

				if err != nil {
					log.Println("Unable to send the image:", err)
				} else {
					log.Println("Image sent!")
				}
			}

			prevImg = currImg
		}
	}
}

// Return true if the two images are equal
func imgsAreEqual(img1 *image.RGBA, img2 *image.RGBA, rectNum uint) bool {
	// Images divided into rectangles
	img1Div := make([]*image.RGBA, rectNum)
	img2Div := make([]*image.RGBA, rectNum)

	// Height and width of each rectangle
	width := uint(img1.Rect.Max.X / rectNum)
	height := uint(img1.Rect.Max.Y)

	// Divide the images in rectNum rectangles, each of which is as high as the
	// image and as wide as imgWidth
	for i = 0; i < rectNum; i++ {
		// Area of the sub-image
		subArea := image.Rect(width*i, height*i, width*(i+1)-1, height*(i+1)-1)
		
		// Get sub-images and add them
		subImg1 = img1.SubImage(subArea)
		subImg2 = img1.SubImage(subArea
		img1Div[i] = subImg1
		img2Div[i] = subImg2
	}
	
	areImgsEqual = true
	
	// Launch a number of goroutines equal to the number of rectangles
	for i = 0; i < rectNum; i++ {
		go
	}
	
	return areImagesEqual
}

// Return true if the board has been erased
func boardIsErased(img *image.RGBA, erasedBoard *image.RGBA) bool {
	return imgsAreEqual(img, erasedBoard)
}

// Send image to Discord bot
func sendImg(url string, img *image.RGBA) error {
	// Encode the image into JPEG and send it into a buffer
	var buffer bytes.Buffer
	err := jpeg.Encode(&buffer, img, nil)

	if err != nil {
		return err
	}

	// Send the image through POST from the buffer
	risp, err := http.Post(url, "image/jpeg", &buffer)

	return nil
}
