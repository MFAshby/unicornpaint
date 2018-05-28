// Version 2 of unicorn, uses gif.GIF as store of pixels
// & a separate goroutine to render it so it can do
// animated GIFs
package unicorn

import (
	"image/gif"
)

// Unicorn2 ...
// Interface for interacting with the Unicorn HAT HD
// Implemented by both real & fake unicorns.
type Unicorn2 interface {
	GetGif() *gif.GIF
	SetGif(*gif.GIF)

	StartRender()
	// Required for os to not think we're stuck
	MainLoop()
}

type BaseUnicorn2 struct {
	g *gif.GIF
}

func (u *BaseUnicorn2) GetGif() *gif.GIF {
	return u.g
}

func (u *BaseUnicorn2) SetGif(g *gif.GIF) {
	u.g = g
}
