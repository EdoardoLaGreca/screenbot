package network

import (
	"os"
	"io"
	"fmt"
	"net"
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

				img, err := png.Decode(file)
				if err != nil {
					return err
				}

				err = SendImg(path, img.(*image.RGBA))
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
func SendImg(url string, img *image.RGBA) error {
	// Encode the image into JPEG and send it into a buffer
	var buffer bytes.Buffer
	errJpeg := png.Encode(&buffer, img)

	if errJpeg != nil {
		return errJpeg
	}

	conn, err := net.Dial("tcp", url)
	if err != nil {
		return err
	}

	// Secret word to check packets
	fmt.Fprintf(conn, "BOT")

	// Image length as uint64 and encoded into bytes (big-endian)
	length := uint64(len(img.Pix))
	var lenEnc [8]byte

	// Translate the image length into an array of 8 bytes, starting from the
	// most significant byte
	for i := 0; i < 8; i++ {
		// i-th byte counting from right
		b := byte((length & (0xFF << (i*8))) >> (i*8))

		lenEnc[8-i-1] = b
	}
	
	// Send length
	fmt.Fprintf(conn, "%s", string(lenEnc[:]))

	io.Copy(conn, &buffer)

	return nil
}
