package auth

import (
	"context"
	"log"

	pb "github.com/daioru/marketplace/internal/generated/api/proto"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Println("Register called for:", req.Email)
	return &pb.RegisterResponse{UserId: "12345"}, nil
}

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Println("Login called for:", req.Email)
	return &pb.LoginResponse{Token: "mocked-jwt-token"}, nil
}

func (s *AuthServiceServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	log.Println("ValidateToken called for token:", req.Token)
	return &pb.ValidateTokenResponse{Valid: true, UserId: "12345"}, nil
}
