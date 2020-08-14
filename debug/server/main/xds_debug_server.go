package main

import (
	"XDSDebugTools/debug/server"
	"XDSDebugTools/debug/server/redis"
	"context"
	"fmt"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"os"
	//"google.golang.org/grpc"
	//"net"
)

var (
	port   uint   = 18000
	nodeID string = "test-id"
)

func main() {

	// Create a cache
	cache := cachev3.NewSnapshotCache(false, cachev3.IDHash{}, nil)

	// Create the snapshot that we'll serve to Envoy

	snapshot := redis.GenerateSnapshot()
	if err := snapshot.Consistent(); err != nil {
		fmt.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}
	fmt.Println("will serve snapshot %+v", snapshot)

	// Add the snapshot to the cache
	if err := cache.SetSnapshot(nodeID, snapshot); err != nil {
		fmt.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

	// Run the xDS server
	ctx := context.Background()
	cb := &server.Callbacks{}
	srv := serverv3.NewServer(ctx, cache, cb)
	server.RunServer(ctx, srv, port)

	fmt.Println("BYE!!!")

}
