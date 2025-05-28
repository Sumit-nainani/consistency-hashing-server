import asyncio,os
import tornado.ioloop
import tornado.web
import tornado.websocket
import websockets
from protopy import hashing_pb2
from dotenv import load_dotenv
from google.protobuf.json_format import MessageToJson
from grpcClient.grpcClient import fetch_data_from_grpc

load_dotenv()
clients = set()

class MainHandler(tornado.web.RequestHandler):
    def get(self):
        self.render("static/index.html")

class WebClientSocket(tornado.websocket.WebSocketHandler):
    def open(self):
        clients.add(self)
        self.write_message(MessageToJson(fetch_data_from_grpc()))
        print("Frontend connected")

    def on_close(self):
        clients.remove(self)
        print("Frontend disconnected")

# Websocket client for Golang Websocket server.
async def consume_from_go_ws():
    async for ws in websockets.connect(os.getenv('WEBSOCKET_URL')):
        try:
            async for raw_data in ws:
                proto_msg = hashing_pb2.WebSocketMetadata()
                proto_msg.ParseFromString(raw_data)
                for client in clients:
                    client.write_message(MessageToJson(proto_msg))
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
