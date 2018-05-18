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

        let style = Object.assign({},
          {background: `rgb(${r},${g},${b})`},
          styles.paletteItem,
          isSelected && styles.selected
        )

        return <div
          onClick={() => this.props.onSelectColor([r, g, b])}
          style={style} 
          key={r*10000+g*1000+b}/>
      })
      let list1 = paletteListItems.slice(0, paletteListItems.length / 2)
      let list2 = paletteListItems.slice(paletteListItems.length / 2)
      return (
        <div>
          {list1}<br/>
          {list2}
        </div>)
    }
  }

  const styles = {
    paletteItem: {
      width: "30px",
      height: "30px",
      border: "3px solid black",
      display: "inline-block",
    },
    selected: {
      border: "3px solid red"
    }
  }