import React, { Component } from 'react'

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