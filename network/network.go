package network

import (
	"os"
	"net/http"
	"bytes"
	"image"
	"image/jpeg"
)

// Store image in case the remote host is (temporarily) offline. Make sure the
// file does not exist or it will be overwritten.
func storeImg(img image.Image, name string) error {
	// Check if the "offline" directory exists. If it doesn't, create it.
	_, statErr := os.Stat("offline")
	if os.IsNotExist(statErr) {
		os.Mkdir("offline", 0644)
	}

	// Create a new file
	file, fileErr := os.Create(name)
	if fileErr != nil {
		return fileErr
	}

	// Save the image in that file
	saveErr := jpeg.Encode(file, img, nil)
	if saveErr != nil {
		return saveErr
	}

	return nil
}


// Send all the images that have been previously stored and delete them once
// they got sent. If it's not possible to send them they will be kept stored.
func sendStored() error {
	// Open directory
	dir, err := os.Open("offline")
	if err != nil {
		return err
	}

	// Read all entries
	entries, err := dir.ReadDir(0)
	if err != nil {
		return err
	}

	for _, en := range entries {
		if !en.IsDir() {
			name := en.Name()

			if len(name) > 5 && (name[:4] == ".jpg" || name[:5] == ".jpeg") {
				path := "offline" + name

				file, err := os.Open(path)
				if err != nil {
					return err
				}

				image, err := jpeg.Decode(file)
				if err != nil {
					return err
				}

				SendImg(path, image)
			}
		}
	}

	return nil
}


// Send image to Discord bot
func SendImg(url string, img image.Image) error {
	// Encode the image into JPEG and send it into a buffer
	var buffer bytes.Buffer
	errJpeg := jpeg.Encode(&buffer, img, nil)

	if errJpeg != nil {
		return errJpeg
	}

	// Send the image through POST from the buffer
	_, errHTTP := http.Post(url, "image/jpeg", &buffer)

	return errHTTP
}
