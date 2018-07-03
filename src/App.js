import React, { Component } from 'react'
import './App.css'
import { rgb, xy, findContiguousPixels, getPixel, rotatePixelsClock, rotatePixelsCounterClock, readGifFrames } from './Utils'
import { setPixel, clear, noop, save, load, addFrame, removeFrame } from './Actions'
import Palette from './Palette'
import PaintArea from './PaintArea'
import Toolkit from './Toolkit'
import ColorIndicator from './ColorIndicator'
import ConnectedIndicator from './ConnectedIndicator'
import LoadDialog from './LoadDialog'
import SaveDialog from './SaveDialog'
import FrameControl from './FrameControl'
import { Timeline } from 'react-twitter-widgets'
import '@fortawesome/fontawesome-free-webfonts/css/fontawesome.css'
import '@fortawesome/fontawesome-free-webfonts/css/fa-solid.css'
import '@fortawesome/fontawesome-free-webfonts/css/fa-regular.css'
import '@fortawesome/fontawesome-free-webfonts/css/fa-brands.css'

const tools = [
  {
    name: "paint",
    icon: "fas fa-pencil-alt",
    action: function (x, y, frame) {
      let { r, g, b } = rgb(this.state.selectedColor)
      setPixel(this._websocket, x, y, r, g, b, frame)
    }
  },
  {
    name: "fill",
    icon: "fab fa-bitbucket",
    action: function (x, y, frame) {
      let pixelsToColor = findContiguousPixels(x, y, this.state.frames[frame])
      pixelsToColor.forEach((coord) => {
        let px = { ...xy(coord), ...rgb(this.state.selectedColor) }
        setPixel(this._websocket, px.x, px.y, px.r, px.g, px.b, frame)
      })
    }
  },
  {
    name: "erase",
    icon: "fas fa-eraser",
    action: function (x, y, frame) {
      setPixel(this._websocket, x, y, 0, 0, 0, frame)
    },
  },
  {
    name: "pick",
    icon: "fas fa-eye-dropper",
    action: function (x, y, frame) {
      let pixels = this.state.frames[frame]
      let color = getPixel(x, y, pixels)
      this.setState({ selectedColor: color })
    },
  },
  {
    name: "add frame",
    icon: "far fa-plus-square",
    onSelect: function () {
      addFrame(this._websocket, this.state.selectedFrame+1, 50)
    }
  },
  {
    name: "remove frame",
    icon: "far fa-minus-square",
    onSelect: function () {
      removeFrame(this._websocket, this.state.selectedFrame)
    }
  },
  {
    name: "rotate-clockwise",
    icon: "fas fa-redo",
    onSelect: function () {
      let pixels = this.state.frames[this.state.selectedFrame]
      let newPixels = rotatePixelsClock(pixels)
      this._setAllPixels(newPixels)
    }
  },
  {
    name: "rotate-anticlockwise",
    icon: "fas fa-undo",
    onSelect: function () {
      let pixels = this.state.frames[this.state.selectedFrame]
      let newPixels = rotatePixelsCounterClock(pixels)
      this._setAllPixels(newPixels)
    }
  },
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
      frames: [],
      saves: [],
      imageData: "",

      // Local data
      selectedFrame: 0,
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
    let { imageData, saves } = JSON.parse(data)
    let frames = readGifFrames(imageData)

    let selectedFrame = this.state.selectedFrame >= frames.length ? 
      frames.length - 1 :
      this.state.selectedFrame
    
    this.setState({
      frames: frames,
      saves: saves,
      selectedFrame: selectedFrame,
      imageData: imageData,
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
    let webSocketProto = window.location.protocol === "https:" ? "wss:" : "ws:"
    var host = window.location.host
    // if (true) {
    //   console.log("Dev mode overriding port to 3001")
    //   host = window.location.hostname + ":3001"
    // }
    this._websocket = new WebSocket(`${webSocketProto}//${host}/ws`)
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
    let selectedFrame = this.state.selectedFrame
    action.bind(this)(x, y, selectedFrame)
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
        let selectedFrame = this.state.selectedFrame
        setPixel(this._websocket, x, y, r, g, b, selectedFrame)
      }
    }
  }

  render() {
    let { selectedFrame, 
      frames, 
      imageData,
      connected, 
      selectedTool, 
      selectedColor,
      saves,
      showingLoad,
      showingSave } = this.state

    let pixels = frames[selectedFrame] || []

    return (
      <div className="container">
        <div className="header">
          <h1>Unicorn Paint!</h1>
          <ConnectedIndicator connected={connected} />
        </div>
        <div className="framePaintContainer">
          <div className="paintContainer">
            <PaintArea
              data={pixels}
              onTool={this._applyTool} />
            <Toolkit
              tools={tools}
              selectedTool={selectedTool}
              onSelectTool={this._selectTool} />
            <Palette
              selectedColor={selectedColor}
              onSelectColor={(color) => this.setState({ selectedColor: color })} />  
            <ColorIndicator color={selectedColor} />
          </div>
          <FrameControl 
              frames={frames} 
              selectedFrame={selectedFrame}
              onFrameSelected={ frame => this.setState({selectedFrame: frame}) }
              imageData={imageData}/>
        </div>
        <div className="liveFeed">
          <Timeline 
            dataSource={{
              sourceType: 'profile',
              screenName: 'UnicornPaint',
            }}
            options={{
              tweetLimit: 1,
              chrome: "noheader nofooter noborders noscrollbar",
              width: "300px"
            }}/>
        </div>
        <div>
          {
            showingLoad
            && <LoadDialog
              saves={saves}
              onLoad={(drawing) => this._loadDrawing(drawing)}
              onClose={() => this.setState({ showingLoad: false })} />
          }
        </div>
        <div>
          {
            showingSave
            && <SaveDialog
              saves={saves}
              onSave={(name) => this._saveDrawing(name)}
              onClose={() => this.setState({ showingSave: false })} />
          }
        </div>
      </div>
    );
  }
}

export default App;
