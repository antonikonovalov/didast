package benchmarks

import (
	"testing"

	"github.com/antonikonovalov/didast/timeid"
	"github.com/gogo/protobuf/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/antonikonovalov/didast/example/users"
)

func BenchmarkPutObjects(b *testing.B) {
	conn, err := grpc.Dial(`localhost:4567`, grpc.WithInsecure())
	if err != nil {
		b.Fatalf("failed to listen inspector services: %v", err)
	}
	defer conn.Close()

	store := pb.NewStoreClient(conn)
	user := &pb.User{
		ID:        timeid.New(),
		Name:      []byte(`antni.konovalovsd`),
		Email:     `antoni.konovalov@sdfsdfksd.com`,
		UpdatedAt: timeid.New(),
	}

	data, err := proto.Marshal(user)
	if err != nil {
		b.Fatalf("failed get services: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pbt *testing.PB) {
		for pbt.Next() {
			obj := &pb.Object{
				ID:     timeid.New(),
				Entity: `users`,
				Data:   string(data),
			}
			_, err = store.Put(context.Background(), obj)
			if err != nil {
				b.Error(err)
			}
		}
	})
}
