package kmdb

import (
	"github.com/tecbot/gorocksdb"
	"log"
)

type Sencondary struct {
	Db	gorocksdb.DB
	SlaveOf	string
}
