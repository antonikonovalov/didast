package main

import (
	"flag"
	"net"

	"github.com/antonikonovalov/didast/bss"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"fmt"
	_ "net/http/pprof"

	pb "github.com/antonikonovalov/didast/example/users"
	"net/http"
)

var (
	addr = flag.String(`addr`, `localhost:4567`, `address for listen service`)
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	go http.ListenAndServe(fmt.Sprintf(":%d", 36662), nil)
	pb.RegisterUsererServer(grpcServer, bss.NewDataService())
	grpcServer.Serve(lis)
}
