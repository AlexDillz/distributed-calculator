package server

import (
	"net"
	"os"

	"github.com/AlexDillz/distributed-calculator/internal/agent"
	pb "github.com/AlexDillz/distributed-calculator/internal/proto"
	"google.golang.org/grpc"
)

func RunGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, &agent.Worker{})
	os.Stdout.WriteString("GRPC Server running on :50051\n")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
