import React, { Component } from 'react'
import './fontawesome-all.min.js'

class ToolKitItem extends Component {
  render() {
    let iconClass = this.props.tool.icon + " fa-2x" 
    var divClass = "toolkititem"
    if (this.props.selected) {
      divClass += " selected"
    }
    return <div className={divClass}
        onClick={() => this.props.onSelectTool(this.props.tool)}>
        <i className={iconClass}></i>
      </div>
  }
}

export default class Toolkit extends Component {
  render() {
    let toolComponents = this.props.tools.map((tool) => {
      return <ToolKitItem 
        tool={tool} 
        key={tool.name}
        selected={tool === this.props.selectedTool}
        onSelectTool={this.props.onSelectTool}/>
    })
    return <div className="toolkit">
      {toolComponents}
    </div>
  }
}
