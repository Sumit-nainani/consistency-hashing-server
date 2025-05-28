import grpc,os
from protopy import hashing_pb2,hashing_pb2_grpc
from dotenv import load_dotenv

load_dotenv()
# gRPC python client for fetching initial node and client data.
def fetch_data_from_grpc():
    with grpc.insecure_channel(os.getenv("GRPC_SERVER_URL")) as channel:
        response = hashing_pb2_grpc.NodeStub(channel).GetHashRingData(hashing_pb2.Empty())
        return response