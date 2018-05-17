const NO_OP = 'NO_OP'
const SET_PIXEL = 'SET_PIXEL'
const CLEAR = 'CLEAR'
const SAVE = 'SAVE'
const LOAD = 'LOAD'

function sendAction(websocket, action) {
  let actionStr = JSON.stringify(action)
  websocket.send(actionStr)
}

function save(websocket, saveName) {
  sendAction(websocket, {
    type: SAVE,
    saveName: saveName
  })
}

function load(websocket, saveName) {
  sendAction(websocket, {
    type: LOAD,
    saveName: saveName
  })
}

function setPixel(websocket, x, y, r, g, b) {
  sendAction(websocket, {
    type: SET_PIXEL,
    x: x,
    y: y,
    r: r, 
    g: g,
    b: b
  })
}

function clear(websocket) {
  sendAction(websocket, { type: CLEAR })
}

function noop(websocket) {
  sendAction(websocket, { type: NO_OP })
}

export {
    setPixel, 
    clear, 
    noop,
    save,
    load
}