// +build linux,arm linux,arm64

package unicorn

import (
	//"golang.org/x/exp/io/spi"
	"log"

	"github.com/ecc1/spi"
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
