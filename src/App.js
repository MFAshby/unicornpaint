import React, { Component } from 'react'
import './App.css'
import { rgb, xy, findContiguousPixels, getPixel, rotatePixelsClock, rotatePixelsCounterClock } from './Utils'
import { setPixel, clear, noop, save, load } from './Actions'
import Palette from './Palette'
import PaintArea from './PaintArea'
import Toolkit from './Toolkit'
import ColorIndicator from './ColorIndicator'
import ConnectedIndicator from './ConnectedIndicator'
import LoadDialog from './LoadDialog'
import SaveDialog from './SaveDialog'

const tools = [
  {
    name: "paint",
    icon: "fas fa-pencil-alt",
    action: function (x, y) {
      let { r, g, b } = rgb(this.state.selectedColor)
      setPixel(this._websocket, x, y, r, g, b)
    }
  },
  {
    name: "fill",
    icon: "fab fa-bitbucket",
    action: function (x, y) {
      let pixelsToColor = findContiguousPixels(x, y, this.state.pixels)
      pixelsToColor.forEach((coord) => {
        let px = { ...xy(coord), ...rgb(this.state.selectedColor) }
        setPixel(this._websocket, px.x, px.y, px.r, px.g, px.b)
      })
    }
  },
  {
    name: "erase",
    icon: "fas fa-eraser",
    action: function (x, y) {
      setPixel(this._websocket, x, y, 0, 0, 0)
    },
  },
  {
    name: "pick",
    icon: "fas fa-eye-dropper",
    action: function (x, y) {
      let color = getPixel(x, y, this.state.pixels)
      this.setState({ selectedColor: color })
    },
  },
  {
    name: "rotate-clockwise",
    icon: "fas fa-redo",
    onSelect: function () {
      let newPixels = rotatePixelsClock(this.state.pixels)
      this._setAllPixels(newPixels)
    }
  },
  {
    name: "rotate-anticlockwise",
    icon: "fas fa-undo",
    onSelect: function () {
      let newPixels = rotatePixelsCounterClock(this.state.pixels)
      this._setAllPixels(newPixels)
    }
  },
  // { 
  //   name: "lighten", 
  //   icon: "far fa-sun"
  // },
  // { 
  //   name: "darken", 
  //   icon: "fas fa-sun"
  // },
  {
    name: "save",
    icon: "fas fa-save",
    onSelect: function () {
      this.setState({ showingSave: true })
    }
  },
  {
    name: "load",
    icon: "fas fa-save",
    onSelect: function () {
      this.setState({ showingLoad: true })
    }
  },
  {
    name: "trash",
    icon: "fas fa-trash",
    onSelect: function () {
      clear(this._websocket)
    }
  },
]

class App extends Component {
  constructor(props) {
    super(props)

    this.state = {
      connected: false,
      // Data from server
      pixels: [],
      saves: [],
      // Local data
      selectedColor: [0, 0, 0],
      selectedTool: tools[0],
      showingSave: false,
      showingLoad: false,
    }
    this._applyTool = this._applyTool.bind(this)
    this._selectTool = this._selectTool.bind(this)
    this._connectWebsocket = this._connectWebsocket.bind(this)
    this._onMessage = this._onMessage.bind(this)
    this._onOpen = this._onOpen.bind(this)
    this._onClose = this._onClose.bind(this)
    this._onError = this._onError.bind(this)
    this._loadDrawing = this._loadDrawing.bind(this)
    this._saveDrawing = this._saveDrawing.bind(this)
    this._connectWebsocket()
  }

  _onMessage({ data }) {
    let state = JSON.parse(data)
    this.setState({
      ...state // Includes pixels and saves
    })
  }

  _onOpen() {
    this.setState({ connected: true })
    noop(this._websocket)
  }

  _onClose() {
    this.setState({ connected: false })
    this._connectWebsocket()
  }

  _onError() {
    this.setState({ connected: false })
    this._connectWebsocket()
  }

  _connectWebsocket() {
    // this._websocket = new WebSocket('ws://' + window.location.hostname + ':3001/ws')
    this._websocket = new WebSocket('ws://shinypi:3001/ws')
    this._websocket.onmessage = this._onMessage
    this._websocket.onopen = this._onOpen
    this._websocket.onclose = this._onClose
    this._websocket.onerror = this._onError
  }

  _applyTool(x, y) {
    let tool = this.state.selectedTool
    if (!tool) {
      return
    }
    let action = tool.action
    if (!action) {
      return
    }
    action.bind(this)(x, y)
  }

  _selectTool(tool) {
    let selectAction = tool.onSelect
    if (selectAction) {
      selectAction.bind(this)()
    } else {
      this.setState({ selectedTool: tool })
    }
  }

  _loadDrawing(name) {
    load(this._websocket, name)
    this.setState({ showingLoad: false })
  }

  _saveDrawing(name) {
    save(this._websocket, name)
    this.setState({ showingSave: false })
  }

  _setAllPixels(newPixels) {
    let width = newPixels.length
    let height = newPixels[0].length
    for (var x = 0; x < width; x++) {
      for (var y = 0; y < height; y++) {
        let px = getPixel(x, y, newPixels)
        let { r, g, b } = rgb(px)
        setPixel(this._websocket, x, y, r, g, b)
      }
    }
  }

  render() {
    return (
      <div className="App">
        <ConnectedIndicator connected={this.state.connected} />
        <Toolkit
          tools={tools}
          selectedTool={this.state.selectedTool}
          onSelectTool={this._selectTool} />
        <PaintArea
          data={this.state.pixels}
          onTool={this._applyTool} />
        <Palette
          selectedColor={this.state.selectedColor}
          onSelectColor={(color) => this.setState({ selectedColor: color })} />
        <ColorIndicator color={this.state.selectedColor} />
        <div>
          {
            this.state.showingLoad
            && <LoadDialog
              saves={this.state.saves}
              onLoad={(drawing) => this._loadDrawing(drawing)}
              onClose={() => this.setState({ showingLoad: false })} />
          }
        </div>
        <div>
          {
            this.state.showingSave
            && <SaveDialog
              saves={this.state.saves}
              onSave={(name) => this._saveDrawing(name)}
              onClose={() => this.setState({ showingSave: false })} />
          }
        </div>

      </div>
    );
  }
}

export default App;
