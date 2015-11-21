package store

import (
	"github.com/shafreeck/hermes/store/leveldb"
)

func OpenLevelDB(path string) (Store, error) {
	db, err := leveldb.Open(path)
	return db, err
}
