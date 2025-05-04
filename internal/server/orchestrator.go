package server

import (
	"context"
	"errors"
	"time"

	pb "github.com/AlexDillz/distributed-calculator/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const agentAddr = "localhost:50051"

func EvaluateExpression(expr string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.Dial(agentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	stream, err := client.ExecuteTask(ctx)
	if err != nil {
		return 0, err
	}

	if err := stream.Send(&pb.TaskRequest{Expression: expr}); err != nil {
		return 0, err
	}
	if err := stream.CloseSend(); err != nil {
		return 0, err
	}

	resp, err := stream.Recv()
	if err != nil {
		return 0, err
	}
	if resp.Error != "" {
		return 0, errors.New(resp.Error)
	}
	return resp.Result, nil
}
