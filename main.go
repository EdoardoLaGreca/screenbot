package main

import (
	"os"
	"fmt"
	"log"
	"bytes"
	"image"
	"image/jpeg"
	"runtime"
	"net/http"
	"github.com/EdoardoLaGreca/screenbot/analysis"
	"github.com/kbinani/screenshot"
)

func main() {
	// Check if URL is provided
	if len(os.Args[1:]) < 1 {
		fmt.Println("ERROR: No URL provided.\nUsage: screenbot <URL>")
		return
	}
	
	var screenArea image.Rectangle
	var prevImg, erasedImg *image.RGBA
	runtime.GOMAXPROCS(8) // Tune this value
	
	// URL to send the images to
	remoteUrl := os.Args[1]
	
	// Amount of squares to divide the images into
	squares := 8

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

		if !analysis.AreImgsEqual(prevImg, currImg) {
			if analysis.BoardIsErased(currImg, erasedImg) {
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
