// +build linux,386 linux,amd64

package unicorn

import (
	"github.com/veandco/go-sdl2/sdl"
)

type FakeUnicorn struct {
	BaseUnicorn
	*BaseFakeUnicorn
}

// NewUnicorn ...
// Constructs a new fake unicorn out of paint and glue
func NewUnicorn() (*FakeUnicorn, error) {
	width := uint8(16)
	height := uint8(16)

	baseFake, err := NewBaseFakeUnicorn(300, 300)
	if err != nil {
		return nil, err
	}

	unicorn := &FakeUnicorn{
		BaseUnicorn{
			pixels: makePixels(width, height),
		},
		baseFake,
	}

	return unicorn, nil
}

func (f *FakeUnicorn) Show() {
	width, height := f.GetWidth(), f.GetHeight()
	for x := uint8(0); x < width; x++ {
		for y := uint8(0); y < height; y++ {
			r, g, b := Rgb(f.pixels[x][y])
			if err := f.renderer.SetDrawColor(r, g, b, uint8(255)); err != nil {
				panic(err)
			}
			cellWidth := f.displayWidth / int32(width)
			cellHeight := f.displayHeight / int32(height)
			if err := f.renderer.FillRect(&sdl.Rect{
				X: cellWidth * int32(x),
				Y: f.displayHeight - (cellHeight * int32(y)) - cellHeight, // SDL Y coordinate is from the top
				W: cellWidth,
				H: cellHeight,
			}); err != nil {
				panic(err)
			}
		}
	}
	f.renderer.Present()
}

func (f *FakeUnicorn) Off() {
	f.Close()
}
