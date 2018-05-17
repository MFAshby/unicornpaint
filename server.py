import json
from flask import Flask, send_from_directory
from flask_sockets import Sockets
from gevent import pywsgi
from geventwebsocket.handler import WebSocketHandler
import os
import os.path
import pickle

try:
    import unicornhathd as unicorn
    print("unicorn hat hd detected")
except ImportError:
    from unicorn_hat_sim import unicornhathd as unicorn

# Actions that can be sent 
NO_OP = 'NO_OP'
SET_PIXEL = 'SET_PIXEL'
CLEAR = 'CLEAR'
SAVE = 'SAVE'
LOAD = 'LOAD'

# GLobal references to app & clients 
all_clients = set()
app = Flask(__name__, static_folder='')
sockets = Sockets(app)

def send_state():
    global all_clients
    state = {
        "saves": os.listdir(SAVES_DIR),
        "pixels": unicorn.get_pixels().tolist()
    }
    stateJson = json.dumps(state)
    for client in all_clients:
        client.send(stateJson)

def execute_command(command):
    cmd_type = command['type']
    if cmd_type == NO_OP: # Used just to get state without doing an action
        pass 
    elif cmd_type == SET_PIXEL:
        x = int(command['x'])
        y = int(command['y'])
        r = int(command['r'])
        g = int(command['g'])
        b = int(command['b'])
        unicorn.set_pixel(x, y, r, g, b)
        unicorn.show()
    elif cmd_type == CLEAR:
        unicorn.clear()
        unicorn.show()
    elif cmd_type == SAVE:
        saveName = command["saveName"]
        save(saveName)
    elif cmd_type == LOAD:
        saveName = command["saveName"]
        load(saveName)

# Constants
SAVES_DIR = 'saves/'

def save(saveName):
    with open(os.path.join(SAVES_DIR, saveName), "wb") as f:
        pickle.dump(unicorn.get_pixels(), f)

def load(saveName):
    with open(os.path.join(SAVES_DIR, saveName), "rb") as f:
        pixels = pickle.load(f)
        for x, row in enumerate(pixels):
            for y, pixel in enumerate(row):
                unicorn.set_pixel(x, y, *pixel)
        unicorn.show()

@sockets.route('/ws')
def do_websocket(websocket):
    global all_clients
    all_clients.add(websocket)
    try:
        while not websocket.closed:
            command = websocket.receive()
            if not command: 
                break
            cmd = json.loads(command)
            execute_command(cmd)
            send_state()
    finally:
        all_clients.remove(websocket)


@app.route('/<path:path>')
def send_static(path):
    return send_from_directory('build/', path)

def main():
    try:
        os.mkdir(SAVES_DIR)
    except:
        pass
    try:
        print("Serving on port 3001")
        server = pywsgi.WSGIServer(('', 3001), app, handler_class=WebSocketHandler)
        server.serve_forever()
    finally:
        unicorn.off()

if __name__=="__main__":
    main()
    
