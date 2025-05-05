package agent

import (
	"context"
	"net"
	"testing"

	pb "github.com/AlexDillz/distributed-calculator/internal/proto"
	"github.com/AlexDillz/distributed-calculator/pkg/calculation"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func dialer() func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	pb.RegisterTaskServiceServer(srv, &Agent{})
	go srv.Serve(lis)
	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}
}

func setupClient(t *testing.T) pb.TaskServiceClient {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(dialer()),
		grpc.WithInsecure(),
	)
	assert.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	return pb.NewTaskServiceClient(conn)
}

func TestAgent_HappyPath(t *testing.T) {
	client := setupClient(t)
	ctx := context.Background()
	stream, err := client.ExecuteTask(ctx)
	assert.NoError(t, err)

	exprs := []string{
		"2+3",
		"4*5",
		"10/2",
		"(2+3)*4",
		"-5+10",
		"2e3+100",
	}
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
}

func TestAgent_ErrorCases(t *testing.T) {
	client := setupClient(t)
	ctx := context.Background()
	stream, err := client.ExecuteTask(ctx)
	assert.NoError(t, err)

	tests := []struct {
		expr        string
		expectError string
	}{
		{"2++2", "invalid expression"},
		{"5/0", "division by zero"},
		{"(2+3", "invalid expression"},
		{"", "invalid expression"},
	}
	for _, tc := range tests {
		assert.NoError(t, stream.Send(&pb.TaskRequest{Expression: tc.expr}))
		resp, err := stream.Recv()
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Error)
		assert.Contains(t, resp.Error, tc.expectError)
	}
	assert.NoError(t, stream.CloseSend())
}
