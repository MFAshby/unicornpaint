package main

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"os"

	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/websocket"

	"github.com/MFAshby/unicornpaint/unicorn"
)

var (
	un unicorn.Unicorn2

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(*http.Request) bool { return true },
	}
	delCh       = make(chan *websocket.Conn)
	addCh       = make(chan *websocket.Conn)
	broadcastCh = make(chan interface{})
	clients     = map[*websocket.Conn]bool{}
)

// Types of commands
type commandType string

const (
	noop        commandType = "NO_OP"
	setPixel    commandType = "SET_PIXEL"
	clear       commandType = "CLEAR"
	save        commandType = "SAVE"
	load        commandType = "LOAD"
	addFrame    commandType = "ADD_FRAME"
	removeFrame commandType = "REMOVE_FRAME"
)

type Command struct {
	Type     commandType `json:"type"`
	X        uint8       `json:"x"`
	Y        uint8       `json:"y"`
	R        uint8       `json:"r"`
	G        uint8       `json:"g"`
	B        uint8       `json:"b"`
	Frame    int         `json:"frame"`
	Delay    int         `json:"delay"`
	SaveName string      `json:"saveName"`
}

const (
	savesDir = "saves"
)

type State struct {
	Saves     []string `json:"saves"`
	ImageData []byte   `json:"imageData"`
}

func getGifBytes() ([]byte, error) {
	g := un.GetGif()
	buf := &bytes.Buffer{}
	err := gif.EncodeAll(buf, g)
	return buf.Bytes(), err
}

func getState() *State {
	infos, err := ioutil.ReadDir(savesDir)
	if err != nil {
		log.Printf("Error reading saves dir %v", err)
	}
	saveFileNames := make([]string, len(infos))
	for ix, info := range infos {
		saveFileNames[ix] = info.Name()
	}

	// Write out as GIF file
	gifBytes, err := getGifBytes()
	if err != nil {
		log.Printf("Error getting GIF bytes %v", err)
	}

	return &State{
		ImageData: gifBytes,
		Saves:     saveFileNames,
	}
}

func savePic(saveFileName string) {
	f, err := os.Create(path.Join(savesDir, saveFileName))
	if err != nil {
		log.Printf("Failed to open file for saving %v", err)
		return
	}
	defer f.Close()

	g := un.GetGif()
	err = gif.EncodeAll(f, g)
}

func loadPic(saveFileName string) {
	f, err := os.Open(path.Join(savesDir, saveFileName))
	if err != nil {
		log.Printf("Failed to read file %v", err)
		return
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		log.Printf("Failed to decode GIF file %v", err)
		return
	}

	un.SetGif(g)
}

func doSetPixel(x, y, r, g, b uint8, frame int) {
	gf := un.GetGif()
	if int(frame) >= len(gf.Image) {
		log.Printf("Tried to set pixel in frame out of range %v", frame)
		return
	}
	im := gf.Image[frame]
	im.Set(int(x), int(y), color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(255),
	})
}

func doAddFrame(frame int, delay int) {
	gf := un.GetGif()
	if frame < 1 || frame > len(gf.Image) {
		log.Printf("Trying to add frame in invalid position %v", frame)
		return
	}

	// Copy the image before it
	sourceImage := gf.Image[frame-1]

	buf := &bytes.Buffer{}
	if err := gif.Encode(buf, sourceImage, nil); err != nil {
		log.Printf("Failed to encode image %v", err)
		return
	}
	im, err := gif.Decode(buf)
	if err != nil {
		log.Printf("Failed to decode image %v", err)
		return
	}

	newImage, ok := im.(*image.Paletted)
	if !ok {
		log.Printf("Wrong image type %v", err)
	}

	gf.Image = append(gf.Image[:frame], append([]*image.Paletted{newImage}, gf.Image[frame:]...)...)
	gf.Delay = append(gf.Delay[:frame], append([]int{delay}, gf.Delay[frame:]...)...)
	gf.Disposal = append(gf.Disposal[:frame], append([]byte{gif.DisposalBackground}, gf.Disposal[frame:]...)...)
}

func doRemoveFrame(frame int) {
	// Can't remove the first frame
	gf := un.GetGif()
	frameCount := len(gf.Image)
	if frame < 0 || frame >= frameCount {
		log.Printf("Trying to remove frame %v which is invalid", frame)
		return
	}

	if frameCount == 1 {
		log.Printf("Trying to remove the last frame, ignoring")
		return
	}

	gf.Image = append(gf.Image[:frame], gf.Image[frame+1:]...)
	gf.Delay = append(gf.Delay[:frame], gf.Delay[frame+1:]...)
	gf.Disposal = append(gf.Disposal[:frame], gf.Disposal[frame+1:]...)
}

func doClear(frame int) {
	gf := un.GetGif()
	if frame >= len(gf.Image) {
		log.Printf("Trying to clear invalid frame %v", frame)
		return
	}

	im := gf.Image[frame]
	b := im.Bounds()
	width := b.Dx()
	height := b.Dy()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			im.SetColorIndex(x, y, 0) // 0 in WebSafe colors is black
		}
	}
}

func upgradeHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade someone %v", err)
		return
	}

	// Add you to my list, remember to remove later
	addCh <- conn
	defer func() { delCh <- conn }()

	// Get you up to speed
	err = conn.WriteJSON(getState())
	if err != nil {
		log.Printf("Failed to send initial state %v", err)
		return
	}

	// Read & execute commands in a loop until error
	for {
		cmd := Command{}
		err = conn.ReadJSON(&cmd)
		if err != nil {
			log.Printf("Error reading from client %v", err)
			break
		}

		switch cmd.Type {
		case noop:
			// Don't do anything
		case setPixel:
			doSetPixel(cmd.X, cmd.Y, cmd.R, cmd.G, cmd.B, cmd.Frame)
		case clear:
			doClear(cmd.Frame)
		case save:
			savePic(cmd.SaveName)
		case load:
			loadPic(cmd.SaveName)
		case addFrame:
			doAddFrame(cmd.Frame, cmd.Delay)
		case removeFrame:
			doRemoveFrame(cmd.Frame)
		}
		// Pretty much all commands demand a change of state,
		// do a broadcast each time
		broadcastCh <- getState()
	}
}

func handleClients() {
	for {
		select {
		case c := <-addCh:
			clients[c] = true
		case c := <-delCh:
			delete(clients, c)
		case msg := <-broadcastCh:
			doBroadcast(msg)
		}
	}
}

func main() {
	uni, err := unicorn.NewUnicorn2()
	if err != nil {
		log.Fatalf("Couldn't get a unicorn :( %v", err)
	}
	un = uni
	un.StartRender()

	log.Println("Starting server on port 3001")
	http.Handle("/", http.FileServer(http.Dir("build")))
	http.HandleFunc("/ws", upgradeHandler)
	go http.ListenAndServe(":3001", nil)
	go handleClients()
	un.MainLoop()
}

func doBroadcast(obj interface{}) {
	for client := range clients {
		err := client.WriteJSON(obj)

		if err != nil {
			log.Printf("Error writing to client, closing %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
