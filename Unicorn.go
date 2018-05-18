package main

import (
	"log"
)

// Unicorn ...
// Object representing the Unicorn HAT to be controlled
type Unicorn interface {
	// Not all unicorns are the same size
	GetWidth() uint8
	GetHeight() uint8

	// Array of pixels, indexed x, then y, then color (rgb)
	GetPixels() [][][]uint8

	// Set an individual pixel
	SetPixel(x, y, r, g, b uint8)

	// Flip the display buffer
	Show()

	// Set all pixels back to black
	Clear()

	// Turns off the LEDs
	Off()
}

// GetUnicorn ...
// Get a unicorn. Tries to get you a real one,
// if it can't find one then gives you a fake one.
func GetUnicorn() (unicorn Unicorn, err error) {
	unicorn, err = NewReal()
	if err != nil {
		log.Println("Couldn't get a real unicorn, trying a fake one")
		unicorn, err = NewFake(uint8(16), uint8(16))
	}
	unicorn, err = NewFake(uint8(16), uint8(16))
	return
}

type BaseUnicorn struct {
	pixels        [][][]uint8
}

func (f *BaseUnicorn) GetWidth() uint8 {
	return uint8(len(f.pixels))
}
func (f *BaseUnicorn) GetHeight() uint8 {
	if len(f.pixels) > 0 {
		return uint8(len(f.pixels[0]))
	}
	return 0
}
func (f *BaseUnicorn) GetPixels() [][][]uint8 {
	return f.pixels
}
func (f *BaseUnicorn) SetPixel(x, y, r, g, b uint8) {
	f.pixels[x][y] = []uint8{r, g, b}
}
func (f *BaseUnicorn) Clear() {
	f.pixels = makePixels(f.GetWidth(), f.GetHeight())
}

func makePixels(width, height uint8) [][][]uint8 {
	pixels := make([][][]uint8, width)
	for x := uint8(0); x < width; x++ {
		pixels[x] = make([][]uint8, height)
		for y := uint8(0); y < height; y++ {
			pixels[x][y] = []uint8{0, 0, 0}
		}
	}
	return pixels
}


func rgb(pixel []uint8) (uint8, uint8, uint8) {
	return pixel[0], pixel[1], pixel[2]
}