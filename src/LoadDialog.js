import React, { Component } from 'react'
import { ModalContainer, ModalDialog } from 'react-modal-dialog'

export default class LoadDialog extends Component {
    render() {
        let savesListItems = this.props.saves.map((save) => {
            return <a key={save} onClick={() => this.props.onLoad(save)}>
                <li>{save}</li>
            </a>
        })

        return <ModalContainer onClose={this.props.onClose}>
            <ModalDialog onClose={this.props.onClose}>
            <h1>Load</h1>
            <ul>
            {savesListItems}
            </ul>
            </ModalDialog>
        </ModalContainer>
    }
}