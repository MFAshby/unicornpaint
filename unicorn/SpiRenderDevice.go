package unicorn

import (
	"github.com/ecc1/spi"
)

type SpiRenderDevice struct {
	device *spi.Device
}

func NewSpiRenderDevice() (*SpiRenderDevice, error) {
	dev, err := spi.Open("/dev/spidev0.0", 9000000, 0)
	if err != nil {
		return nil, err
	}
	return &SpiRenderDevice{
		device: dev,
	}, nil
}
