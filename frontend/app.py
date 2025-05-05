import streamlit as st
import numpy as np
import matplotlib.pyplot as plt
import asyncio
import threading
import websockets
import json
import hashlib
import os
from google.protobuf.json_format import MessageToDict
import hashing_pb2  # Your compiled proto file
from streamlit_autorefresh import st_autorefresh
# Constants
TOTAL_POSITIONS = 64
CIRCLE_RADIUS = 250
STATE = {"servers": []}  # Shared state

# Lock for thread-safe state access
state_lock = threading.Lock()


# Get coordinates on the ring
def get_coordinates(index):
    angle = (2 * np.pi * index) / TOTAL_POSITIONS
    x = CIRCLE_RADIUS * np.cos(angle)
    y = CIRCLE_RADIUS * np.sin(angle)
    return x, y


# WebSocket consumer function
async def consume_data():
    uri = "ws://localhost:8085/ws"
    try:
        async with websockets.connect(uri) as websocket:
            while True:
                raw_data = await websocket.recv()
                msg = hashing_pb2.WebSocketMetadata()
                msg.ParseFromString(raw_data)
                data = MessageToDict(msg)
                with state_lock:
                    print(data,"inside")
                    if data["type"] == "pod" and data["action"] == "add":
                        node = data["nodeMetaData"]
                        STATE["servers"].append({
                            "name": node["nodeName"],
                            "ip": node["nodeIp"],
                            "hash": int(node["nodeHash"])
                        })

    except Exception as e:
        print("WebSocket error:", e)


# Thread runner
def start_websocket_client():
    asyncio.new_event_loop().run_until_complete(consume_data())


# Draw the ring with server positions
def draw_ring():
    fig, ax = plt.subplots(figsize=(6, 6), subplot_kw={'polar': True})
    ax.set_xticks(np.linspace(0, 2 * np.pi, TOTAL_POSITIONS, endpoint=False))
    ax.set_yticklabels([])
    ax.set_xticklabels([str(i) for i in range(TOTAL_POSITIONS)], fontsize=7)
    ax.set_theta_zero_location("N")
    ax.set_theta_direction(-1)

    # Plot border
    angles = np.linspace(0, 2 * np.pi, TOTAL_POSITIONS, endpoint=False)
    ax.plot(angles, [CIRCLE_RADIUS] * TOTAL_POSITIONS, 'o', color='gray', markersize=4)

    with state_lock:
        for server in STATE.get("servers", []):
            index = server["hash"]
            angle = angles[index]
            ax.plot(angle, CIRCLE_RADIUS, 'o', color='blue', markersize=10)
            ax.text(angle, CIRCLE_RADIUS + 20, server["name"][:5], ha='center', fontsize=7, color='blue')

    ax.set_title("Consistent Hashing Ring", va='bottom')
    st.pyplot(fig)


# Streamlit app
st.set_page_config(layout="centered")
st.title("Consistent Hashing Ring (Live View)")

st_autorefresh(interval=1000, key="datarefresh")

# Start WebSocket client once
if "websocket_started" not in st.session_state:
    threading.Thread(target=start_websocket_client, daemon=True).start()
    st.session_state.websocket_started = True

# Refresh button
# if st.button("🔁 Refresh View"):
#     st.rerun()

# Draw the ring
draw_ring()
