package rdb

/*
#include <stdlib.h>
#include <rocksdb/c.h>
#cgo LDFLAGS: -lstdc++ -lrocksdb
*/
import "C"
import (
	"os"
	"runtime"
)

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func NewCache(dir string, ttl int) *rocksdbCache {
	if !exists(dir) {
		os.Mkdir(dir, 0755)
	}

	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU()))
	C.rocksdb_options_set_create_if_missing(options, 1)
	var e *C.char
	db := C.rocksdb_open_with_ttl(options, C.CString(dir), C.int(ttl), &e)
	if e != nil {
		panic(C.GoString(e))
	}

	C.rocksdb_options_destroy(options)
	c := make(chan *pair, 5)
	wo := C.rocksdb_writeoptions_create()
	go write_func(db, c, wo)
	return &rocksdbCache{db, C.rocksdb_readoptions_create(), C.rocksdb_writeoptions_create(), e, c}
}
