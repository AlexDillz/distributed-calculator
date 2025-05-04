package agent

import (
	"io"
	"testing"

	"context"
	"net"

	pb "github.com/AlexDillz/distributed-calculator/internal/proto"
	"github.com/AlexDillz/distributed-calculator/pkg/calculation"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const buferSize = 1024 * 1024

func dialer() func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(buferSize)
	srv := grpc.NewServer()
	pb.RegisterTaskServiceServer(srv, &Agent{})
	go srv.Serve(lis)
	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestAgentExecuteTask(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(dialer()),
		grpc.WithInsecure(),
	)
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	stream, err := client.ExecuteTask(ctx)
	assert.NoError(t, err)

	exprs := []string{"2+3", "4*5", "10/2"}
	for _, e := range exprs {
		assert.NoError(t, stream.Send(&pb.TaskRequest{Expression: e}))
	}
	assert.NoError(t, stream.CloseSend())

	for _, e := range exprs {
		resp, err := stream.Recv()
		assert.NoError(t, err)
		exp, _ := calculation.Calc(e)
		assert.Equal(t, exp, resp.Result)
		assert.Empty(t, resp.Error)
	}
	_, err = stream.Recv()
	assert.Equal(t, io.EOF, err)
}
