package main

import (
	"flag"
	pb "github.com/antonikonovalov/didast/example/users"
	"github.com/antonikonovalov/didast/timeid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"time"
)

var DefaultInspector string = `localhost:4567`

var (
	email = flag.String(`email`, `antoni.konovalov@gmail.com`, `email`)
	id    = flag.Int64(`id`, 0, `id your data`)
)

func main() {
	flag.Parse()

	command := flag.Arg(0)

	// Set up a connection to the lookupd services
	conn, err := grpc.Dial(DefaultInspector, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("failed to listen inspector services: %v", err)
	}
	defer conn.Close()

	store := pb.NewUsererClient(conn)

	switch command {
	case `put`:
		if *id == 0 {
			idN := timeid.New()
			id = &idN
		}
		user := &pb.User{
			ID:         *id,
			Name:       []byte(*email),
			Email:      *email,
			UpdateddAt: timeid.New(),
		}

		start := time.Now()
		_, err = store.Put(context.Background(), user)
		end := time.Since(start)

		if err != nil {
			grpclog.Fatalf("failed get services: %v \n %s", err, end.String())
		} else {
			println(`success put:`, user.String(), "\n", end.String())
		}
	case `get`:
		if *id == 0 {
			println(`id == 0`)
			return
		}
		start := time.Now()
		user, err := store.Get(context.Background(), &pb.ID{ID: *id})
		end := time.Since(start)
		if err != nil {
			grpclog.Fatalf("failed get services: %v \n %s", err, end.String())
		} else {
			println(`success get:`, user.String(), "\n", end.String())
		}
	}
}
