package server

import (
	"context"
	"log"
	"sync"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

type Callbacks struct {
	Signal   chan struct{}
	Fetches  int
	Requests int
	mu       sync.Mutex
}

func (cb *Callbacks) Report() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	log.Printf("server callbacks fetches=%d requests=%d\n", cb.Fetches, cb.Requests)
}
func (cb *Callbacks) OnStreamOpen(_ context.Context, id int64, typ string) error {

	log.Printf("stream %d open for %s\n", id, typ)

	return nil
}
func (cb *Callbacks) OnStreamClosed(id int64) {

	log.Printf("stream %d closed\n", id)

}
func (cb *Callbacks) OnStreamRequest(int64, *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Requests++
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	return nil
}

func (cb *Callbacks) OnStreamResponse(id int64, request *discovery.DiscoveryRequest, response *discovery.DiscoveryResponse) {

	log.Println("request:" + request.String())

	log.Println("response:" + response.String())

}
func (cb *Callbacks) OnFetchRequest(_ context.Context, req *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Fetches++
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	return nil
}
func (cb *Callbacks) OnFetchResponse(*discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {}
