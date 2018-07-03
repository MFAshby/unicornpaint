import React, { Component } from 'react'
import { rgb, getPixel } from './Utils'

class FramePreview extends Component {
    render() {
        let { pixels, selected, index, onClick } = this.props
        let width = pixels.length
        let height = pixels[0].length

        let cells = []
        for (var y=height-1; y>=0; y--) {
            for (var x=0; x<width; x++) {
                let {r, g, b} = rgb(getPixel(x, y, pixels))
                cells.push(<div
                    key={`FramePreview ${index} ${x} ${y}`}
                    className="framePreviewCell"
                    style={{
                        background: `rgb(${r},${g},${b})`,
                      }}/>)
            }
        } 

        let wrapperClassName = selected ? "framePreviewWrapper selected" : "framePreviewWrapper"
        return <div 
            className={wrapperClassName}
            onClick={onClick}>
            <span>Frame {index + 1}</span>
            <div className="framePreview">
            {cells}
            </div>
        </div>
    }
}

function animationPreview(props) {
    let { imageData } = props
    let imgUrl = `url(data:image/gif;base64,${imageData})`
    return <div className="animationPreviewWrapper">
        <span>Preview:</span>
        <div className="animationPreview" 
            style={{backgroundImage: imgUrl}}/>
    </div>
}

export default class FrameControl extends Component {
    render() {
        // A series of divs, 1 per frame, 
        let { frames, 
            selectedFrame,
            onFrameSelected,
            imageData } = this.props

        let framePreviews = frames.map((frame, ix) => 
            <FramePreview 
                key={ix}
                index={ix}
                pixels={frame} 
                selected={ix === selectedFrame}
                onClick={() => onFrameSelected(ix)}/>)

        

        return <div className="frameControl">
            {animationPreview({imageData:imageData})}
            {framePreviews}
        </div>
    }
}
