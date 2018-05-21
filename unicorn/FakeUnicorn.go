// +build linux,386 linux,amd64

package unicorn

import (
	"github.com/veandco/go-sdl2/sdl"
)

type FakeUnicorn struct {
	BaseUnicorn
	displayWidth  int32
	displayHeight int32
	window        *sdl.Window
	renderer      *sdl.Renderer
}

// NewUnicorn ...
// Constructs a new fake unicorn out of paint and glue
func NewUnicorn() (*FakeUnicorn, error) {
	width := uint8(16)
	height := uint8(16)
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}

	unicorn := &FakeUnicorn{
		BaseUnicorn{
			pixels: makePixels(width, height),
		},
		300,
		300,
		nil,
		nil,
	}
	if err := unicorn.createWindow(); err != nil {
		unicorn.Close()
		return nil, err
	}
	if err := unicorn.createRenderer(); err != nil {
		unicorn.Close()
		return nil, err
	}
	return unicorn, nil
}

func (f *FakeUnicorn) createWindow() error {
	window, err := sdl.CreateWindow("Fake Unicorn",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		f.displayWidth,
		f.displayHeight,
		sdl.WINDOW_SHOWN)
	f.window = window
	return err
}

func (f *FakeUnicorn) createRenderer() error {
	renderer, err := sdl.CreateRenderer(f.window, -1, sdl.RENDERER_ACCELERATED)
	f.renderer = renderer
	return err
}

func (f *FakeUnicorn) Close() error {
	if f.window != nil {
		f.window.Destroy()
	}
	if f.renderer != nil {
		f.renderer.Destroy()
	}
	return nil
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

// MainLoop ...
// Handle UI events so OS doesn't think we're frozen
func (f *FakeUnicorn) MainLoop() {
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
