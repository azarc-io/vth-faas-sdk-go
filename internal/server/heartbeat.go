package server

import (
	"context"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
	"time"
)

type Heartbeat struct {
	config *config.Config
	client sdk_v1.ManagerServiceClient
	ticker *time.Ticker
	done   chan struct{}
}

func NewHeartbeat(config *config.Config, client sdk_v1.ManagerServiceClient) *Heartbeat {
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
				_, _ = h.client.RegisterHeartbeat(context.Background(), &sdk_v1.RegisterHeartbeatRequest{AgentId: h.config.App.InstanceId}) //TODO handle error
			}
		}
	}()

}

func (h *Heartbeat) Stop() {
	h.ticker.Stop()
	h.done <- struct{}{}
}
