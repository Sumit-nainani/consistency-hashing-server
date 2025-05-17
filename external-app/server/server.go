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

func (ns *NodeService) GetHashRingData(ctx context.Context, req *pb.Empty) (*pb.WebSocketMetadataList, error) {
	WebSocketMetadataList := &pb.WebSocketMetadataList{}
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
	// server := grpc.NewServer()
	// nodeService := &NodeService{HashRing: hashRing}
	// pb.RegisterNodeServer(server, nodeService)

	// Wrap gRPC server with grpc-web
	// wrappedGrpc := grpcweb.WrapServer(server,
	// 	grpcweb.WithOriginFunc(func(origin string) bool { return true }), // Allow all origins
	// )

	// Create HTTP handler for grpc-web
	// httpServer := http.Server{
	// 	Addr: ":8080",
	// 	Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
	// 		if wrappedGrpc.IsGrpcWebRequest(req) || wrappedGrpc.IsAcceptableGrpcCorsRequest(req) {
	// 			wrappedGrpc.ServeHTTP(resp, req)
	// 		} else {
	// 			resp.WriteHeader(http.StatusNotFound)
	// 		}
	// 	}),
	// }

	// log.Println("gRPC-Web Server listening on :8080")
	// if err := httpServer.ListenAndServe(); err != nil {
	// 	log.Fatalf("Failed to serve HTTP: %v", err)
	// }
}
