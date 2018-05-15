import React, { Component } from 'react'
import logo from './logo.svg'
import './App.css'

const NO_OP = 'NO_OP'
const SET_PIXEL = 'SET_PIXEL'
const CLEAR = 'CLEAR'

function sendAction(websocket, action) {
  let actionStr = JSON.stringify(action)
  websocket.send(actionStr)
}

function setPixel(websocket, x, y, r, g, b) {
  sendAction(websocket, {
    type: SET_PIXEL,
    x: x,
    y: y,
    r: r, 
    g: g,
    b: b
  })
}

function clear(websocket) {
  sendAction(websocket, { type: CLEAR })
}

function noop(websocket) {
  sendAction(websocket, { type: NO_OP })
}

function rgb(item) {
  return {
    r: item[0],
    g: item[1],
    b: item[2]
  }
}

const colorPalette = [
  [0,0,0],
  [132,0,0],
  [0,132,0],
  [132,132,0],
  [0,0,132],
  [132,0,132],
  [0,132,132],
  [132,132,132],
  [198,198,198],
  [255,0,0],
  [0,255,0],
  [255,255,0],
  [0,0,255],
  [255,0,255],
  [0,255,255],
  [255,255,255],
]


class Palette extends Component {
  render() {
    let paletteListItems = colorPalette.map((item) => {
      var className = "paletteitem"
      let { r, g, b } = rgb(item)
      let selected = rgb(this.props.selectedColor)
      if (r === selected.r && g === selected.g && b === selected.b) {
        className += " selected"
      }
      return <div
        onClick={() => this.props.onSelectColor([r, g, b])}
        style={{background: `rgb(${r},${g},${b})`}} 
        className={className}
        key={r*10000+g*1000+b}/>
    })
    return (
      <div className="palette">
        {paletteListItems}
      </div>)
  }
}

class PaintArea extends Component {
  constructor(props) {
    super(props)
    this.handleMouseMove = this.handleMouseMove.bind(this)

    this._onMouseDown = this._onMouseDown.bind(this)
    this._onMouseUp = this._onMouseUp.bind(this)
    this.state = {
      mouseDown: false
    }
  }

  _onMouseDown() {
    this.setState({ mouseDown: true })
  }

  _onMouseUp() {
    this.setState({ mouseDown: false })
  }

  componentWillMount() {
    document.addEventListener('mousedown', this._onMouseDown)
    document.addEventListener('mouseup', this._onMouseUp)
  }

  componentWillUnmount() {
    document.removeEventListener('mousedown', this._onMouseDown)
    document.removeEventListener('mouseup', this._onMouseUp)
  }

  handleMouseMove(x, y) {
    if (this.state.mouseDown) {
      this.props.onTool(x, y)
    }
  }

  render() {
    let cells = this.props.data.map((row, iy) => {
      let rowCells = row.map((cell, ix) => {
          let r = cell[0]
          let g = cell[1]
          let b = cell[2]
          return <td
            onMouseMove={() => this.handleMouseMove(ix, iy)}
            onClick={() => this.props.onTool(ix, iy)}
            className="paintareacell"
            style={{
              background: `rgb(${r},${g},${b})`
            }}
            key={(ix * 100000) + iy}/>
      })
      return <tr key={iy}>{rowCells}</tr>
    })
    
    return (
      <table 
        className="paintarea"
        draggable={false}>
        <tbody>
          {cells}
        </tbody>
      </table>
    )
  }
}

class App extends Component {
  constructor(props) {
    super(props)

    this.state = {
      connected: false,
      pixels: [],
      selectedColor: [0, 0, 0]
    }
    this._connectWebsocket = this._connectWebsocket.bind(this)
    this._onMessage = this._onMessage.bind(this)
    this._onOpen = this._onOpen.bind(this)
    this._onClose = this._onClose.bind(this)
    this._onError = this._onError.bind(this)

    this.paintCell = this.paintCell.bind(this)
    this._connectWebsocket()
  }

  _onMessage({data}) {
    let pixels = JSON.parse(data)
    pixels.forEach(row => {
      let rowNums = row.map((col) => col.some((it) => it > 0))
      console.log(...rowNums)
    });

    this.setState({
      pixels: pixels
    })
  }

  _onOpen() {
    this.setState({connected: true})
    noop(this._websocket)
  }

  _onClose() {
    this.setState({connected: false})
    this._connectWebsocket()
  }

  _onError() {
    this.setState({connected: false})
    this._connectWebsocket()
  }

  _connectWebsocket() {
    this._websocket = new WebSocket('ws://' + window.location.hostname + ':3001/ws')
    this._websocket.onmessage = this._onMessage
    this._websocket.onopen = this._onOpen
    this._websocket.onclose = this._onClose
    this._websocket.onerror = this._onError
  }

  paintCell(x, y) {
    let { r, g, b } = rgb(this.state.selectedColor)
    setPixel(this._websocket, x, y, r, g, b)
  }

  render() {
    let connectedText = this.state.connected ? "Connected" : "Not connected"
    return (
      <div className="App">
        <div>
          <span>{connectedText}</span>
        </div>
        <PaintArea 
          data={this.state.pixels}
          onTool={this.paintCell}/>
        <Palette 
          selectedColor={this.state.selectedColor}
          onSelectColor={(color) => this.setState({selectedColor: color})} />
      </div>
    );
  }
}

export default App;
