package kmdb

import (
	"github.com/tecbot/gorocksdb"
	"log"
)

type KMDB struct {
	Db           *gorocksdb.DB
	Options      *gorocksdb.Options
	ReadOptions  *gorocksdb.ReadOptions
	WriteOptions *gorocksdb.WriteOptions
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
	var dbDir string
	if config.Type.Primary {
		dbDir = config.Store.Dir + "/primary"
	} else {
		dbDir = config.Store.Dir + "/backup"
	}
	rocksdb, err := gorocksdb.OpenDb(opts, dbDir)
	if err != nil {
		log.Panic(err)
	}
	kmdb = KMDB{rocksdb, opts, gorocksdb.NewDefaultReadOptions(), gorocksdb.NewDefaultWriteOptions(), true}
	return &kmdb
}

func (self *KMDB) Put(key, value []byte) error {
	return self.Db.Put(self.WriteOptions, key, value)
}

func (self *KMDB) Get(key []byte) ([]byte, error) {
	slice, err := self.Db.Get(self.ReadOptions, key)
	if err != nil {
		log.Println("Get key:", key, "error:", err)
	}
	log.Println("Get key:", key, "value:", slice.Data())
	var value []byte
	value = append(value, slice.Data()...)
	slice.Free()
	return value, err
}

func (self *KMDB) Del(key []byte) error {
	return self.Db.Delete(self.WriteOptions, key)
}
