import React, { Component } from 'react'
import { rgb, getPixel } from './Utils'

/**
 * Expects props: 
 * onTool = (x, y) => {}
 *  Callback for when a cell is clicked
 * 
 * data = [[]]
 *  3d array - rows, cells, RGB value for each cell
 */
export default class PaintArea extends Component {
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
    let data = this.props.data
    let height = data[0] ? data[0].length : 0
    let rows = []
    for (var y=height-1; y>=0; y--) {
      let cells = []
      for (var x=0; x<data.length; x++) {
        let {r, g, b} = rgb(getPixel(x, y, data))
        let ix = x
        let iy = y
        cells.push(<td
          onMouseMove={() => this.handleMouseMove(ix, iy)}
          onClick={() => this.props.onTool(ix, iy)}
          className="paintareacell"
          style={{
            background: `rgb(${r},${g},${b})`
          }}
          key={(ix * 100000) + iy}/>)
      }
      rows.push(<tr key={y}>{cells}</tr>)
    }
    
    return (
      <table 
        className="paintarea"
        draggable={false}>
        <tbody>
          {rows}
        </tbody>
      </table>
    )
  }
}