package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "github.com/daioru/marketplace/internal/generated/api/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	registerResp, err := client.Register(context.Background(), &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Register failed: %v", err)
	}
	log.Println("Register response", registerResp.UserId)

	loginResp, err := client.Login(context.Background(), &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	log.Println("Login response", loginResp.Token)

	validateResp, err := client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: loginResp.Token,
	})
	if err != nil {
		log.Fatalf("ValidateToken failed: %v", err)
	}
	log.Println("ValidateToken response", validateResp.Valid, validateResp.UserId)
}
