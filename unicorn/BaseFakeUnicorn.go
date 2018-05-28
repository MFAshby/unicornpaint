// +build linux,386 linux,amd64

package unicorn

import "github.com/veandco/go-sdl2/sdl"

// BaseFakeUnicorn ...
// The base for FakeUnicorn & FakeUnicorn2
// Share the SDL code.
type BaseFakeUnicorn struct {
	displayWidth  int32
	displayHeight int32
	window        *sdl.Window
	renderer      *sdl.Renderer
}

func NewBaseFakeUnicorn(width, height int32) (*BaseFakeUnicorn, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}

	unicorn := &BaseFakeUnicorn{
		displayWidth:  width,
		displayHeight: height,
		window:        nil,
		renderer:      nil,
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

func (f *BaseFakeUnicorn) createWindow() error {
	window, err := sdl.CreateWindow("Fake Unicorn",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		f.displayWidth,
		f.displayHeight,
		sdl.WINDOW_SHOWN)
	f.window = window
	return err
}

func (f *BaseFakeUnicorn) createRenderer() error {
	renderer, err := sdl.CreateRenderer(f.window, -1, sdl.RENDERER_ACCELERATED)
	f.renderer = renderer
	return err
}

func (f *BaseFakeUnicorn) Close() error {
	if f.window != nil {
		f.window.Destroy()
	}
	if f.renderer != nil {
		f.renderer.Destroy()
	}
	return nil
}

// MainLoop ...
// Handle UI events so OS doesn't think we're frozen
func (f *BaseFakeUnicorn) MainLoop() {
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
