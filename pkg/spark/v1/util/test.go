package util

import (
	"fmt"
	gnats "github.com/nats-io/nats-server/v2/server"
	gnatsTest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"net"
)

func RunServerOnPort(port int, dir string) (*gnats.Server, error) {
	opts := gnatsTest.DefaultTestOptions
	opts.Port = port
	opts.JetStream = true
	opts.StoreDir = dir
	return RunServerWithOptions(&opts)
}

func RunServerWithOptions(opts *gnats.Options) (*gnats.Server, error) {
	return gnats.NewServer(opts)
}

// GetFreeTCPPort returns free open TCP port
func GetFreeTCPPort() (port int, err error) {
	ln, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		return 0, err
	}
	port = ln.Addr().(*net.TCPAddr).Port
	err = ln.Close()
	return
}

func GetNatsClient(port int) (*nats.Conn, jetstream.JetStream) {
	sUrl := fmt.Sprintf("nats://127.0.0.1:%d", port)
	nc, err := nats.Connect(sUrl)
	if err != nil {
		panic(err)
	}

	if !nc.IsConnected() {
		errorMsg := fmt.Errorf("could not establish connection to nats-server")
		panic(errorMsg)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		panic(err)
	}

	return nc, js
}
