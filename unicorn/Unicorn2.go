// Version 2 of unicorn, uses gif.GIF as store of pixels
// & a separate goroutine to render it so it can do
// animated GIFs
package unicorn

import (
	"image"
	"image/color/palette"
	"image/gif"
	"time"
)

// Unicorn2 ...
// Interface for interacting with the Unicorn HAT HD
// Implemented by both real & fake unicorns.
type Unicorn2 interface {
	// Change the image
	GetGif() *gif.GIF
	SetGif(*gif.GIF)

	// Starts the rendering goroutine
	StartRender() chan bool

	// Must be implemented to actually render the image to device
	renderImage(image.Image)

	// Required for os to not think we're stuck
	MainLoop()
}

// BaseUnicorn2 ...
// Common to both real & fake unicorns!
// timing code for rendering & stopping rendering
type BaseUnicorn2 struct {
	g *gif.GIF
}

func (u *BaseUnicorn2) GetGif() *gif.GIF {
	return u.g
}

func (u *BaseUnicorn2) SetGif(g *gif.GIF) {
	u.g = g
}

func NewBaseUnicorn2() *BaseUnicorn2 {
	im := image.NewPaletted(
		image.Rect(0, 0, 16, 16),
		palette.WebSafe)

	gf := &gif.GIF{
		Image:           []*image.Paletted{im},
		Delay:           []int{50},
		Disposal:        []byte{gif.DisposalBackground},
		BackgroundIndex: 0, // This is black in the websafe palette
	}

	return &BaseUnicorn2{
		g: gf,
	}
}

// StartRenderBase ...
// Deals with the timing aspect of animated GIFs
func (u *BaseUnicorn2) StartRenderBase(renderImage func(image.Image)) chan bool {
	stopChan := make(chan bool)
	go func() {
		timer := time.NewTimer(0)
		imageIndex := 0
		running := true
		for running {
			select {
			case <-stopChan:
				timer.Stop()
				running = false
			case <-timer.C:
				gf := u.GetGif()

				// Image could change underneath us, but there should always be 1 image at least.
				if imageIndex >= len(gf.Image) {
					imageIndex = 0
				}

				im := gf.Image[imageIndex]
				delay := gf.Delay[imageIndex] //100ths of a second, 10^-2
				renderImage(im)

				timer.Reset(time.Duration(delay * 10000000)) // nanoseconds 10^-9 sec
				imageIndex = (imageIndex + 1) % len(gf.Image)
			}
		}
	}()
	return stopChan
}
