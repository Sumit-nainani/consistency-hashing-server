package server

import (
	"context"
	"fmt"
	ph "hashing/hashing"
	"hashing/hashring"
	"log"
	"net"

	"google.golang.org/grpc"
)

type NodeService struct {
	HashRing *hashring.HashRing
	ph.UnimplementedNodeServer
}

// GetNodeForRequest returns the node responsible for the given gRPC request
func (ns *NodeService) GetNodeForRequest(ctx context.Context, req *ph.NodeRequest) (*ph.NodeResponse, error) {
	fmt.Println("hellooooo")
	node, hash := ns.HashRing.GetNode(req.Key)
	fmt.Println(node, "node")
	return &ph.NodeResponse{Node: node, Hash: int32(hash)}, nil
}

// AddNodeForRequest adds a new node for the given gRPC request
func (ns *NodeService) AddNodeForRequest(ctx context.Context, req *ph.NodeRequest) (*ph.AddNodeResponse, error) {
	hashValue := ns.HashRing.AddNode(req.Node, req.Ip)
	return &ph.AddNodeResponse{Hash: int32(hashValue)}, nil
}

// RemoveNodeForRequest removes a node for the given gRPC request
func (ns *NodeService) RemoveNodeForRequest(ctx context.Context, req *ph.NodeRequest) (*ph.DeleteNodeResponse, error) {
	ns.HashRing.RemoveNode(req.Node)
	return &ph.DeleteNodeResponse{}, nil
}

func StartGrpcServer() {

	hashRing := hashring.GetRingInstance()

	server := grpc.NewServer()
	nodeService := &NodeService{HashRing: hashRing}
	ph.RegisterNodeServer(server, nodeService)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	log.Println("Server listening on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
