package main

import (
	"flag"
	"time"

	pb "github.com/antonikonovalov/didast/example/users"
	"github.com/antonikonovalov/didast/timeid"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var DefaultInspector string = `localhost:4567`

var (
	email  = flag.String(`email`, `antoni.konovalov@gmail.com`, `email`)
	entity = flag.String(`entity`, `users`, `name of your entity`)
	id     = flag.Int64(`id`, 0, `id your data`)
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

	store := pb.NewStoreClient(conn)

	switch command {
	case `put`:
		if *id == 0 {
			idN := timeid.New()
			id = &idN
		}

		user := &pb.User{
			ID:        *id,
			Name:      []byte(*email),
			Email:     *email,
			UpdatedAt: timeid.New(),
		}

		data, err := proto.Marshal(user)
		if err != nil {
			grpclog.Fatalf("failed get services: %v", err)
		}
		println(string(data), data)
		u2 := &pb.User{}
		err = proto.Unmarshal(data, u2)
		if err != nil {
			grpclog.Fatalf("failed get services: %v ", err)
		}

		start := time.Now()
		_, err = store.Put(context.Background(), &pb.Object{
			ID:     *id,
			Entity: *entity,
			Data:   string(data),
		})

		end := time.Since(start)

		if err != nil {
			grpclog.Fatalf("failed get services: %v \n %s", err, end.String())
		} else {
			println(`success put to `, *entity, u2.String(), "\n", end.String())
		}
	case `get`:
		if *id == 0 {
			println(`id == 0`)
			return
		}
		start := time.Now()
		data, err := store.Get(context.Background(), &pb.ID{ID: *id, Entity: *entity})
		end := time.Since(start)
		if err != nil {
			grpclog.Fatalf("failed get services: %v \n %s", err, end.String())
		} else {
			user := &pb.User{}
			println(data.String(), string(data.Data))
			println(data)
			err = proto.Unmarshal([]byte(data.Data), user)
			if err != nil {
				grpclog.Fatalf("failed get services: %v \n %s", err, end.String())
			}
			println(`success get:`, user.String(), "\n", end.String())
		}
	}
}
