package sdk_v1

import (
	"context"
	"time"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
)

type Heartbeat struct {
	config *config.Config
	client ManagerServiceClient
	ticker *time.Ticker
	done   chan struct{}
}

func NewHeartbeat(config *config.Config, client ManagerServiceClient) *Heartbeat {
	return &Heartbeat{config: config, client: client}
}

func (h *Heartbeat) Start() {
	h.ticker = time.NewTicker(h.config.ManagerService.HeartBeatInterval)
	h.done = make(chan struct{})

	go func() {
		for {
			select {
			case <-h.done:
				return
			case <-h.ticker.C:
				_, _ = h.client.RegisterHeartbeat(context.Background(), &RegisterHeartbeatRequest{AgentId: h.config.App.InstanceID}) // TODO handle error
			}
		}
	}()
}

func (h *Heartbeat) Stop() {
	h.ticker.Stop()
	h.done <- struct{}{}
}
