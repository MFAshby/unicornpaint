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
        let { r, g, b } = rgb(item)
        let selected = rgb(this.props.selectedColor)
        let isSelected = r === selected.r && g === selected.g && b === selected.b

        let className = isSelected ? "paletteItem selected" : "paletteItem"

        return <div
          key={r*10000+g*1000+b}
          className={className}
          style={{background: `rgb(${r},${g},${b})`}} 
          onClick={() => this.props.onSelectColor([r, g, b])}/>
      })
      return (
        <div className="palette">
          {paletteListItems}
        </div>)
    }
  }