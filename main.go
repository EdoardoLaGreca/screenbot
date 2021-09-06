package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"image"
	"runtime"
	"github.com/EdoardoLaGreca/screenbot/analysis"
	"github.com/EdoardoLaGreca/screenbot/network"
	"github.com/kbinani/screenshot"
)

func main() {
	// Check if URL is provided
	if len(os.Args[1:]) < 1 {
		fmt.Println("ERROR: No URL provided.\nUsage: ./screenbot <URL>")
		return
	}

	runtime.GOMAXPROCS(8) // Tune this value if needed
	
	// URL to send the images to
	remoteUrl := os.Args[1]

	var area image.Rectangle

	fmt.Print("Enter the coordinates of the rectangle to observe as\n" +
		" <x1> <y1> <x2> <y2>\n" +
		"where (x1, y1) are the coordinates of the top-left corner and " +
		"(x2, y2) are the coordinates of the bottom-right corner: ")

	fmt.Scanf("%d %d %d %d", &area.Min.X, &area.Min.Y, &area.Max.X, &area.Max.Y)

	fmt.Println("\nAssuming that currently the board is empty...")

	var prevImg, erasedImg *image.RGBA

	// Get the first image
	for {
		erasedImg = getScreenshot(area)
		if erasedImg != nil {
			break
		}
	}

	prevImg = erasedImg

	sendLoop(area, prevImg, erasedImg, remoteUrl)
}

// Get the screenshot and print an error message if it failed
func getScreenshot(area image.Rectangle) *image.RGBA {
	img, err := screenshot.CaptureRect(area)

	if err != nil {
		log.Println("Unable to capture the screenshot, error:", err)
		return nil
	}

	return img
}

func sendLoop(area image.Rectangle, prevImg, erasedImg *image.RGBA, url string) {

	for {
		currImg := getScreenshot(area)

		log.Println("Screenshot captured")

		if !analysis.AreImgsEqual(prevImg, currImg) {
			if analysis.BoardIsErased(currImg, erasedImg) {

				// Send image through a goroutine to avoid a blocking behaviour
				go func() {
					err := network.SendImg(url, prevImg)

					if err != nil {
						log.Println("Unable to send the image, error:", err)

						// Save the image as local file and send it once there
						// will be a connection again

					} else {
						log.Println("Image sent!")
					}
				}()

			}

			prevImg = currImg
		}

		time.Sleep(1*time.Second)
	}
}
