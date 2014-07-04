package kmdb

import (
	"github.com/tecbot/gorocksdb"
	"log"
)

type KMDB struct {
	Db           gorocksdb.DB
	Options      gorocksdb.Options
	ReadOptions  gorocksdb.ReadOptions
	WriteOptions gorocksdb.WriteOptions
	Primary      bool
}

func (self *KMDB) Close() {
	if self.Db != nil {
		self.Db.Close()
	}
}

func Open(config *Config) *KMDB {
	kmdb := KMDB{}
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	rocksdb, err := gorocksdb.OpenDb(opts, config.Store.Dir)
	if err != nil {
		log.Panic(err)
	}
	kmdb.Db = KMDB{rocksdb, opts, gorocksdb.NewDefaultReadOptions(), gorocksdb.NewDefaultWriteOptions(), true}
	return &kmdb
}

func (self *KMDB) Put(key, value []byte) error {
	return self.Db.Put(self.WriteOptions, key, value)
}

func (self *KMDB) Get(key []byte) (*[]byte, error) {
	slice, err := self.Db.Get(self.ReadOptions, key)
	value := slice.Data()
	slice.Free()
	return value, err
}

func (self *KMDB) Del(key []byte) error {
	return self.Db.Delete(self.WriteOptions, key)
}
