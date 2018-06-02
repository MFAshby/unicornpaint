import React, { Component } from 'react'

class ToolKitItem extends Component {
  render() {
    let iconClass = this.props.tool.icon + " fa-2x" 
    let style = styles.toolkititem
    if (this.props.selected) {
      style = Object.assign({}, style, styles.selected)
    }
    return <div style={style}
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
    return <div style={styles.toolkit}>
      {toolComponents}
    </div>
  }
}

const styles = {
  toolkit: {
    background: "grey"
  },
  toolkititem: {
    padding: "5px",
  },
  selected: {
    background: "lightcoral"
  }
}