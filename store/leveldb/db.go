package leveldb

/*
#cgo LDFLAGS: -lleveldb
#include <leveldb/c.h>
*/
import "C"

import (
	"unsafe"
)

type DBOptions struct {
	Opts *C.leveldb_options_t
}
type ReadOptions struct {
	Opts *C.leveldb_readoptions_t
}
type WriteOptions struct {
	Opts *C.leveldb_writeoptions_t
}

type LevelDB struct {
	db        *C.leveldb_t
	path      string
	dbOpts    *DBOptions
	readOpts  *ReadOptions
	writeOpts *WriteOptions
}

type DBError struct {
	Err string
}

func (e *DBError) Error() string {
	return e.Err
}

func Open(path string) (*LevelDB, error) {
	var err *C.char
	db := &LevelDB{}
	db.dbOpts = &DBOptions{C.leveldb_options_create()}
	db.readOpts = &ReadOptions{C.leveldb_readoptions_create()}
	db.writeOpts = &WriteOptions{C.leveldb_writeoptions_create()}

	C.leveldb_options_set_create_if_missing(db.dbOpts.Opts, 1)
	db.db = C.leveldb_open(db.dbOpts.Opts, C.CString(path), &err)
	if db.db == nil {
		e := &DBError{}
		e.Err = "Open db failed"
		return nil, e
	}
	if err != nil {
		e := &DBError{}
		e.Err = C.GoString(err)
		return db, e
	}
	return db, nil
}
func (db *LevelDB) Get(key []byte) ([]byte, error) {
	var vlen C.size_t
	var err *C.char
	val := C.leveldb_get(db.db, db.readOpts.Opts, C.CString(string(key)), C.size_t(len(key)), &vlen, &err)
	if err != nil {
		e := &DBError{}
		e.Err = C.GoString(err)
		return nil, e
	}
	if vlen == 0 {
		return nil, nil
	}
	return C.GoBytes(unsafe.Pointer(val), C.int(vlen)), nil
}
func (db *LevelDB) Set(key, value []byte) error {
	var err *C.char
	C.leveldb_put(db.db, db.writeOpts.Opts, C.CString(string(key)), C.size_t(len(key)),
		C.CString(string(value)), C.size_t(len(value)), &err)
	if err != nil {
		e := &DBError{}
		e.Err = C.GoString(err)
		return e
	}
	return nil
}
func (db *LevelDB) Delete(key []byte) error {
	var err *C.char
	C.leveldb_delete(db.db, db.writeOpts.Opts, C.CString(string(key)), C.size_t(len(key)),
		&err)
	if err != nil {
		e := &DBError{}
		e.Err = C.GoString(err)
		return e
	}
	return nil
}
func (db *LevelDB) Scanner(begin, end []byte) {
}
