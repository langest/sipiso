package main

import "os"

func main() {
	//TODO
	//load image file
	//sort pixels
	// preferably with modular algorithm
	// also this step should be able to run in multiple go routines
	//save output
}

func openImage() error {
	file, err := os.Open("path/to/image.img")
	defer file.Close()
	if err != nil {
		return err
	}

	//img := image.NewNRGBA
	return nil
}
