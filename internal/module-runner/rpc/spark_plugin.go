//go:build exclude

package module_runner

import (
	"github.com/rs/zerolog/log"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Blackboard interface {
	SetValue(data any) error
}

// SparkRpcApi is the interface that we're exposing as a plugin.
type SparkRpcApi interface {
	Greet(bb Blackboard) string
}

// Here is an implementation that talks over RPC
type sparkRpc struct {
	client *rpc.Client
	b      *plugin.MuxBroker
}

type receiver struct {
	bb Blackboard
}

func (r *receiver) SetValue(args string, reply *any) error {
	err := r.bb.SetValue(args)
	log.Info().Msgf("%v: %v", args, reply)
	return err
}

func (g *sparkRpc) Greet(bb Blackboard) string {
	var resp string

	reqId := g.b.NextId()
	go g.b.AcceptAndServe(reqId, &receiver{bb: bb})

	err := g.client.Call("Plugin.Greet", reqId, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

type GreeterRPCServer struct {
	// This is the real implementation
	Impl SparkRpcApi
	mux  *plugin.MuxBroker
}

func (s *GreeterRPCServer) Greet(i uint32, resp *string) error {
	log.Info().Msgf("got responder: %d", i)
	conn, err := s.mux.Dial(i)
	if err != nil {
		return err
	}

	bc := &blackboardClient{rpcClient: rpc.NewClient(conn)}

	*resp = s.Impl.Greet(bc)
	return nil
}

type blackboardClient struct {
	rpcClient *rpc.Client
}

func (b blackboardClient) SetValue(data any) error {
	var nr any
	if err := b.rpcClient.Call("Plugin.SetValue", data, &nr); err != nil {
		return err
	}

	log.Info().Msgf("%v", nr)
	return nil
}

type SparkPlugin struct {
	// Impl Injection
	Impl SparkRpcApi
}

func (p *SparkPlugin) Server(mux *plugin.MuxBroker) (interface{}, error) {
	return &GreeterRPCServer{Impl: p.Impl, mux: mux}, nil
}

func (*SparkPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &sparkRpc{client: c, b: b}, nil
}

type IBlackboard struct {
	Value  string
	GetVal func() string
}

func (b *IBlackboard) SetValue(data any) error {
	var ov string
	if b.GetVal != nil {
		ov = b.GetVal()
	}

	log.Info().Msgf("IBlackboard (%s) received (%s): %v", b.Value, ov, data)
	return nil
}
