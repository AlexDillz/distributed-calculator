package main

import (
	"log"
	"net"

	"github.com/AlexDillz/distributed-calculator/internal/agent"
	pb "github.com/AlexDillz/distributed-calculator/internal/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterTaskServiceServer(srv, &agent.Agent{})
	log.Println("Agent is running on :50051")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
