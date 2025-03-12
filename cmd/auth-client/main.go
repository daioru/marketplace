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
	log.Println("Login response access", loginResp.AccessToken)
	log.Println("Login response refresh", loginResp.RefreshToken)

	validateResp, err := client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: loginResp.AccessToken,
	})
	if err != nil {
		log.Fatalf("ValidateToken failed: %v", err)
	}
	log.Println("ValidateToken response", validateResp.Valid, validateResp.UserId)
}
