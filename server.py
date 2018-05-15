import json
from flask import Flask, send_from_directory
from flask_sockets import Sockets
from gevent import pywsgi
from geventwebsocket.handler import WebSocketHandler

try:
    import unicornhathd as unicorn
    print("unicorn hat hd detected")
except ImportError:
    from unicorn_hat_sim import unicornhathd as unicorn

NO_OP = 'NO_OP'
SET_PIXEL = 'SET_PIXEL'
CLEAR = 'CLEAR'

all_clients = set()

app = Flask(__name__, static_folder='build/')
sockets = Sockets(app)

def send_state():
    global all_clients
    jsonPixels = json.dumps(unicorn.get_pixels())
    for client in all_clients:
        client.send(jsonPixels)

def execute_command(command):
    cmd_type = command['type']
    if cmd_type == NO_OP:
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

@sockets.route('/ws')
def do_websocket(websocket):
    global all_clients
    all_clients.add(websocket)
    try:
        while not websocket.closed:
            command = websocket.receive()
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
        server = pywsgi.WSGIServer(('', 3001), app, handler_class=WebSocketHandler)
        server.serve_forever()
    except:
        unicorn.off()

if __name__=="__main__":
    main()