// +build linux,arm linux,arm64

package unicorn

import (
	"github.com/ecc1/spi"
	"log"
	"os"
)

type RealUnicorn struct {
	BaseUnicorn
	device *spi.Device
}

// NewUnicorn ...
// Constructs a new real unicorn from fairy dust and sprinkles
func NewUnicorn() (*RealUnicorn, error) {
	dev, err := spi.Open("/dev/spidev0.0", 9000000, 0)
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
	width := int(u.GetWidth())
	height := int(u.GetHeight())
	sz := (width * height * 3) + 1
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

// MainLoop ...
// Do nothing until SIGTERM, then close the SPI library
func MainLoop() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<- c
	Close()
}
