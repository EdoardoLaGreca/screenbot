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
	entries, err := os.Open("offline").ReadDir(0)
	if err != nil {
		return err
	}

	for _, v := range entries {
		
	}
}


// Send image to Discord bot
func SendImg(url string, img *image.RGBA) error {
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
