package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	gw "github.com/ouotuo/go-example/rpc/grpc/gateway/proto"
)

var (
	echoEndpoint = flag.String("echo_endpoint", "localhost:9090", "endpoint of YourService")
)

func run() error {
	glog.Infoln("run...")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":9080", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	glog.Info("gw_server running...")

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
