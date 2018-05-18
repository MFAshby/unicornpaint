import React, { Component } from 'react'

export default class ConnectedIndicator extends Component {
  render() {
    let connectedText = this.props.connected ? "Connected" : "Not connected"
    let color = this.props.connected ? "green" : "red"
    return <div><span style={{
      background:color,
      display: "inline-block"
    }}>{connectedText}</span></div>
  }
}