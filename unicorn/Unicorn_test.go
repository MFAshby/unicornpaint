package unicorn

import (
	"reflect"
	"testing"
	"time"
)

func TestFakeUnicorn(t *testing.T) {
	unicorn, err := NewUnicorn()
	if err != nil {
		t.Errorf("Got an error making a fake unicorn, shouldn't happen")
	}
	defer unicorn.Close()

	// Check simple functions
	if unicorn.GetHeight() != 16 {
		t.Errorf("Height was wrong, expecting 16")
	}
	if unicorn.GetWidth() != 16 {
		t.Errorf("Width was wrong, expecting 16")
	}
	// Pixels should be black to start with
	pixels := unicorn.GetPixels()
	for x := uint8(0); x < 16; x++ {
		for y := uint8(0); y < 16; y++ {
			if !reflect.DeepEqual(pixels[x][y], []uint8{0, 0, 0}) {
				t.Errorf("Expecting black pixels to start with")
			}
		}
	}

	// Should be able to set a pixel, no others should change
	unicorn.SetPixel(0, 0, uint8(255), uint8(255), uint8(255))
	pixels = unicorn.GetPixels()
	if !reflect.DeepEqual(pixels[0][0], []uint8{255, 255, 255}) {
		t.Errorf("Pixel wasn't set when it should be")
	}
	for x := uint8(0); x < 16; x++ {
		for y := uint8(0); y < 16; y++ {
			if x == 0 && y == 0 {
				continue
			}
			if !reflect.DeepEqual(pixels[x][y], []uint8{0, 0, 0}) {
				t.Errorf("Expecting black pixels to start with")
			}
		}
	}

	// Should be able to set a second pixel
	unicorn.SetPixel(3, 4, uint8(4), uint8(5), uint8(6))
	pixels = unicorn.GetPixels()
	for x := uint8(0); x < 16; x++ {
		for y := uint8(0); y < 16; y++ {
			checkcolor := []uint8{0, 0, 0}
			if x == 0 && y == 0 {
				checkcolor = []uint8{255, 255, 255}
			} else if x == 3 && y == 4 {
				checkcolor = []uint8{4, 5, 6}
			}
			if !reflect.DeepEqual(pixels[x][y], checkcolor) {
				t.Errorf("Got incorrect pixel color at %d %d", x, y)
			}
		}
	}

	unicorn.Show()
	time.Sleep(time.Duration(500) * time.Millisecond)
	unicorn.SetPixel(10, 10, uint8(255), uint8(255), uint8(0))
	unicorn.Show()
	time.Sleep(time.Duration(500) * time.Millisecond)

	unicorn.SetPixel(0, 15, uint8(255), uint8(0), uint8(0))
	unicorn.Show()
	time.Sleep(time.Duration(500) * time.Millisecond)
}
