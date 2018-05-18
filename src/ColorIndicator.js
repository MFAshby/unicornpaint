import React, { Component } from 'react'
import { rgb } from './Utils'

function toHex(n) {
    var s = Number(n).toString(16)
    if (s.length === 1) {
        s = "0" + s
    }
    return s
}

export default class ColorIndicator extends Component {
    render() {
        let { r, g, b } = rgb(this.props.color)
        let colorDesc = `#${toHex(r)}${toHex(g)}${toHex(b)}`
        var foreground = "black"
        if (r < 133 && g < 133 && b < 133) {
            foreground = "white"
        }

        return <div 
            style={{
                background: `rgb(${r},${g},${b})`, 
                color: foreground,
                ...styles.colorIndicator
            }}>
            <span>{colorDesc}</span>
        </div>
    }
}

const styles = {
    colorIndicator: {
        width: "100px",
        height: "100px"
    }
}