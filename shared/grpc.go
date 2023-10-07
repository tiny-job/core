package shared

import (
	"context"

	"github.com/tiny-job/core/proto"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct {
	client proto.JobClient
}

func (m *GRPCClient) Run(ctx context.Context, params map[string]string) (map[string]string, error) {
	resp, err := m.client.Run(ctx, &proto.RunRequest{
		Params: params,
	})
	return resp.Result, err
}

// GRPCServer Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	proto.UnimplementedJobServer
	// This is the real implementation
	Impl Job
}

func (m *GRPCServer) Run(ctx context.Context, req *proto.RunRequest) (*proto.RunResponse, error) {
	v, err := m.Impl.Run(ctx, req.GetParams())
	return &proto.RunResponse{Result: v}, err
}
