// Package shared contains shared data between the host and plugins.
package shared

import (
	"context"

	"github.com/tiny-job/core/proto"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   PluginJobEnv,
	MagicCookieValue: PluginJob,
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	PluginJob: &JobRunPlugin{},
}

// Job 任务接口
type Job interface {
	Run(ctx context.Context, params map[string]string) (map[string]string, error)
}

// JobRunPlugin 任务运行插件
type JobRunPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Job
}

func (p *JobRunPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterJobServer(s, &GRPCServer{Impl: p.Impl, broker: broker})
	return nil
}

func (p *JobRunPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewJobClient(c), broker: broker}, nil
}
