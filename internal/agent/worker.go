package agent

import (
	"log"

	pb "github.com/AlexDillz/distributed-calculator/internal/proto"
	"github.com/AlexDillz/distributed-calculator/pkg/calculation"
)

type Worker struct {
	pb.UnimplementedTaskServiceServer
}

func (w *Worker) ExecuteTask(stream pb.TaskService_ExecuteTaskServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		log.Printf("Received expression: %s", req.Expression)

		result, calcErr := calculation.Calc(req.Expression)
		resp := &pb.TaskResponse{
			Result: result,
		}
		if calcErr != nil {
			resp.Error = calcErr.Error()
		}

		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
