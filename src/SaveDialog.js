import React, { Component } from 'react'
import { ModalContainer, ModalDialog } from 'react-modal-dialog'

export default class SaveDialog extends Component {
    constructor(props) {
        super(props)
        this.state = {
            name: ""
        }   
    }

    render() {
        return <ModalContainer onClose={this.props.onClose}>
        <ModalDialog onClose={this.props.onClose}>
          <h1>Save</h1>
          <form onSubmit={(event) => {
              this.props.onSave(this.state.name)
              event.preventDefault()}}>
            <input 
                value={this.state.name}
                onChange={(event) => this.setState({name: event.target.value})}/>
            <input
                type="submit" value="Save"/>
          </form>
        </ModalDialog>
      </ModalContainer>
    }
}
