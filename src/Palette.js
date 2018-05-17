import React, {Component} from 'react'
import { rgb } from './Utils.js'

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
  
  
  export default class Palette extends Component {
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