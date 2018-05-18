package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
    //"github.com/veandco/go-sdl2/sdl"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

var (
	unicorn Unicorn

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
	noop     commandType = "NO_OP"
	setPixel commandType = "SET_PIXEL"
	clear    commandType = "CLEAR"
	save     commandType = "SAVE"
	load     commandType = "LOAD"
)

type Command struct {
	Type     commandType `json:"type"`
	X        uint8       `json:"x"`
	Y        uint8       `json:"y"`
	R        uint8       `json:"r"`
	G        uint8       `json:"g"`
	B        uint8       `json:"b"`
	SaveName string      `json:"saveName"`
}

const (
	savesDir = "saves"
)

type State struct {
	Saves  []string     `json:"saves"`
	Pixels [][]uint8arr `json:"pixels"`
}

// This is a trick to avoid the JSON serializer from
// interpreting uint8's as bytes and encoding them in base64
type uint8arr []uint8

func (u uint8arr) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
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

	// Irritating conversion function
	pixels := unicorn.GetPixels()
	width := unicorn.GetWidth()
	height := unicorn.GetHeight()

	px2 := make([][]uint8arr, width)
	for x := uint8(0); x < width; x++ {
		px2[x] = make([]uint8arr, height)
		for y := uint8(0); y < height; y++ {
			px2[x][y] = uint8arr(pixels[x][y])
		}
	}

	return &State{
		Pixels: px2,
		Saves:  saveFileNames,
	}
}

func savePic(saveFileName string) {
	pixels := unicorn.GetPixels()
	data, err := json.Marshal(pixels)
	if err != nil {
		log.Printf("Failed to save picture to JSON %v", err)
		return
	}
	err = ioutil.WriteFile(path.Join(savesDir, saveFileName), data, 0644)
	if err != nil {
		log.Printf("Failed to write to save file %v", err)
		return
	}
}

func loadPic(saveFileName string) {
	data, err := ioutil.ReadFile(path.Join(savesDir, saveFileName))
	if err != nil {
		log.Printf("Failed to read file %v", err)
		return
	}

	newPixels := [][][]uint8{}
	err = json.Unmarshal(data, &newPixels)
	if err != nil {
		log.Printf("Failed to parse file %v", err)
		return
	}

	width := len(newPixels)
	height := len(newPixels[0])
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b := rgb(newPixels[x][y])
			unicorn.SetPixel(uint8(x), uint8(y), r, g, b)
		}
	}
	unicorn.Show()
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
			unicorn.SetPixel(cmd.X, cmd.Y, cmd.R, cmd.G, cmd.B)
			unicorn.Show()
		case clear:
			unicorn.Clear()
			unicorn.Show()
		case save:
			savePic(cmd.SaveName)
		case load:
			loadPic(cmd.SaveName)
		}
		// Pretty much all commands demand a change of state,
		// do a broadcast each time
		broadcastCh <- getState()
	}
}

/*func handleSdlEvents() {
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			}
		}
	}
}*/

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
	uni, err := GetUnicorn()
	if err != nil {
		log.Fatalf("Couldn't get a unicorn :( %v", err)
	}
	unicorn = uni

	log.Println("Starting server on port 3001")
	http.Handle("/", http.FileServer(http.Dir("build")))
	http.HandleFunc("/ws", upgradeHandler)
	go http.ListenAndServe(":3001", nil)
	//go handleClients()
    handleClients()
	//handleSdlEvents()
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


