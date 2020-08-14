package main

import (
	"context"
	"fmt"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	base "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	rv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/golang/protobuf/ptypes"
	"log"
	"time"
)
import "google.golang.org/grpc"

const (
	address     = "localhost:18000"
	clusterName = "example_proxy_cluster"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	cds := v3.NewClusterDiscoveryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := rv3.DiscoveryRequest{
		VersionInfo:   "2",
		Node:          &base.Node{Id: "test-id"},
		ResourceNames: []string{clusterName},
		ErrorDetail:   nil}

	r, err := cds.FetchClusters(ctx, &request)

	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	clusters := []*cluster.Cluster{}

	for _, r := range r.GetResources() {
		var any ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(r, &any); err != nil {
			return
		}
		if c, ok := any.Message.(*cluster.Cluster); ok { // v2
			clusters = append(clusters, c)
		}
	}

	if len(clusters) == 0 {
		return
	}

	log.Printf("SUCCESS:", clusters[0].String())

	fmt.Println("BYE!!!")
}
