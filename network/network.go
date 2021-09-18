package network

import (
	"os"
	"net/http"
	"bytes"
	"image"
	"image/png"
)

// Store image in case the remote host is (temporarily) offline. Make sure the
// file does not exist or it will be overwritten.
func StoreImg(img image.Image, name string) error {
	// Create a new file
	file, fileErr := os.Create(name)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	// Save the image in that file
	saveErr := png.Encode(file, img)
	if saveErr != nil {
		return saveErr
	}

	return nil
}


// Send all the images that have been previously stored and delete them once
// they got sent. If it's not possible to send them they will be kept stored.
func SendStored() error {
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

			if len(name) > 4 && name[:4] == ".png" {
				path := "offline/" + name

				file, err := os.Open(path)
				if err != nil {
					return err
				}

				image, err := png.Decode(file)
				if err != nil {
					return err
				}

				err = SendImg(path, image)
				if err != nil {
					// The connection has been lost again, stop sending images
					break
				} else {
					// Remove the image that has been sent
					err = os.Remove(path)
				}
			}
		}
	}

	return nil
}


// Send image to Discord bot
func SendImg(url string, img image.Image) error {
	// Encode the image into JPEG and send it into a buffer
	var buffer bytes.Buffer
	errJpeg := png.Encode(&buffer, img)

	if errJpeg != nil {
		return errJpeg
	}

	// Send the image through POST from the buffer
	_, errHTTP := http.Post(url, "image/png", &buffer)

	return errHTTP
}
