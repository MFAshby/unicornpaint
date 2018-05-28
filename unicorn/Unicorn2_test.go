package unicorn

import (
	"bytes"
	"image/gif"
	"testing"
	"time"
)

func gifAsset(name string) (*gif.GIF, error) {
	data, err := Asset(name)
	if err != nil {
		return nil, err
	}

	g, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return g, nil
}

func TestAnimated(t *testing.T) {
	un, err := NewUnicorn2()
	if err != nil {
		t.Errorf("Failed to create fake unicorn :( %v", err)
		return
	}
	defer un.Close()

	g, err := gifAsset("data/sample.gif")
	if err != nil {
		t.Errorf("Failed to load asset %v", err)
		return
	}

	un.SetGif(g)
	stopChan := un.StartRender()

	// Stop after 3
	time.Sleep(3 * time.Second)
	stopChan <- true

	// Leave it for a sec
	time.Sleep(1 * time.Second)
	g2, err := gifAsset("data/sample2.gif")
	if err != nil {
		t.Errorf("Failed to load asset %v", err)
		return
	}
	un.SetGif(g2)
	stopChan = un.StartRender()

	// Stop after 5
	time.Sleep(5 * time.Second)
	stopChan <- true

	// Make sure it's stopped
	time.Sleep(2 * time.Second)
}
