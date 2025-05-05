import asyncio
import websockets
import hashing_pb2  # Replace with the actual name of the generated file (e.g., file_pb2)
from google.protobuf.json_format import MessageToJson

# WebSocket client to consume data from the server
async def consume_data():
    uri = "ws://localhost:8085/ws"  # WebSocket server URI

    async with websockets.connect(uri) as websocket:
        while True:
            # Receive the Protobuf message from the WebSocket server
            raw_data = await websocket.recv()
            
            # Deserialize the Protobuf data
            metadata_list = hashing_pb2.WebSocketMetadata()
            metadata_list.ParseFromString(raw_data)
            
            # Convert the Protobuf object to a JSON string for human readability
            readable_data = MessageToJson(metadata_list)
            
            # Print the readable data
            print("Received data:", readable_data)

# Run the WebSocket client
asyncio.get_event_loop().run_until_complete(consume_data())
