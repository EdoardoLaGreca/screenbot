package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"image"
	"image/png"
	"runtime"
	"github.com/EdoardoLaGreca/screenbot/analysis"
	"github.com/EdoardoLaGreca/screenbot/network"
	"github.com/kbinani/screenshot"
)

func main() {
	// Check if URL is provided
	if len(os.Args[1:]) < 1 {
		fmt.Println("ERROR: No URL provided.\nUsage: ./screenbot <IP>:<port>")
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

	// Check if the "offline" directory exists. If it doesn't, create it.
	_, statErrOffline := os.Stat("offline")
	if os.IsNotExist(statErrOffline) {
		os.Mkdir("offline", 0755)
	}

	// Check if the "imgs" directory exists. If it doesn't, create it.
	_, statErrImgs := os.Stat("imgs")
	if os.IsNotExist(statErrImgs) {
		os.Mkdir("imgs", 0755)
	}

	fmt.Scanf("%d %d %d %d", &area.Min.X, &area.Min.Y, &area.Max.X, &area.Max.Y)

	var prevImg, erasedImg *image.RGBA

	erasedImg, err := getErasedBoard()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("File erased.png found!")

		if erasedImg.Bounds() != area {
			fmt.Println("The image of file erased.png has different size " +
				"(width and height) compared to the chosen coordinates, so it " +
				"will be ignored.")
			erasedImg = nil
		}
	}

	// Get the first image
	if erasedImg == nil {
		fmt.Println("\nAssuming that currently the board is empty...")
		for {
			erasedImg = getScreenshot(area)
			if erasedImg != nil {
				break
			}
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

		if currImg == nil {
			// Continue if an issue occurred while capturing the screenshot
			time.Sleep(1*time.Second)
			continue
		}

		log.Println("Screenshot captured")

		if !analysis.AreImgsEqualConv(prevImg, currImg) {
			if analysis.BoardIsErased(currImg, erasedImg) {
				// Name for the stored image
				timestamp := time.Now().Format(time.RFC822)
				imgName := timestamp + ".png"

				// Store all images
				storeImgErr(prevImg, "imgs/" + imgName)
				
				// Send image through a goroutine to avoid a blocking behaviour
				go func(image *image.RGBA, filename string) {
					log.Println("Sending image...")
					err := network.SendImg(url, image)

					if err != nil {
						log.Println("Unable to send the image, error:", err)

						// Save the image as local file and send it once there
						// will be a connection again
						storeImgErr(image, "offline/" + filename)
					} else {
						log.Println("Image sent!")

						// Also send the images stored that has not been sent yet
						network.SendStored()
					}
				}(prevImg, imgName)

			}

			prevImg = currImg
		}

		time.Sleep(1*time.Second)
	}
}

// Store image and print errors, if any
func storeImgErr(img *image.RGBA, filename string) {
	err := network.StoreImg(img, filename)
	if err != nil {
		log.Println("Unable to store image, error:", err)
	}
}

// Get erased board image
func getErasedBoard() (*image.RGBA, error) {
	_, err := os.Stat("erased.png")
	if err != nil {
		msg := "The image for the erased board (erased.png) was not found, " +
			"using the first screenshot as erased board."
		return nil, fmt.Errorf(msg)
	}

	file, err := os.Open("erased.png")

	// Decode the image into image.Image
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img.(*image.RGBA), nil
}
