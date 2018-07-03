import React, { Component } from 'react'

export default class YoutubeEmbed extends Component {
    shouldComponentUpdate() {
        return false
    }
    render() {
        return <iframe 
            title="live stream"
            width="350" 
            height="350" 
            src="https://www.youtube.com/embed/HFlAhWwCdfE?autoplay=1" 
            frameBorder="0" 
            allow="autoplay; encrypted-media" 
            allowFullScreen></iframe>
    }
}