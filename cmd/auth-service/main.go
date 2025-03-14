package main

import (
	"log"
	"net"

	pb "github.com/daioru/marketplace/api/proto/auth"
	"github.com/daioru/marketplace/internal/auth"
	"github.com/daioru/marketplace/internal/auth/db"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
)

func main() {
	database := db.InitDB()

	err := goose.Up(database.DB, "migrations/auth")
	if err != nil {
		log.Fatalf("Auth migration error: %v", err)
	}

	server := grpc.NewServer()
	authService := &auth.AuthService{DB: database}
	pb.RegisterAuthServiceServer(server, authService)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Println("Auth-service started on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
