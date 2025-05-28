package server

import (
	"context"
	pb "hashing/hashing"
	"hashing/hashring"
	"log"
	"net"

	"google.golang.org/grpc"
)

type NodeService struct {
	HashRing *hashring.HashRing
	pb.UnimplementedNodeServer
}

// This is RPC service which is registered in proto file.
// It is used for fetching the initial data of availbale nodes and registered clients.
// This service is called from frontend using a RPC python client and called first time when web interface is rendered.
// This service passes a list of nodes and clients in the protobuff structure registered in proto file.
func (ns *NodeService) GetHashRingData(ctx context.Context, req *pb.Empty) (*pb.WebSocketMetadataList, error) {
	WebSocketMetadataList := &pb.WebSocketMetadataList{}

	// Inserting node metadatas.
	for _, node_meta_data := range hashring.GetRingInstance().NodeNameToNodeMetaData {
		WebSocketMetadataList.Item = append(WebSocketMetadataList.Item, &pb.WebSocketMetadata{
			Type:   "pod",
			Action: "add",
			Data: &pb.WebSocketMetadata_NodeMetaData{
				NodeMetaData: &pb.NodeMetaData{
					NodeHash: node_meta_data.NodeHash,
					NodeIp:   node_meta_data.NodeIP,
					NodeName: node_meta_data.NodeName,
				},
			},
		})
	}
    
	// Inserting client metadatas.
	for _, request_meta_data := range hashring.GetRingInstance().RequestIpToMetaData {
		WebSocketMetadataList.Item = append(WebSocketMetadataList.Item, &pb.WebSocketMetadata{
			Type: "client",
			Data: &pb.WebSocketMetadata_RequestMetaData{
				RequestMetaData: &pb.RequestMetaData{
					AssignedNodeName: request_meta_data.AssignedNodeName,
					AssignedNodeIp:   request_meta_data.AssignedNodeIP,
					RequestHash:      request_meta_data.RequestHash,
					AssignedNodeHash: request_meta_data.AssignedNodeHash,
				},
			},
		})
	}

	return WebSocketMetadataList, nil
}

// This method is starting gRPC server in a separate goroutine.
func StartGrpcServer() {
	hashRing := hashring.GetRingInstance()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterNodeServer(grpcServer, &NodeService{HashRing: hashRing})
	log.Println("gRPC server on :50051")
	grpcServer.Serve(lis)
}
