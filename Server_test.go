package main

import (
	"testing"
	"encoding/json"
	"reflect"
)

func TestStateJson(t *testing.T) {
	testPixels := make([][]uint8arr, 16)
	for x:=0; x<16; x++ {
		testPixels[x] = make([]uint8arr, 16)
		for y:=0; y<16; y++ {
			testPixels[x][y] = uint8arr{0, 2, 3}
		}
	}

	testSaves := []string{"bob", "sally", "blah"}
	s1 := &State{
		Pixels: testPixels,
		Saves: testSaves,
	}

	data, err := json.Marshal(s1)
	if err != nil {
		t.Errorf("Failed to write state to JSON %v", err)
	}

	s2 := &State{}
	err = json.Unmarshal(data, s2)
	if err != nil {
		t.Errorf("Failed to read state from JSON %v", err)
	}

	if !reflect.DeepEqual(s1, s2) {
		t.Errorf("Differences after serializing state %v %v", s1, s2)
	}
}

func TestCommandJson(t *testing.T) {
	cmd := Command{
		Type: noop,
		X: uint8(1),
		Y: uint8(1),
		R: uint8(1),
		G: uint8(1),
		B: uint8(1),
		SaveName: "testing",
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		t.Errorf("Error encoding command to JSON")
	}

	cmd2 := Command{}
	err = json.Unmarshal(data, &cmd2)
	if err != nil {
		t.Errorf("Error decoding command from JSON")
	}

	if !reflect.DeepEqual(cmd, cmd2) {
		t.Errorf("Differences after encoding JSON %v %v", cmd, cmd2)
	}

	cmd3 := Command{}
	testData := []byte(`{ "type": "SET_PIXEL", "x": 1, "y": 2, "r": 255, "g": 255, "b": 255 }`)
	err = json.Unmarshal(testData, &cmd3)
	if err != nil {
		t.Errorf("Error unmarshalling test JSON %v %s", err, testData)
	}

	cmd4 := Command{
		Type: setPixel,
		X: uint8(1),
		Y: uint8(2),
		R: uint8(255),
		G: uint8(255), 
		B: uint8(255),
	}

	if !reflect.DeepEqual(cmd3, cmd4) {
		t.Errorf("Json unmarshalled incorrectly to %v", cmd4)
	}
}