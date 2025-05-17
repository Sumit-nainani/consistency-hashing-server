import grpc
from protopy import hashing_pb2,hashing_pb2_grpc

def fetch_data_from_grpc():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = hashing_pb2_grpc.NodeStub(channel)
        response = stub.GetHashRingData(hashing_pb2.Empty())
        print(response,"response")
        return response