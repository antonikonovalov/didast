package bss

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/antonikonovalov/didast/example/users"
	"github.com/antonikonovalov/didast/timeid"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/net/context"
)

type DataService struct {
	users.StoreServer

	interval time.Duration
	path     string
	opened   map[string]*leveldb.DB
}

func NewDataService(path string) users.StoreServer {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if path == "" {
		path = currentDir
	}

	return &DataService{
		interval: time.Hour,
		path:     path + `/.didastore`,
		opened:   make(map[string]*leveldb.DB),
	}
}

func (s *DataService) particalStore(entity string, id int64) *leveldb.DB {
	gap := timeid.DetectInterval(id, s.interval)
	dbName := partName(entity, gap[0], gap[1])
	db, ok := s.opened[dbName]
	if !ok {
		db, _ = leveldb.OpenFile(s.path+`/`+dbName, nil)
		s.opened[dbName] = db
	}
	return db
}

func partName(entity string, start, end int64) string {
	return fmt.Sprintf(`%s_%d-%d`, entity, start, end)
}

func (s *DataService) Put(c context.Context, u *users.Object) (*users.Empty, error) {
	db := s.particalStore(u.Entity, u.ID)
	if db == nil {
		return nil, errors.New(`db must be`)
	}
	err := db.Put([]byte(fmt.Sprintf(`%s`, u.ID)), []byte(u.Data), nil)
	if err != nil {
		return nil, err
	}
	return new(users.Empty), nil
}
func (s *DataService) Get(c context.Context, id *users.ID) (*users.Object, error) {
	db := s.particalStore(id.Entity, id.ID)
	if db == nil {
		return nil, errors.New(`db must be`)
	}
	data, err := db.Get([]byte(fmt.Sprintf(`%s`, id.ID)), nil)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &users.Object{
		ID:     id.ID,
		Entity: id.Entity,
		Data:   string(data),
	}, nil
}

// streams
func (s *DataService) Putter(stream users.Store_PutterServer) error {
	return nil
}

func (s *DataService) Getter(stream users.Store_GetterServer) error {
	return nil
}
