package main

import (
	"log"
	"net"

	"github.com/daioru/marketplace/internal/auth"
	"google.golang.org/grpc"

	pb "github.com/daioru/marketplace/internal/generated/api/proto"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &auth.AuthServiceServer{})

	log.Println("Auth service is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
