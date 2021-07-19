package rdb

/*
#include <stdlib.h>
#include <rocksdb/c.h>
#cgo LDFLAGS: -lstdc++ -lrocksdb
*/
import "C"
import (
	"errors"
	"time"
	"unsafe"
)

const BATCH_SIZE = 100

type pair struct {
	k string
	v []byte
}

type rocksdbCache struct {
	db *C.rocksdb_t              // rocksdb_t 类型的指针
	ro *C.rocksdb_readoptions_t  // rocksdb 读操作选项类型
	wo *C.rocksdb_writeoptions_t // rocksdb 写操作选项类型
	e  *C.char                   //char*
	ch chan *pair
}

func (c *rocksdbCache) Get(key string) ([]byte, error) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var length C.size_t
	v := C.rocksdb_get(c.db, c.ro, k, C.size_t(len(key)), &length, &c.e)
	if c.e != nil {
		return nil, errors.New(C.GoString(c.e))

	}
	defer C.free(unsafe.Pointer(v))
	return C.GoBytes(unsafe.Pointer(v), C.int(length)), nil

}

func (c *rocksdbCache) Set(key string, value []byte) error {
	c.ch <- &pair{key, value}
	return nil

}

func (c *rocksdbCache) Del(key string) error {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	C.rocksdb_delete(c.db, c.wo, k, C.size_t(len(key)), &c.e)
	if c.e != nil {
		return errors.New(C.GoString(c.e))

	}
	return nil

}

func flush_batch(db *C.rocksdb_t, b *C.rocksdb_writebatch_t, o *C.rocksdb_writeoptions_t) {
	var e *C.char
	C.rocksdb_write(db, o, b, &e)
	if e != nil {
		panic(C.GoString(e))

	}
	C.rocksdb_writebatch_clear(b)
}

func write_func(db *C.rocksdb_t, c chan *pair, o *C.rocksdb_writeoptions_t) {
	count := 0
	t := time.NewTimer(time.Second)
	b := C.rocksdb_writebatch_create()
	for {
		select {
		case p := <-c:
			count++
			key := C.CString(p.k)
			value := C.CBytes(p.v)
			C.rocksdb_writebatch_put(b, key, C.size_t(len(p.k)), (*C.char)(value), C.size_t(len(p.v)))
			C.free(unsafe.Pointer(key))
			C.free(value)
			if count == BATCH_SIZE {
				flush_batch(db, b, o)
				count = 0

			}
			// 如果在处理非计时器触发的事件时计时时间到 timer中就会有一个计时器留在管道口
			// 重置计时器需要小心 如果不能确定一个计时器已经触发 就需要先调用stop()停止这个计时器
			// 然后取走timer管道中的计时器
			if !t.Stop() {
				<-t.C
			}
			t.Reset(time.Second)
		case <-t.C:
			if count != 0 {
				flush_batch(db, b, o)
				count = 0

			}
			t.Reset(time.Second)
		}
	}
}
