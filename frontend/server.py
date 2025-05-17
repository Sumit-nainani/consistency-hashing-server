import sys
sys.path.append("/Users/sumitnainani/Desktop/goserver/frontend")
import asyncio
import tornado.ioloop
import tornado.web
import tornado.websocket
from tornado import escape
import websockets
import json
from google.protobuf.json_format import MessageToDict
from google.protobuf.json_format import MessageToJson
# from grpc.grpcClient import fetch_data_from_grpc
from protopy import hashing_pb2
from grpcClient.grpcClient import fetch_data_from_grpc

clients = set()


    
class MainHandler(tornado.web.RequestHandler):
    def get(self):
        self.render("static/index.html")

class WebClientSocket(tornado.websocket.WebSocketHandler):
    def open(self):
        clients.add(self)
        initial_data = fetch_data_from_grpc()
        json_data = MessageToJson(initial_data)
        print(json_data,"json_data")
        self.write_message(json_data)
        print("Frontend connected")

    def on_close(self):
        clients.remove(self)
        print("Frontend disconnected")

async def consume_from_go_ws():
    uri = "ws://localhost:8085/ws"
    async for ws in websockets.connect(uri):
        try:
            async for raw_data in ws:
                proto_msg = hashing_pb2.WebSocketMetadata()
                proto_msg.ParseFromString(raw_data)
                json_data = MessageToDict(proto_msg)
                print(json_data,"data")
                # Send to connected frontend clients
                for client in clients:
                    print("sending data")
                    client.write_message(json.dumps(json_data))
        except Exception as e:
            print("WebSocket connection failed:", e)
            await asyncio.sleep(5)

def make_app():
    return tornado.web.Application([
        (r"/", MainHandler),
        (r"/ws-client", WebClientSocket),
        (r"/static/(.*)", tornado.web.StaticFileHandler, {"path": "./static"}),
    ])

if __name__ == "__main__":
    app = make_app()
    app.listen(8888)
    loop = asyncio.get_event_loop()
    loop.create_task(consume_from_go_ws())
    tornado.ioloop.IOLoop.current().start()
