package main

import (
	"flag"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/antonikonovalov/didast/bss"
	pb "github.com/antonikonovalov/didast/example/users"
)

var (
	addr = flag.String(`addr`, `localhost:4567`, `address for listen service`)
	path = flag.String(`path`, ``, `address for listen service`)
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterStoreServer(grpcServer, bss.NewDataService(*path))
	grpcServer.Serve(lis)
}
