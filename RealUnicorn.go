package main

import (
	"golang.org/x/exp/io/spi"
	"log"
)

type RealUnicorn struct {
	BaseUnicorn
	device *spi.Device
}

// NewReal ...
// Constructs a new real unicorn from fairy dust and sprinkles
func NewReal() (*RealUnicorn, error) {
	dev, err := spi.Open(&spi.Devfs{
		Dev:      "/dev/spidev0.0",
		Mode:     spi.Mode0,
		MaxSpeed: 9000000,
	})
	if err != nil {
		return nil, err
	}

	return &RealUnicorn{
		BaseUnicorn{
			pixels: makePixels(16, 16),
		},
		dev,
	}, nil
}

func (u *RealUnicorn) Show() {
	// Width * height * colours + leading bit
	width := u.GetWidth()
	height := u.GetHeight()
	write := make([]byte, (width*height*3)+1)

	// Add the leading bit
	write[0] = 0x72
	// Add all the pixel values
	ix := 1
	for x := uint8(0); x < width; x++ {
		for y := uint8(0); y < height; y++ {
			for j := 0; j < 3; j++ {
				write[ix] = u.pixels[x][y][j]
				ix++
			}
		}
	}
	// Write to the device
	err := u.device.Tx(write, nil)
	if err != nil {
		log.Printf("Error writing to SPI device %v", err)
	}
}
func (u *RealUnicorn) Off() {
	u.Close()
}
func (u *RealUnicorn) Close() error {
	return u.device.Close()
}
