const NO_OP = 'NO_OP'
const SET_PIXEL = 'SET_PIXEL'
const CLEAR = 'CLEAR'
const SAVE = 'SAVE'
const LOAD = 'LOAD'
const ADD_FRAME = "ADD_FRAME"
const REMOVE_FRAME = "REMOVE_FRAME"

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

function setPixel(websocket, x, y, r, g, b, frame) {
  sendAction(websocket, {
    type: SET_PIXEL,
    x: x,
    y: y,
    r: r, 
    g: g,
    b: b,
    frame: frame
  })
}

function clear(websocket) {
  sendAction(websocket, { type: CLEAR })
}

function noop(websocket) {
  sendAction(websocket, { type: NO_OP })
}

function addFrame(websocket, frame = 1, delay = 50) {
  sendAction(websocket, { 
    type: ADD_FRAME, 
    frame: frame, 
    delay: delay
  })
}

function removeFrame(websocket, frame = 1) {
  sendAction(websocket, {
    type: REMOVE_FRAME, 
    frame: frame
  })
}

export {
    setPixel, 
    clear, 
    noop,
    save,
    load,
    addFrame,
    removeFrame
}