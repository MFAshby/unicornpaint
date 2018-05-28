// Version 2 of unicorn, uses gif.GIF as store of pixels
// & a separate goroutine to render it so it can do
// animated GIFs
package unicorn

import (
	"image"
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
	StartRender()

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

// StartRender ...
// Starts rendering the image. If it's an animated image,
// renders animation frames. Return a channel to stop the
// image rendering.
func (u *FakeUnicorn2) StartRender() chan bool {
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
				im := gf.Image[imageIndex]
				delay := gf.Delay[imageIndex] //100ths of a second, 10^-2
				u.renderImage(im)
				
				timer.Reset(time.Duration(delay * 10000000)) // nanoseconds 10^-9 sec
				imageIndex = (imageIndex + 1) % len(gf.Image)
			}
		}
	}()
	return stopChan
}