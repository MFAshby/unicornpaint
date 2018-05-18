package main

import (
	//"golang.org/x/exp/io/spi"
    "github.com/ecc1/spi"
	"log"
)

type RealUnicorn struct {
	BaseUnicorn
	device *spi.Device
}

// NewReal ...
// Constructs a new real unicorn from fairy dust and sprinkles
func NewReal() (*RealUnicorn, error) {
	/*dev, err := spi.Open(&spi.Devfs{
		Dev:      "/dev/spidev0.0",
		Mode:     spi.Mode3,
		MaxSpeed: 9000000,
	})*/
    dev, err := spi.Open("/dev/spidev0.0", 9000000, 0)
	if err != nil {
		return nil, err
	}
    //dev.SetBitOrder(spi.LSBFirst)

	return &RealUnicorn{
		BaseUnicorn{
			pixels: makePixels(16, 16),
		},
		dev,
	}, nil
}

func (u *RealUnicorn) Show() {
	// Width * height * colours + leading bit
	width := int(u.GetWidth())
	height := int(u.GetHeight())
    sz := (width*height*3)+1
	write := make([]byte, sz)

	// Add the leading bit
	write[0] = 0x72
	// Add all the pixel values
	ix := 1
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			for j := 0; j < 3; j++ {
				write[ix] = byte(u.pixels[x][y][j])
				ix++
			}
		}
	}
	// Write to the device
	//err := u.device.Tx(write, nil)
    err := u.device.Transfer(write)
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
