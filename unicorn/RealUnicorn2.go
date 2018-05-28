// +build linux,arm linux,arm64

package unicorn

import (
	"image"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ecc1/spi"
)

type RealUnicorn2 struct {
	BaseUnicorn2
	device *spi.Device
}

func (u *RealUnicorn2) renderImage(im image.Image) {
	b := im.Bounds()
	width, height := b.Dx(), b.Dy()
	sz := (width * height * 3) + 1
	write := make([]byte, sz)

	// Write leading bit
	write[0] = 0x72

	// Write color values
	ix := 1
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			col := im.At(x, y)
			r, g, b, _ := col.RGBA()
			write[ix] = byte(r)
			ix++
			write[ix] = byte(g)
			ix++
			write[ix] = byte(b)
			ix++
		}
	}
	// Write to the device
	err := u.device.Transfer(write)
	if err != nil {
		log.Printf("Error writing to SPI device %v", err)
	}
}

// NewUnicorn2 ...
// Constructs a new and improved unicorn from stuff and things
func NewUnicorn2() (*RealUnicorn2, error) {
	dev, err := spi.Open("/dev/spidev0.0", 9000000, 0)
	if err != nil {
		return nil, err
	}
	return &RealUnicorn2{
		BaseUnicorn2{},
		dev,
	}, nil
}

// StartRender ...
// Passes through to base to actually do the render
func (u *RealUnicorn2) StartRender() chan bool {
	return u.StartRenderBase(u.renderImage)
}

// MainLoop ...
// Just blocks until sigterm
func (u *RealUnicorn2) MainLoop() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
