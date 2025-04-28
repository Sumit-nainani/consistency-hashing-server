package client

import (
	"log"
	"sync"

	ph "hashing/hashing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ClientInstance ph.NodeClient
	connInstance   *grpc.ClientConn
	once           sync.Once
)

func InitGrpcClinet(){
	// Connect to the gRPC server
	once.Do(func(){
		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		connInstance = conn
		ClientInstance = ph.NewNodeClient(conn)
	})
}
	
func GetClient() ph.NodeClient{
	InitGrpcClinet()
	return ClientInstance
}

func Close() {
	if connInstance != nil {
		connInstance.Close()
	}
}
