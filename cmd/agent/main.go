package main

import (
	"net"

	"github.com/AlexDillz/distributed-calculator/internal/agent"
	"github.com/AlexDillz/distributed-calculator/internal/config"
	pb "github.com/AlexDillz/distributed-calculator/internal/proto"
	"github.com/AlexDillz/distributed-calculator/pkg/logging"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	logging.InitLogger()
	logger := logging.GetLogger()

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, &agent.Agent{})

	logger.Printf("Agent gRPC Server running on %s\n", cfg.GRPCPort)
	grpcServer.Serve(lis)
}
