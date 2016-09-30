package bss

import (
	"errors"
	"fmt"
	"github.com/antonikonovalov/didast/example/users"
	"github.com/antonikonovalov/didast/timeid"
	"github.com/golang/protobuf/proto"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/net/context"
	"time"
)

type DataService struct {
	users.UsererServer

	interval time.Duration
	path     string
	table    string
	opened   map[string]*leveldb.DB
}

func NewDataService() users.UsererServer {
	return &DataService{
		interval: time.Hour,
		path:     `/Users/antoniko/go/src/github.com/antonikonovalov/didast/.didastore`,
		table:    `users`,
		opened:   make(map[string]*leveldb.DB),
	}
}

func (s *DataService) particalStore(id int64) *leveldb.DB {
	gap := timeid.DetectInterval(id, s.interval)
	dbName := s.partName(gap[0], gap[1])
	db, ok := s.opened[dbName]
	if !ok {
		db, _ = leveldb.OpenFile(s.path+`/`+dbName, nil)
		s.opened[dbName] = db
	}
	return db
}

func (s *DataService) partName(start, end int64) string {
	return fmt.Sprintf(`%s_%d-%d`, s.table, start, end)
}

func (s *DataService) Put(c context.Context, u *users.User) (*users.Empty, error) {
	db := s.particalStore(u.ID)
	if db == nil {
		return nil, errors.New(`db must be`)
	}
	data, err := proto.Marshal(u)
	if err != nil {
		return nil, err
	}
	err = db.Put([]byte(fmt.Sprintf(`%s`, u.ID)), data, nil)
	if err != nil {
		return nil, err
	}
	return new(users.Empty), nil
}
func (s *DataService) Get(c context.Context, id *users.ID) (*users.User, error) {
	db := s.particalStore(id.ID)
	if db == nil {
		return nil, errors.New(`db must be`)
	}
	data, err := db.Get([]byte(fmt.Sprintf(`%s`, id.ID)), nil)
	if err != nil {
		return nil, err
	}
	user := &users.User{}
	err = proto.Unmarshal(data, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// streams
func (s *DataService) Putter(stream users.Userer_PutterServer) error {
	return nil
}

func (s *DataService) Getter(stream users.Userer_GetterServer) error {
	return nil
}
