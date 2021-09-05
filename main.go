package main

import (
	"os"
	"fmt"
	"log"
	"time"
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
		fmt.Println("ERROR: No URL provided.\nUsage: ./screenbot <URL>")
		return
	}
	
	// URL to send the images to
	remoteUrl := os.Args[1]

	runtime.GOMAXPROCS(8) // Tune this value

	var screenArea image.Rectangle
	var prevImg, erasedImg *image.RGBA

	fmt.Print("Enter the coordinates of the rectangle to observe as\n" +
		" <x1> <y1> <x2> <y2>\n" +
		"where (x1, y1) are the coordinates of the top-left corner and " +
		"(x2, y2) are the coordinates of the bottom-right corner: ")

	fmt.Scanf("%d %d %d %d", &screenArea.Min.X, &screenArea.Min.Y,
		&screenArea.Max.X, &screenArea.Max.Y)

	fmt.Println("\nAssuming that currently the board is empty...")

	firstImg := true

	for {
		currImg, err := screenshot.CaptureRect(screenArea)

		if err != nil {
			log.Println("Unable to capture the screenshot, error:", err)
			continue
		}

		log.Println("Screenshot captured")

		// Skip analysis if the current image is the first image
		if firstImg {
			prevImg = currImg
			erasedImg = currImg
			firstImg = false
		} else if !analysis.AreImgsEqual(prevImg, currImg) {
			if analysis.BoardIsErased(currImg, erasedImg) {
				err := sendImg(remoteUrl, prevImg)

				if err != nil {
					log.Println("Unable to send the image, error:", err)
					// Save the image as local file and send it once there
					// will be a connection again
				} else {
					log.Println("Image sent!")
				}
			}

			prevImg = currImg
		}

		time.Sleep(1*time.Second)
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
	_, err = http.Post(url, "image/jpeg", &buffer)

	return err
}
