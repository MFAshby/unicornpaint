import React, { Component } from 'react'
import { rgb, getPixel } from './Utils'

class FramePreview extends Component {
    render() {
        let width = this.props.pixels.length
        let height = this.props.pixels[0].length

        let rows = []
        for (var y=height-1; y>=0; y--) {
            let cells = []
            for (var x=0; x<width; x++) {
                let {r, g, b} = rgb(getPixel(x, y, this.props.pixels))
                cells.push(<td
                    style={{
                        background: `rgb(${r},${g},${b})`,
                        ...styles.previewPixel
                      }}/>)
            }
            rows.push(<tr key={y}>{cells}</tr>)
        } 

        let bgColor = this.props.selected ? "red" : "grey"
        return <table onClick={this.props.onClick} style={{background: bgColor}}><tbody> {rows} </tbody></table>
    }
}

export default class FrameControl extends Component {
    // frames: []
    // selectedFrame: number
    render() {
        // A series of divs, 1 per frame, 
        let frames = this.props.frames 
        let selectedFrame = this.props.selectedFrame 
        let framePreviews = frames.map((frame, ix) => 
            <div>
                <span>Frame {ix+1}</span>
                <FramePreview 
                    pixels={frame} 
                    selected={ix === this.props.selectedFrame}
                    onClick={() => this.props.onFrameSelected(ix)}/>
            </div>)
        return <div style={styles.previewContainer}>{framePreviews}</div>
    }
}

const styles = {
    previewcontainer: {
    },
    previewTable: {
    },
    previewPixel: {
        width: "4px",
        height: "4px",
    }
}